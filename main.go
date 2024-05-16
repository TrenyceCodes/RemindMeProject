package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()

	server.GET("/", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{"message": "hello world"})
	})

	if err := server.Run(":3001"); err != nil {
		fmt.Println("Server is experiencing problems running")
	}
}
