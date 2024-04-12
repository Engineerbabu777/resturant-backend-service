package controllers

import (
	"context"
	"fmt"
	"resturant-backend/database"
	"resturant-backend/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")

func GetFoods() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func GetFood() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		foodId := c.Param("food_id")
		var food models.Food

		err := foodCollection.FindOne(ctx, bson.M{"food_id": foodId}).Decode(&food)
		defer cancel()
		if err != nil {
			c.JSON(504, gin.H{"error": err.Error()})
		}
		c.JSON(200, food)
	}
}

func CreateFood() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var food models.Food
		var menu models.Menu

		if err := c.BindJSON(&food); err!=nil {
			c.JSON(504, gin.H{"error": err.Error()})
			return;
		}

		validationErr := validate.Struct(food);
		if validationErr != nil {
			c.JSON(504, gin.H{"error": validationErr.Error()})
			return;
		}
		err := menuCollection.FindOne(ctx, bson.M{"menu_id": food.MenuId}).Decode(&menu)
		defer cancel()
		if err != nil {
		    msg := fmt.Sprintf("menu was not food");
			c.JSOn(504, gin.H{"error":msg})
			return;
		}
		food.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		food.Updated_at, _ =time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		food.ID = primitive.NewObjectID()
		food.Food_id = food.ID.Hex()
		var num = toFixed(*food.Price, 2)
		food.Price = &num
		result, insertErr := foodCollection.InsertOne(ctx, food)
		if insertErr != nil {
		    msg := fmt.Sprintf("food item was not created");
			c.JSOn(504, gin.H{"error":msg})
			return;
		}
		c.JSON(200, result)
	}
}

func UpdateFood() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func round(num float64) int {

}

func toFixed(num float64, precision int) float64 {

}
