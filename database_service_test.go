package main

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *DatabaseService {
	// Create temporary database file
	dbPath := "test_database.db"

	// Clean up any existing test database
	os.Remove(dbPath)

	// Create new database service
	db, err := NewDatabaseService(dbPath)
	require.NoError(t, err)

	// Clean up function
	t.Cleanup(func() {
		db.Close()
		os.Remove(dbPath)
	})

	return db
}

func TestDatabaseService_TweetOperations(t *testing.T) {
	db := setupTestDB(t)

	// Test data
	tweet := TweetModel{
		ID:          "tweet_123",
		Text:        "This is a test tweet about $DARK",
		CreatedAt:   time.Now(),
		ReplyCount:  5,
		UserID:      "user_456",
		InReplyToID: "",
	}

	// Test SaveTweet
	t.Run("SaveTweet", func(t *testing.T) {
		err := db.SaveTweet(tweet)
		assert.NoError(t, err)
	})

	// Test TweetExists
	t.Run("TweetExists", func(t *testing.T) {
		exists := db.TweetExists("tweet_123")
		assert.True(t, exists)

		notExists := db.TweetExists("nonexistent_tweet")
		assert.False(t, notExists)
	})

	// Test GetTweet
	t.Run("GetTweet", func(t *testing.T) {
		retrievedTweet, err := db.GetTweet("tweet_123")
		assert.NoError(t, err)
		assert.Equal(t, tweet.ID, retrievedTweet.ID)
		assert.Equal(t, tweet.Text, retrievedTweet.Text)
		assert.Equal(t, tweet.ReplyCount, retrievedTweet.ReplyCount)
		assert.Equal(t, tweet.UserID, retrievedTweet.UserID)
	})

	// Test GetTweetReplyCount
	t.Run("GetTweetReplyCount", func(t *testing.T) {
		replyCount, err := db.GetTweetReplyCount("tweet_123")
		assert.NoError(t, err)
		assert.Equal(t, 5, replyCount)
	})

	// Test UpdateTweetReplyCount
	t.Run("UpdateTweetReplyCount", func(t *testing.T) {
		err := db.UpdateTweetReplyCount("tweet_123", 10)
		assert.NoError(t, err)

		replyCount, err := db.GetTweetReplyCount("tweet_123")
		assert.NoError(t, err)
		assert.Equal(t, 10, replyCount)
	})

	// Test DeleteTweet
	t.Run("DeleteTweet", func(t *testing.T) {
		err := db.DeleteTweet("tweet_123")
		assert.NoError(t, err)

		exists := db.TweetExists("tweet_123")
		assert.False(t, exists)
	})
}

func TestDatabaseService_UserOperations(t *testing.T) {
	db := setupTestDB(t)

	// Test data
	user := UserModel{
		ID:       "user_123",
		Username: "testuser",
		Name:     "Test User",
		IsFUD:    false,
	}

	// Test SaveUser
	t.Run("SaveUser", func(t *testing.T) {
		err := db.SaveUser(user)
		assert.NoError(t, err)
	})

	// Test UserExists
	t.Run("UserExists", func(t *testing.T) {
		exists := db.UserExists("user_123")
		assert.True(t, exists)

		notExists := db.UserExists("nonexistent_user")
		assert.False(t, notExists)
	})

	// Test UserExistsByUsername
	t.Run("UserExistsByUsername", func(t *testing.T) {
		exists := db.UserExistsByUsername("testuser")
		assert.True(t, exists)

		notExists := db.UserExistsByUsername("nonexistent_username")
		assert.False(t, notExists)
	})

	// Test GetUser
	t.Run("GetUser", func(t *testing.T) {
		retrievedUser, err := db.GetUser("user_123")
		assert.NoError(t, err)
		assert.Equal(t, user.ID, retrievedUser.ID)
		assert.Equal(t, user.Username, retrievedUser.Username)
		assert.Equal(t, user.Name, retrievedUser.Name)
		assert.Equal(t, user.IsFUD, retrievedUser.IsFUD)
	})

	// Test GetUserByUsername
	t.Run("GetUserByUsername", func(t *testing.T) {
		retrievedUser, err := db.GetUserByUsername("testuser")
		assert.NoError(t, err)
		assert.Equal(t, user.ID, retrievedUser.ID)
		assert.Equal(t, user.Username, retrievedUser.Username)
	})

	// Test DeleteUser
	t.Run("DeleteUser", func(t *testing.T) {
		err := db.DeleteUser("user_123")
		assert.NoError(t, err)

		exists := db.UserExists("user_123")
		assert.False(t, exists)
	})
}

func TestDatabaseService_FUDUserOperations(t *testing.T) {
	db := setupTestDB(t)

	// Test data
	fudUser := FUDUserModel{
		UserID:         "fud_user_123",
		Username:       "fuduser",
		FUDType:        "professional_trojan_horse",
		FUDProbability: 0.95,
		DetectedAt:     time.Now(),
		MessageCount:   1,
		LastMessageID:  "msg_123",
	}

	// Test SaveFUDUser
	t.Run("SaveFUDUser", func(t *testing.T) {
		err := db.SaveFUDUser(fudUser)
		assert.NoError(t, err)
	})

	// Test IsFUDUser
	t.Run("IsFUDUser", func(t *testing.T) {
		isFUD := db.IsFUDUser("fud_user_123")
		assert.True(t, isFUD)

		notFUD := db.IsFUDUser("clean_user_123")
		assert.False(t, notFUD)
	})

	// Test GetFUDUser
	t.Run("GetFUDUser", func(t *testing.T) {
		retrievedFUDUser, err := db.GetFUDUser("fud_user_123")
		assert.NoError(t, err)
		assert.Equal(t, fudUser.UserID, retrievedFUDUser.UserID)
		assert.Equal(t, fudUser.Username, retrievedFUDUser.Username)
		assert.Equal(t, fudUser.FUDType, retrievedFUDUser.FUDType)
		assert.Equal(t, fudUser.FUDProbability, retrievedFUDUser.FUDProbability)
		assert.Equal(t, fudUser.MessageCount, retrievedFUDUser.MessageCount)
	})

	// Test IncrementFUDUserMessageCount
	t.Run("IncrementFUDUserMessageCount", func(t *testing.T) {
		err := db.IncrementFUDUserMessageCount("fud_user_123", "msg_456")
		assert.NoError(t, err)

		updatedFUDUser, err := db.GetFUDUser("fud_user_123")
		assert.NoError(t, err)
		assert.Equal(t, 2, updatedFUDUser.MessageCount)
		assert.Equal(t, "msg_456", updatedFUDUser.LastMessageID)
	})

	// Test GetAllFUDUsers
	t.Run("GetAllFUDUsers", func(t *testing.T) {
		// Add another FUD user
		fudUser2 := FUDUserModel{
			UserID:         "fud_user_456",
			Username:       "anotherfuduser",
			FUDType:        "emotional_escalation",
			FUDProbability: 0.85,
			DetectedAt:     time.Now(),
			MessageCount:   3,
			LastMessageID:  "msg_789",
		}
		err := db.SaveFUDUser(fudUser2)
		assert.NoError(t, err)

		allFUDUsers, err := db.GetAllFUDUsers()
		assert.NoError(t, err)
		assert.Len(t, allFUDUsers, 2)
	})

	// Test DeleteFUDUser
	t.Run("DeleteFUDUser", func(t *testing.T) {
		err := db.DeleteFUDUser("fud_user_123")
		assert.NoError(t, err)

		isFUD := db.IsFUDUser("fud_user_123")
		assert.False(t, isFUD)
	})
}

func TestDatabaseService_RelationshipOperations(t *testing.T) {
	db := setupTestDB(t)

	// Create a user first
	user := UserModel{
		ID:       "user_rel_123",
		Username: "reluser",
		Name:     "Relationship User",
	}
	err := db.SaveUser(user)
	require.NoError(t, err)

	// Create tweets
	mainTweet := TweetModel{
		ID:         "main_tweet_123",
		Text:       "This is a main tweet about $DARK",
		CreatedAt:  time.Now(),
		ReplyCount: 2,
		UserID:     "user_rel_123",
	}

	replyTweet1 := TweetModel{
		ID:          "reply_tweet_1",
		Text:        "This is reply 1",
		CreatedAt:   time.Now().Add(1 * time.Minute),
		ReplyCount:  0,
		UserID:      "user_rel_123",
		InReplyToID: "main_tweet_123",
	}

	replyTweet2 := TweetModel{
		ID:          "reply_tweet_2",
		Text:        "This is reply 2",
		CreatedAt:   time.Now().Add(2 * time.Minute),
		ReplyCount:  0,
		UserID:      "user_rel_123",
		InReplyToID: "main_tweet_123",
	}

	// Save tweets
	err = db.SaveTweet(mainTweet)
	require.NoError(t, err)
	err = db.SaveTweet(replyTweet1)
	require.NoError(t, err)
	err = db.SaveTweet(replyTweet2)
	require.NoError(t, err)

	// Test GetTweetsByUser
	t.Run("GetTweetsByUser", func(t *testing.T) {
		userTweets, err := db.GetTweetsByUser("user_rel_123")
		assert.NoError(t, err)
		assert.Len(t, userTweets, 3)
	})

	// Test GetRepliesForTweet
	t.Run("GetRepliesForTweet", func(t *testing.T) {
		replies, err := db.GetRepliesForTweet("main_tweet_123")
		assert.NoError(t, err)
		assert.Len(t, replies, 2)
		assert.Equal(t, "reply_tweet_1", replies[0].ID)
		assert.Equal(t, "reply_tweet_2", replies[1].ID)
	})
}

func TestDatabaseService_SearchOperations(t *testing.T) {
	db := setupTestDB(t)

	// Create test tweets
	tweets := []TweetModel{
		{
			ID:        "search_tweet_1",
			Text:      "This tweet mentions $DARK token",
			CreatedAt: time.Now(),
			UserID:    "user_1",
		},
		{
			ID:        "search_tweet_2",
			Text:      "Another tweet about cryptocurrency",
			CreatedAt: time.Now().Add(-1 * time.Hour),
			UserID:    "user_2",
		},
		{
			ID:        "search_tweet_3",
			Text:      "Dark mode is great for coding",
			CreatedAt: time.Now().Add(-2 * time.Hour),
			UserID:    "user_3",
		},
	}

	for _, tweet := range tweets {
		err := db.SaveTweet(tweet)
		require.NoError(t, err)
	}

	// Test SearchTweets
	t.Run("SearchTweets", func(t *testing.T) {
		results, err := db.SearchTweets("DARK", 10)
		assert.NoError(t, err)
		assert.Len(t, results, 2) // Should find both tweets containing "DARK"
	})

	// Test GetRecentTweets
	t.Run("GetRecentTweets", func(t *testing.T) {
		recentTweets, err := db.GetRecentTweets(2)
		assert.NoError(t, err)
		assert.Len(t, recentTweets, 2)
		// Should be ordered by created_at DESC
		assert.Equal(t, "search_tweet_1", recentTweets[0].ID)
	})

	// Test count operations
	t.Run("CountOperations", func(t *testing.T) {
		tweetCount, err := db.GetTweetCount()
		assert.NoError(t, err)
		assert.Equal(t, int64(3), tweetCount)

		userCount, err := db.GetUserCount()
		assert.NoError(t, err)
		assert.Equal(t, int64(0), userCount) // No users saved in this test

		fudUserCount, err := db.GetFUDUserCount()
		assert.NoError(t, err)
		assert.Equal(t, int64(0), fudUserCount) // No FUD users saved in this test
	})
}

func TestDatabaseService_ComplexScenario(t *testing.T) {
	db := setupTestDB(t)

	// Simulate a complete scenario: user creates tweets, gets flagged as FUD

	// 1. Create user
	user := UserModel{
		ID:       "complex_user_123",
		Username: "complexuser",
		Name:     "Complex Test User",
	}
	err := db.SaveUser(user)
	require.NoError(t, err)

	// 2. User creates some tweets
	tweets := []TweetModel{
		{
			ID:         "complex_tweet_1",
			Text:       "First tweet about $DARK",
			CreatedAt:  time.Now().Add(-3 * time.Hour),
			ReplyCount: 0,
			UserID:     "complex_user_123",
		},
		{
			ID:         "complex_tweet_2",
			Text:       "Second tweet spreading FUD about $DARK",
			CreatedAt:  time.Now().Add(-2 * time.Hour),
			ReplyCount: 5,
			UserID:     "complex_user_123",
		},
		{
			ID:          "complex_tweet_3",
			Text:        "Reply to my own tweet",
			CreatedAt:   time.Now().Add(-1 * time.Hour),
			ReplyCount:  0,
			UserID:      "complex_user_123",
			InReplyToID: "complex_tweet_2",
		},
	}

	for _, tweet := range tweets {
		err := db.SaveTweet(tweet)
		require.NoError(t, err)
	}

	// 3. User gets flagged as FUD
	fudUser := FUDUserModel{
		UserID:         "complex_user_123",
		Username:       "complexuser",
		FUDType:        "professional_direct_attack",
		FUDProbability: 0.92,
		DetectedAt:     time.Now(),
		MessageCount:   1,
		LastMessageID:  "complex_tweet_2",
	}
	err = db.SaveFUDUser(fudUser)
	require.NoError(t, err)

	// 4. Verify the complete scenario

	// Check user exists
	assert.True(t, db.UserExists("complex_user_123"))

	// Check user tweets
	userTweets, err := db.GetTweetsByUser("complex_user_123")
	assert.NoError(t, err)
	assert.Len(t, userTweets, 3)

	// Check FUD status
	assert.True(t, db.IsFUDUser("complex_user_123"))

	// Check replies
	replies, err := db.GetRepliesForTweet("complex_tweet_2")
	assert.NoError(t, err)
	assert.Len(t, replies, 1)
	assert.Equal(t, "complex_tweet_3", replies[0].ID)

	// Search for FUD content
	fudTweets, err := db.SearchTweets("FUD", 10)
	assert.NoError(t, err)
	assert.Len(t, fudTweets, 1)
	assert.Equal(t, "complex_tweet_2", fudTweets[0].ID)
}
