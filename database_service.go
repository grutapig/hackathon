package main

import (
	"fmt"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DatabaseService struct {
	db *gorm.DB
}

// NewDatabaseService creates a new database service instance
func NewDatabaseService(dbPath string) (*DatabaseService, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // Silent to reduce log noise
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	service := &DatabaseService{db: db}

	// Run migrations
	if err := service.runMigrations(); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return service, nil
}

// runMigrations runs database migrations
func (s *DatabaseService) runMigrations() error {
	return s.db.AutoMigrate(&TweetModel{}, &UserModel{}, &FUDUserModel{})
}

// Tweet related methods

// SaveTweet saves or updates a tweet in the database
func (s *DatabaseService) SaveTweet(tweet TweetModel) error {
	tweet.UpdatedAt = time.Now()
	return s.db.Save(&tweet).Error
}

// GetTweet retrieves a tweet by ID
func (s *DatabaseService) GetTweet(id string) (*TweetModel, error) {
	var tweet TweetModel
	err := s.db.Where("id = ?", id).First(&tweet).Error
	if err != nil {
		return nil, err
	}
	return &tweet, nil
}

// TweetExists checks if a tweet exists in the database
func (s *DatabaseService) TweetExists(id string) bool {
	var count int64
	s.db.Model(&TweetModel{}).Where("id = ?", id).Count(&count)
	return count > 0
}

// GetTweetReplyCount gets the reply count for a tweet
func (s *DatabaseService) GetTweetReplyCount(id string) (int, error) {
	var tweet TweetModel
	err := s.db.Select("reply_count").Where("id = ?", id).First(&tweet).Error
	if err != nil {
		return 0, err
	}
	return tweet.ReplyCount, nil
}

// UpdateTweetReplyCount updates the reply count for a tweet
func (s *DatabaseService) UpdateTweetReplyCount(id string, replyCount int) error {
	return s.db.Model(&TweetModel{}).Where("id = ?", id).Update("reply_count", replyCount).Error
}

// GetTweetsByUser retrieves all tweets by a specific user
func (s *DatabaseService) GetTweetsByUser(userID string) ([]TweetModel, error) {
	var tweets []TweetModel
	err := s.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&tweets).Error
	return tweets, err
}

// GetRepliesForTweet retrieves all replies for a specific tweet
func (s *DatabaseService) GetRepliesForTweet(tweetID string) ([]TweetModel, error) {
	var replies []TweetModel
	err := s.db.Where("in_reply_to_id = ?", tweetID).Order("created_at ASC").Find(&replies).Error
	return replies, err
}

// DeleteTweet deletes a tweet from the database
func (s *DatabaseService) DeleteTweet(id string) error {
	return s.db.Delete(&TweetModel{}, "id = ?", id).Error
}

// User related methods

// SaveUser saves or updates a user in the database
func (s *DatabaseService) SaveUser(user UserModel) error {
	user.UpdatedAt = time.Now()
	return s.db.Save(&user).Error
}

// GetUser retrieves a user by ID
func (s *DatabaseService) GetUser(id string) (*UserModel, error) {
	var user UserModel
	err := s.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByUsername retrieves a user by username
func (s *DatabaseService) GetUserByUsername(username string) (*UserModel, error) {
	var user UserModel
	err := s.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UserExists checks if a user exists in the database
func (s *DatabaseService) UserExists(id string) bool {
	var count int64
	s.db.Model(&UserModel{}).Where("id = ?", id).Count(&count)
	return count > 0
}

// UserExistsByUsername checks if a user exists by username
func (s *DatabaseService) UserExistsByUsername(username string) bool {
	var count int64
	s.db.Model(&UserModel{}).Where("username = ?", username).Count(&count)
	return count > 0
}

// DeleteUser deletes a user from the database
func (s *DatabaseService) DeleteUser(id string) error {
	return s.db.Delete(&UserModel{}, "id = ?", id).Error
}

// FUD User related methods

// SaveFUDUser saves or updates a FUD user in the database
func (s *DatabaseService) SaveFUDUser(fudUser FUDUserModel) error {
	fudUser.UpdatedAt = time.Now()
	return s.db.Save(&fudUser).Error
}

// GetFUDUser retrieves a FUD user by user ID
func (s *DatabaseService) GetFUDUser(userID string) (*FUDUserModel, error) {
	var fudUser FUDUserModel
	err := s.db.Where("user_id = ?", userID).First(&fudUser).Error
	if err != nil {
		return nil, err
	}
	return &fudUser, nil
}

// IsFUDUser checks if a user is marked as FUD
func (s *DatabaseService) IsFUDUser(userID string) bool {
	var count int64
	s.db.Model(&FUDUserModel{}).Where("user_id = ?", userID).Count(&count)
	return count > 0
}

// GetAllFUDUsers retrieves all FUD users
func (s *DatabaseService) GetAllFUDUsers() ([]FUDUserModel, error) {
	var fudUsers []FUDUserModel
	err := s.db.Order("detected_at DESC").Find(&fudUsers).Error
	return fudUsers, err
}

// IncrementFUDUserMessageCount increments the message count for a FUD user
func (s *DatabaseService) IncrementFUDUserMessageCount(userID string, messageID string) error {
	return s.db.Model(&FUDUserModel{}).Where("user_id = ?", userID).Updates(map[string]interface{}{
		"message_count":   gorm.Expr("message_count + 1"),
		"last_message_id": messageID,
		"updated_at":      time.Now(),
	}).Error
}

// DeleteFUDUser deletes a FUD user from the database
func (s *DatabaseService) DeleteFUDUser(userID string) error {
	return s.db.Delete(&FUDUserModel{}, "user_id = ?", userID).Error
}

// Search and query methods

// SearchTweets searches tweets by text content
func (s *DatabaseService) SearchTweets(query string, limit int) ([]TweetModel, error) {
	var tweets []TweetModel
	err := s.db.Where("text LIKE ?", "%"+query+"%").Limit(limit).Order("created_at DESC").Find(&tweets).Error
	return tweets, err
}

// GetRecentTweets retrieves recent tweets with optional limit
func (s *DatabaseService) GetRecentTweets(limit int) ([]TweetModel, error) {
	var tweets []TweetModel
	err := s.db.Order("created_at DESC").Limit(limit).Find(&tweets).Error
	return tweets, err
}

// GetTweetCount returns the total number of tweets in the database
func (s *DatabaseService) GetTweetCount() (int64, error) {
	var count int64
	err := s.db.Model(&TweetModel{}).Count(&count).Error
	return count, err
}

// GetUserCount returns the total number of users in the database
func (s *DatabaseService) GetUserCount() (int64, error) {
	var count int64
	err := s.db.Model(&UserModel{}).Count(&count).Error
	return count, err
}

// GetFUDUserCount returns the total number of FUD users in the database
func (s *DatabaseService) GetFUDUserCount() (int64, error) {
	var count int64
	err := s.db.Model(&FUDUserModel{}).Count(&count).Error
	return count, err
}

// Close closes the database connection
func (s *DatabaseService) Close() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
