package models

import (
	"time"
	"github.com/google/uuid"
)

type Chat struct {
	Id        uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserId    uuid.UUID `json:"user_id" gorm:"type:uuid"`
	Title     string    `json:"title"`
	Summary   string    `json:"summary"`
	CreatedAt time.Time `json:"created_at"`
}
