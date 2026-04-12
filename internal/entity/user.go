package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"index"`
	ID        uuid.UUID  `json:"ID" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Username  string     `json:"Username"`
	Email     string     `json:"Email"`
	Password  string     `json:"Password"`
	Role      string     `json:"Role"`
	Verified  bool       `json:"Verified"`
}
