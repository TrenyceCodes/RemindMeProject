package main

import (
	"example/remindme/connection"
	"example/remindme/routes"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	StartGinServer()
}

func StartGinServer() {
	gin.SetMode(gin.ReleaseMode)

	server := gin.New()
	config := cors.DefaultConfig()
	config.AllowCredentials = true
	config.AllowAllOrigins = false
	// I think you should whitelist a limited origins instead:
	//  config.AllowAllOrigins = []{"xxxx", "xxxx"}
	config.AllowOriginFunc = func(origin string) bool {
		return true
	}

	server.Use(func(context *gin.Context) {
		host := context.Request.Header.Get("Origin")
		context.Writer.Header().Set("Access-Control-Allow-Origin", host)
		context.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		context.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		context.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		if context.Request.Method == "OPTIONS" {
			log.Println("Handling OPTIONS request")
			context.AbortWithStatus(http.StatusNoContent)
			return
		}
		log.Println("Executing CORS middleware")
		context.Next()
	})
	server.Use(cors.New(config))

	client, message, err := connection.MongoDatabaseConnection()
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	fmt.Println(message)

	server.GET("/", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{"message": "hello world"})
	})

	routes.HandleRegistratingUser(client, server)

	serverPort := os.Getenv("PORT")
	if serverPort == "" {
		serverPort = "8080"
	}

	if err := server.Run(":" + serverPort); err != nil {
		fmt.Println("Server is experiencing problems running")
	}
}
