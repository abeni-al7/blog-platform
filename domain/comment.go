package domain

import (
	"time"
	//"gorm.io/gorm"
)

type Comment struct {
	//gorm.Model
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Content   string    `json:"content"`
	UserID    int64     `json:"user_id"`                                        // Foreign key column
	User      User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"` // GORM relation
	BlogID    int64     `json:"blog_id"`                                        // Foreign key column
	Blog      Blog      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"` // GORM relation
	CreatedAt time.Time `json:"created_at"`                                     // auto set on insert
	UpdatedAt time.Time `json:"updated_at"`                                     // auto set on update
}
