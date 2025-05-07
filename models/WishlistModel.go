package models

import (
	"time"
	"github.com/google/uuid"
)

type Wishlist struct {
	Id         int64 `gorm:"-" json:"id"`
	UserId     uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	ProductId    string `gorm:"type:varchar(255)" json:"product_id"`;
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
}