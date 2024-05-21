package controller

import (
	"context"
	"example/remindme/model"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
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

		tokenString := ginContext.Request.Header.Get("Authorization")
		if tokenString == "" {
			ginContext.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization token not provided"})
			return
		}

		// Remove "Bearer " prefix if present
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		secretKey := os.Getenv("SECRET_KEY")
		if secretKey == "" {
			ginContext.JSON(http.StatusInternalServerError, gin.H{"message": "SECRET_KEY is not set in the environment variables"})
			ginContext.Abort()
			return
		}

		// Parse and validate the JWT
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secretKey), nil
		})

		if err != nil {
			ginContext.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid authorization token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			ginContext.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
			return
		}

		username, ok := claims["username"].(string)
		if !ok {
			ginContext.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token claims"})
			return
		}

		fmt.Println("Username: ", username)

		// Bind the request payload to the Todo model
		var todo model.Todo
		if err := ginContext.ShouldBindJSON(&todo); err != nil {
			ginContext.JSON(http.StatusBadRequest, gin.H{"message": "There was an issue binding the model"})
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

func GetTodo(client *mongo.Client) gin.HandlerFunc {
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

		// Remove "Bearer " prefix if present
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
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

		// Prepare the MongoDB context and collection
		mongoContext, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		mongoCollection := client.Database(os.Getenv("DATABASE_NAME")).Collection("todos")

		// Define the filter to search by author
		filter := bson.M{"author": username}

		// Find the todos in the database
		cursor, err := mongoCollection.Find(mongoContext, filter)
		if err != nil {
			ginContext.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve todos"})
			return
		}
		defer cursor.Close(mongoContext)

		// Iterate through the cursor and decode each todo
		var todos []model.Todo
		for cursor.Next(mongoContext) {
			var todo model.Todo
			if err := cursor.Decode(&todo); err != nil {
				ginContext.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to decode todo"})
				return
			}
			todos = append(todos, todo)
		}

		// Check for any cursor errors
		if err := cursor.Err(); err != nil {
			ginContext.JSON(http.StatusInternalServerError, gin.H{"message": "Cursor error"})
			return
		}

		// Respond with the retrieved todos
		ginContext.JSON(http.StatusOK, gin.H{"message": "Todos retrieved successfully", "todos": todos})
	}
}

func UpdateTodo(client *mongo.Client) gin.HandlerFunc {
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

		// Remove "Bearer " prefix if present
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
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

		// Define the filter to find the todo by id and author
		filter := bson.M{"_id": todo.Id, "author": username}

		// Update the todo in the database
		update := bson.M{
			"$set": bson.M{
				"title":       todo.Title,
				"description": todo.Description,
			},
		}
		_, err = mongoCollection.UpdateOne(mongoContext, filter, update)
		if err != nil {
			ginContext.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update todo"})
			return
		}

		// Respond with success
		ginContext.JSON(http.StatusOK, gin.H{"message": "Todo updated successfully", "todo": todo})
	}
}

func DeleteTodo(client *mongo.Client) gin.HandlerFunc {
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

		// Remove "Bearer " prefix if present
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
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

		// Get the todo ID from the URL parameter
		todoId := ginContext.Param("id")
		objectId, err := primitive.ObjectIDFromHex(todoId)
		if err != nil {
			ginContext.JSON(http.StatusBadRequest, gin.H{"message": "Invalid todo ID"})
			return
		}

		// Prepare the MongoDB context and collection
		mongoContext, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		mongoCollection := client.Database(os.Getenv("DATABASE_NAME")).Collection("todos")

		// Define the filter to find the todo by id and author
		filter := bson.M{"_id": objectId, "author": username}

		// Delete the todo from the database
		_, err = mongoCollection.DeleteOne(mongoContext, filter)
		if err != nil {
			ginContext.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete todo"})
			return
		}

		// Respond with success
		ginContext.JSON(http.StatusOK, gin.H{"message": "Todo deleted successfully"})
	}
}
