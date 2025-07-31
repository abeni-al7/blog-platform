package main

import (
	"log"

	"github.com/abeni-al7/blog-platform/delivery/routers"
	"github.com/abeni-al7/blog-platform/infrastructure"
	"github.com/abeni-al7/blog-platform/repositories"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
        log.Fatal("No .env file found")
    }

	repositories.ConnectDB()
	infrastructure.ConnectClient()
	routers.Init()
}