package models

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
	"time"
)

type Chat struct {
	Id        uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	UserId    uuid.UUID      `gorm:"type:uuid;not null" json:"user_id"`
	Title     string         `gorm:"type:varchar(255)" json:"title"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	Message   string         `gorm:"-" json:"message" validate:"required"`
	ImageURLs pq.StringArray `gorm:"-" json:"image_urls"`
}
