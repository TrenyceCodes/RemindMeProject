package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
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

	serverPort := os.Getenv("PORT")
	if serverPort == "" {
		serverPort = "8080"
	}

	if err := server.Run(":" + serverPort); err != nil {
		fmt.Println("Server is experiencing problems running")
	}
}
