package main

import (
	"gorm.io/gorm"
	"time"
)

// Tweet model for database storage
type TweetModel struct {
	gorm.Model
	ID          string    `gorm:"primaryKey;column:id" json:"id"`
	Text        string    `gorm:"column:text" json:"text"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
	ReplyCount  int       `gorm:"column:reply_count" json:"reply_count"`
	UserID      string    `gorm:"column:user_id;index" json:"user_id"`
	InReplyToID string    `gorm:"column:in_reply_to_id;index" json:"in_reply_to_id,omitempty"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (TweetModel) TableName() string {
	return "tweets"
}

// User model for database storage
type UserModel struct {
	gorm.Model
	ID        string    `gorm:"primaryKey;column:id" json:"id"`
	Username  string    `gorm:"column:username;uniqueIndex" json:"username"`
	Name      string    `gorm:"column:name" json:"name"`
	IsFUD     bool      `gorm:"column:is_fud;default:false" json:"is_fud"`
	FUDType   string    `gorm:"column:fud_type" json:"fud_type,omitempty"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (UserModel) TableName() string {
	return "users"
}

// FUDUser model for storing FUD users information
type FUDUserModel struct {
	gorm.Model
	ID             uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID         string    `gorm:"column:user_id;uniqueIndex" json:"user_id"`
	Username       string    `gorm:"column:username" json:"username"`
	FUDType        string    `gorm:"column:fud_type" json:"fud_type"`
	FUDProbability float64   `gorm:"column:fud_probability" json:"fud_probability"`
	DetectedAt     time.Time `gorm:"column:detected_at" json:"detected_at"`
	MessageCount   int       `gorm:"column:message_count;default:1" json:"message_count"`
	LastMessageID  string    `gorm:"column:last_message_id" json:"last_message_id"`
	CreatedAt      time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (FUDUserModel) TableName() string {
	return "fud_users"
}
