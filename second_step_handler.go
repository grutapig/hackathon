package main

import (
	"encoding/json"
	"fmt"
	"github.com/grutapig/hackaton/twitterapi"
	"log"
	"time"
)

// SecondStepHandler handles flagged messages second step analysis
func SecondStepHandler(fudChannel chan twitterapi.NewMessage, notificationCh chan FUDAlertNotification, twitterApi *twitterapi.TwitterAPIService, claudeApi *ClaudeApi, systemPromptSecondStep []byte, userStatusManager *UserStatusManager, ticker string, dbService *DatabaseService) {
	defer close(notificationCh)

	for newMessage := range fudChannel {
		// Get user's ticker mentions using advanced search (max 3 pages)
		userTickerMentions := getUserTickerMentions(twitterApi, newMessage.Author.UserName, ticker)
		followers, err := twitterApi.GetUserFollowers(twitterapi.UserFollowersRequest{UserName: newMessage.Author.UserName})
		followings, err := twitterApi.GetUserFollowings(twitterapi.UserFollowingsRequest{UserName: newMessage.Author.UserName})

		// Prepare claude request
		claudeMessages := PrepareClaudeSecondStepRequest(userTickerMentions, followers, followings, userStatusManager)

		// Add thread context in order: grandparent -> parent -> current
		if newMessage.GrandParentTweet.ID != "" {
			claudeMessages = append(claudeMessages, ClaudeMessage{ROLE_USER, "the main post is: " + newMessage.GrandParentTweet.Author + ":" + newMessage.GrandParentTweet.Text})
			claudeMessages = append(claudeMessages, ClaudeMessage{ROLE_USER, "reply in thread: " + newMessage.ParentTweet.Author + ":" + newMessage.ParentTweet.Text})
		} else {
			claudeMessages = append(claudeMessages, ClaudeMessage{ROLE_USER, "the main post is: " + newMessage.ParentTweet.Author + ":" + newMessage.ParentTweet.Text})
		}

		claudeMessages = append(claudeMessages, ClaudeMessage{ROLE_USER, "user reply being analyzed: " + newMessage.Author.UserName + ":" + newMessage.Text})
		claudeMessages = append(claudeMessages, ClaudeMessage{Role: ROLE_ASSISTANT, Content: "{"})

		fmt.Println("send to analyze:", claudeMessages)
		resp, err := claudeApi.SendMessage(claudeMessages, string(systemPromptSecondStep)+"analyzed user is "+newMessage.Author.UserName)
		aiDecision2 := SecondStepClaudeResponse{}
		fmt.Println("claude make a decision for this user:", resp, err)

		if err != nil {
			log.Printf("error claude second step: %s", err)
			continue
		}

		err = json.Unmarshal([]byte("{"+resp.Content[0].Text), &aiDecision2)
		if err != nil {
			log.Printf("error unmarshaling claude response: %s", err)
			continue
		}

		fmt.Println(aiDecision2)

		// Update user status after analysis
		userStatusManager.UpdateUserAfterAnalysis(newMessage.Author.ID, newMessage.Author.UserName, aiDecision2, newMessage.TweetID)

		if aiDecision2.IsFUDMessage {
			// Store FUD user in database
			if aiDecision2.IsFUDUser {
				fudUser := FUDUserModel{
					UserID:         newMessage.Author.ID,
					Username:       newMessage.Author.UserName,
					FUDType:        aiDecision2.FUDType,
					FUDProbability: aiDecision2.FUDProbability,
					DetectedAt:     time.Now(),
					MessageCount:   1,
					LastMessageID:  newMessage.TweetID,
				}

				// Check if FUD user already exists
				if dbService.IsFUDUser(newMessage.Author.ID) {
					// Increment message count for existing FUD user
					err = dbService.IncrementFUDUserMessageCount(newMessage.Author.ID, newMessage.TweetID)
					if err != nil {
						log.Printf("Failed to increment FUD user message count: %v", err)
					}
				} else {
					// Save new FUD user
					err = dbService.SaveFUDUser(fudUser)
					if err != nil {
						log.Printf("Failed to save FUD user: %v", err)
					} else {
						log.Printf("Stored new FUD user: %s", newMessage.Author.UserName)
					}
				}
			}

			// Create FUD alert notification
			alert := FUDAlertNotification{
				FUDMessageID:      newMessage.TweetID,
				FUDUserID:         newMessage.Author.ID,
				FUDUsername:       newMessage.Author.UserName,
				ThreadID:          newMessage.ReplyTweetID,
				DetectedAt:        time.Now().Format(time.RFC3339),
				AlertSeverity:     mapRiskLevelToSeverity(aiDecision2.UserRiskLevel),
				FUDType:           aiDecision2.FUDType,
				FUDProbability:    aiDecision2.FUDProbability,
				MessagePreview:    newMessage.Text,
				RecommendedAction: getRecommendedAction(aiDecision2),
				KeyEvidence:       aiDecision2.KeyEvidence,
				DecisionReason:    aiDecision2.DecisionReason,
			}
			notificationCh <- alert
		}
	}
}

func mapRiskLevelToSeverity(riskLevel string) string {
	switch riskLevel {
	case "critical":
		return "critical"
	case "high":
		return "high"
	case "medium":
		return "medium"
	case "low":
		return "low"
	default:
		return "medium"
	}
}

func getRecommendedAction(decision SecondStepClaudeResponse) string {
	if decision.UserRiskLevel == "critical" {
		return "IMMEDIATE_ACTION_REQUIRED"
	} else if decision.UserRiskLevel == "high" {
		return "MONITOR_CLOSELY"
	} else {
		return "STANDARD_MONITORING"
	}
}
