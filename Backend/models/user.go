package models

import (
	"time"
)

type User struct {
	UserID       string    `gorm:"primaryKey;column:userid"`
	Username     string    `gorm:"unique;column:username"`
	ProfileImage []byte    `gorm:"column:profile_image"`
	Email        string    `gorm:"unique;column:email"`
	Password     string    `gorm:"column:password"`
	Role         string    `gorm:"column:role"`
	EntityID     string    `gorm:"column:entityid"`
	Telephone    string    `gorm:"column:telephone"`
	CreatedAt    time.Time `gorm:"column:created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`
	DeletedAt    time.Time `gorm:"column:deleted_at"`
}
