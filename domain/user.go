package domain

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID             int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Username       string    `gorm:"type:varchar(255)" json:"username"`
	Email          string    `gorm:"type:varchar(500)" json:"email"`
	Password       string    `gorm:"type:varchar(255)" json:"_"`
	Role           string    `gorm:"type:varchar(255)" json:"role"`
	Bio            string    `json:"bio"`
	ProfilePicture string    `gorm:"type:varchar(500)" json:"profile_picture"`
	Phone          string    `gorm:"type:varchar(255)" json:"phone"`
	Status         string    `gorm:"type:varchar(255)" json:"status"`
	CreatedAt      time.Time `json:"created_at"` // auto set on insert
	UpdatedAt      time.Time `json:"updated_at"` // auto set on update
}
