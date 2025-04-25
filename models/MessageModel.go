package models

import (
	"github.com/lib/pq"
	"time"
	"github.com/google/uuid"
)

type Message struct {
	Id         uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	ChatId     uuid.UUID `gorm:"type:uuid;not null" json:"chat_id"`
	SenderRole string    `gorm:"type:varchar(10);not null" json:"sender_role"` 
	Text       string    `gorm:"type:text;not null" json:"text"`
	ImageURLs pq.StringArray `gorm:"type:text[]" json:"image_urls"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
}