package controller

import (
	"context"
	"example/remindme/model"
	"example/remindme/utils"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterUser(client *mongo.Client) gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		//reads in env file
		envFile, err := godotenv.Read(".env")
		if err != nil {
			ginContext.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to load environment variables"})
			return
		}

		//loads in user model, setting user id and id to the same object value
		var user model.User

		if client == nil {
			ginContext.JSON(http.StatusInternalServerError, gin.H{"message": "Client is nil"})
			return
		}

		//handles setting database connection and inserting into database
		mongoContext := context.Background()
		mongoCollection := client.Database(envFile["DATABASE_NAME"]).Collection(envFile["DATABASE_COLLECTION"])

		if err := ginContext.ShouldBindJSON(&user); err != nil {
			ginContext.JSON(http.StatusInternalServerError, gin.H{"message": "ShouldBindJSON user data failed"})
			return
		}

		jsonwebtoken, err := utils.CreateToken(user.Username, user.User_id)
		if err != nil {
			log.Fatal("Error: ", err)
			return
		}

		user = model.User{
			Username:     user.Username,
			Email:        user.Email,
			Password:     user.Password,
			JsonWebToken: jsonwebtoken,
		}

		hashPassword, message, err := utils.HashPassword(user.Password)
		if err != nil {
			log.Fatal("Error: ", err)
			return
		}

		user.Id = primitive.NewObjectID()
		user.User_id = user.Id.Hex()
		user.Password = hashPassword
		user.JsonWebToken = jsonwebtoken
		fmt.Println(message)

		_, err = mongoCollection.InsertOne(mongoContext, user)
		if err != nil {
			ginContext.JSON(http.StatusInternalServerError, gin.H{"message": "Error inserting user"})
			return
		}

		ginContext.JSON(http.StatusOK, gin.H{
			"message": "user created successfully",
			"data": map[string]interface{}{
				"username":     user.Username,
				"email":        user.Email,
				"jsonwebtoken": user.JsonWebToken,
			},
		})
	}
}

func LoginUser(client *mongo.Client) gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		envFile, err := godotenv.Read(".env")
		if err != nil {
			ginContext.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to load environment variables"})
			return
		}

		var user model.User
		var foundUser model.User

		if client == nil {
			ginContext.JSON(http.StatusInternalServerError, gin.H{"message": "Client is nil"})
			return
		}

		mongoContext := context.Background()
		mongoCollection := client.Database(envFile["DATABASE_NAME"]).Collection(envFile["DATABASE_COLLECTION"])

		if err := ginContext.ShouldBindJSON(&user); err != nil {
			ginContext.JSON(http.StatusInternalServerError, gin.H{"message": "ShouldBindJSON user data failed"})
			return
		}

		filter := bson.M{
			"email":    user.Email,
			"username": user.Username,
		}

		if err := mongoCollection.FindOne(mongoContext, filter).Decode(&foundUser); err != nil {
			ginContext.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		isPasswordValid, message := utils.ValidatePassword(user.Password, foundUser.Password)
		if !isPasswordValid {
			fmt.Println(message)
			return
		}

		ginContext.JSON(http.StatusOK, gin.H{
			"message": "User logged in successfully",
			"data": map[string]interface{}{
				"username": user.Username,
				"email":    user.Email,
			},
		})
	}
}
