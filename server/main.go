package main

import (
	"example/remindme/connection"
	"example/remindme/routes"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
)

func main() {
	StartGinServer()
}

// starts up the main server
func StartGinServer() {
	gin.SetMode(gin.ReleaseMode)

	server := gin.New()
	config := cors.DefaultConfig()
	config.AllowCredentials = true
	config.AllowAllOrigins = false
	config.AllowOriginFunc = func(origin string) bool {
		return true
	}

	client, message, err := connection.MongoDatabaseConnection()
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	fmt.Println(message)

	routes.HandleUserAuthentication(client, server)
	routes.HandleTodoRoutes(client, server)
	handleServerMiddleware(server, config)

	serverPort := os.Getenv("PORT")
	if serverPort == "" {
		serverPort = "8080"
	}

	if err := server.Run("localhost:" + serverPort); err != nil {
		fmt.Println("Server is experiencing problems running")
	}
}

// handles middleware of remindme
func handleServerMiddleware(server *gin.Engine, config cors.Config) {
	server.Use(func(context *gin.Context) {
		// Set CORS headers
		context.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		context.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		context.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		context.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")

		// Handle preflight request
		if context.Request.Method == "OPTIONS" {
			context.AbortWithStatus(http.StatusNoContent)
			return
		}

		// Get the JWT from the Authorization header
		authHeader := context.Writer.Header().Get("Authorization")
		if authHeader == "" {
			context.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization header is required"})
			context.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Load environment variables
		err := godotenv.Load(".env")
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to load environment variables"})
			context.Abort()
			return
		}

		secretKey := os.Getenv("SECRET_KEY")
		if secretKey == "" {
			context.JSON(http.StatusInternalServerError, gin.H{"message": "SECRET_KEY is not set in the environment variables"})
			context.Abort()
			return
		}

		// Parse and validate the JWT
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			context.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
			context.Abort()
			return
		}

		// Extract claims and set context
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			username, ok := claims["username"].(string)
			if !ok {
				context.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token claims"})
				context.Abort()
				return
			}
			context.Set("Username", username)
			if userId, ok := claims["userId"].(string); ok {
				context.Set("UserId", userId)
			}
		} else {
			context.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token claims"})
			context.Abort()
			return
		}

		context.Next()
	})

	server.Use(cors.New(config))
}
