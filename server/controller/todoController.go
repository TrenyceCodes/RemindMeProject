package controller

import (
	"context"
	"example/remindme/model"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateTodo(client *mongo.Client) gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		// Load environment variables
		err := godotenv.Load(".env")
		if err != nil {
			ginContext.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to load environment variables"})
			return
		}

		// Get the JWT from the request header
		tokenString := ginContext.Request.Header.Get("Authorization")
		if tokenString == "" {
			ginContext.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization token required"})
			return
		}

		// Verify and decode the JWT
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Ensure the signing method is what you expect
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			// Return the secret key
			return []byte(os.Getenv("SECRET_KEY")), nil
		})

		if err != nil {
			ginContext.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid authorization token"})
			return
		}

		// Extract username from token claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			ginContext.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token claims"})
			return
		}
		username, ok := claims["username"].(string)
		if !ok {
			ginContext.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token claims"})
			return
		}

		// Bind the request payload to the Todo model
		var todo model.Todo
		if err := ginContext.ShouldBindJSON(&todo); err != nil {
			ginContext.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request payload"})
			return
		}

		// Prepare the MongoDB context and collection
		mongoContext, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		mongoCollection := client.Database(os.Getenv("DATABASE_NAME")).Collection("todos")

		// Set the todo fields
		todo.Id = primitive.NewObjectID()
		todo.Author = username
		todo.Todo_id = todo.Id.Hex()

		// Insert the todo into the database
		_, err = mongoCollection.InsertOne(mongoContext, todo)
		if err != nil {
			ginContext.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create todo"})
			return
		}

		// Respond with success
		ginContext.JSON(http.StatusOK, gin.H{"message": "Todo created successfully", "todo": todo})
	}
}
