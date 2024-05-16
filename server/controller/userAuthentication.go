package controller

import (
	"context"
	"example/remindme/model"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterUser(client *mongo.Client) gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		envFile, _ := godotenv.Read(".env")
		var user model.User

		if client == nil {
			ginContext.JSON(http.StatusInternalServerError, gin.H{"message": "Client is nil"})
		}

		mongoContext := context.Background()
		mongoCollection := client.Database(envFile["DATABASE_NAME"]).Collection(envFile["DATABASE_COLLECTION"])

		user = model.User{
			Username: user.Username,
			Email:    user.Email,
			Password: user.Password,
		}

		if err := ginContext.ShouldBindJSON(&user); err != nil {
			ginContext.JSON(http.StatusInternalServerError, gin.H{"message": "ShouldBindJSON user data failed"})
		}

		insertUserData, err := mongoCollection.InsertOne(mongoContext, user)
		if err != nil {
			ginContext.JSON(http.StatusInternalServerError, gin.H{"message": "Error inserting user"})
		}

		fmt.Println("User inserted: ", insertUserData.InsertedID)
		ginContext.JSON(http.StatusOK, gin.H{
			"message": "user created successfully",
			"data": map[string]interface{}{
				"inserted_id": insertUserData.InsertedID,
				"username":    user.Username,
				"email":       user.Email,
			},
		})
	}
}
