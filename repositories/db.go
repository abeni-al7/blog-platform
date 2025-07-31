package repositories

import (
	"fmt"
	"log"
	"os"

	"github.com/abeni-al7/blog-platform/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))

    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    DB = db

	err = DB.AutoMigrate(&domain.User{}, &domain.Blog{}, &domain.Comment{}, &domain.Tag{}, &domain.Tag_Blog{}, &domain.Token{})
    if err != nil {
        log.Fatal("Failed to migrate database:", err)
    }
}