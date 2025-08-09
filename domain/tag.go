package domain

import (
	"time"
	// "gorm.io/gorm"
)

type Tag struct {
	// gorm.Model
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"type:varchar(100);uniqueIndex" json:"name"` // Unique tag name
	Content   string    `gorm:"type:varchar(500)" json:"content"`
	CreatedAt time.Time `json:"created_at"` // auto set on insert
	UpdatedAt time.Time `json:"updated_at"` // auto set on update
}

type Tag_Blog struct {
	// gorm.Model
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	BlogID    int64     `json:"blog_id"`                                        // Foreign key column
	Blog      Blog      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"` // GORM relation
	TagID     int64     `json:"tag_id"`                                         // Foreign key column
	Tag       Tag       `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"` // GORM relation
	CreatedAt time.Time `json:"created_at"`                                     // auto set on insert
	UpdatedAt time.Time `json:"updated_at"`                                     // auto set on update
}
