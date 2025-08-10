package main

import (
	"log"

	"github.com/blog-platform/delivery/routers"
	"github.com/blog-platform/infrastructure"
	"github.com/blog-platform/repositories"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
        log.Fatal("No .env file found")
    }

	repositories.ConnectDB()
	infrastructure.ConnectClient()
	route := routers.Init(gin.Default())
	route.Run()
}