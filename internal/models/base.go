package models

import (
	"time"

	"gorm.io/gorm"
)

// BaseEntity holds the columns every persisted model shares.
// Embed it anonymously (like extending a NestJS BaseEntity) so that both GORM
// and encoding/json flatten its fields into the parent struct.
type BaseEntity struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
