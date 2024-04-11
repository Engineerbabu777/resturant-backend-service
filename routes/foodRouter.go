package routes

import "github.com/gin-gonic/gin"



func foodRoutes(incomingRoutes *gin.Engine){
	incomingRoutes.GET("/foods", controller.GetFoods());
	incomingRoutes.GET("/users/:food_id", controller.GetFood());
	incomingRoutes.PUT("/foods", controller.CreateFood());
	incomingRoutes.PATCH("/users/:food_id", controller.UpdateFood());
}