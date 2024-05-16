package routes

import (
	"example/remindme/controller"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func HandleRegistratingUser(client *mongo.Client, route *gin.Engine) {
	route.POST("/user/register", controller.RegisterUser(client))
}
