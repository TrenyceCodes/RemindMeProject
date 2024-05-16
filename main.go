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
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	serverPort := os.Getenv("PORT")

	server := gin.Default()

	server.GET("/", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{"message": "hello world"})
	})

	if err := server.Run(":" + serverPort); err != nil {
		fmt.Println("Server is experiencing problems running")
	}
}
