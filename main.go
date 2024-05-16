package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	StartGinServer()
}

func StartGinServer() {
	gin.SetMode(gin.ReleaseMode)

	server := gin.New()

	server.GET("/", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{"message": "hello world"})
	})

	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	serverPort := os.Getenv("PORT")

	if err := server.Run(":" + serverPort); err != nil {
		fmt.Println("Server is experiencing problems running")
	}
}
