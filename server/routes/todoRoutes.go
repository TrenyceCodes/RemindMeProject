package routes

import (
	"example/remindme/controller"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func HandleTodoRoutes(client *mongo.Client, route *gin.Engine) {
	route.POST("/todo", controller.CreateTodo(client))
	route.GET("/todo/:id", controller.GetTodoById(client))
}
