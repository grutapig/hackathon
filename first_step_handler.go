package main

import (
	"encoding/json"
	"fmt"
	"github.com/grutapig/hackaton/twitterapi"
	"log"
)

// FirstStepHandler handles new message first step analysis
func FirstStepHandler(newMessageCh chan twitterapi.NewMessage, fudChannel chan twitterapi.NewMessage, claudeApi *ClaudeApi, twitterApi *twitterapi.TwitterAPIService, systemPromptFirstStep []byte, userStatusManager *UserStatusManager, dbService *DatabaseService) {
	defer close(fudChannel)

	for newMessage := range newMessageCh {
		log.Println("Got a new message:", newMessage.Author.UserName, " - ", newMessage.Text, "reply to:", newMessage.ParentTweet.Text)

		// First step analysis for new or clean users
		messages := ClaudeMessages{}

		// Add thread context in order: grandparent -> parent -> current
		if newMessage.GrandParentTweet.ID != "" {
			messages = append(messages, ClaudeMessage{ROLE_USER, "the main post is: " + newMessage.GrandParentTweet.Author + ":" + newMessage.GrandParentTweet.Text})
			messages = append(messages, ClaudeMessage{ROLE_USER, "reply in thread: " + newMessage.ParentTweet.Author + ":" + newMessage.ParentTweet.Text})
		} else {
			messages = append(messages, ClaudeMessage{ROLE_USER, "the main post is: " + newMessage.ParentTweet.Author + ":" + newMessage.ParentTweet.Text})
		}

		messages = append(messages, ClaudeMessage{ROLE_USER, "user reply being analyzed: " + newMessage.Author.UserName + ":" + newMessage.Text})
		messages = append(messages, ClaudeMessage{ROLE_ASSISTANT, "{"})
		fmt.Println(messages)
		resp, err := claudeApi.SendMessage(messages, fmt.Sprintf("%s\n<instruction>you must analyze %s user messages in the context of the full thread</instruction>", string(systemPromptFirstStep), newMessage.Author.UserName))
		if err != nil {
			log.Printf("error claude: %s", err)
			continue
		}

		fmt.Println(resp.Content[0].Text)
		aiDecision := FirstStepClaudeResponse{}
		err = json.Unmarshal([]byte("{"+resp.Content[0].Text), &aiDecision)
		fmt.Println("ai decision", aiDecision, resp.Content[0].Text)

		if aiDecision.IsFud {
			// Mark user as being analyzed to prevent duplicate analysis
			userStatusManager.SetUserAnalyzing(newMessage.Author.ID, newMessage.Author.UserName)
			fudChannel <- newMessage
		}
	}
}
