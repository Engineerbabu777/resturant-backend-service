package controllers

import (
	"context"
	"log"
	"net/http"
	"resturant-backend/database"
	"resturant-backend/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type OrderItemPack struct {
	Table_id    *string
	Order_items []models.OrderItem
}

var orderItemCollection *mongo.Collection = database.OpenCollection(database.Client, "orderItem")

func GetOrderItems() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)  

		result, err := orderItemCollection.Find(context.TODO(), bson.M{})

		defer cancel();

		if err != nil {
			c.JSON(504, gin.H{"error": err.Error()})
		}

		var allOrdersItems []bson.M;
		if err = result.All(ctx, &allOrdersItems); err!=nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allOrdersItems)
	}
}

func GetOrderItemByOrder() gin.HandlerFunc {
	return func(c *gin.Context) {

		orderId := c.Param("order_id");

		allOrdersItems, err := ItemsByOrder(orderId);

		if err!= nil {
			c.JSON(504, gin.H{"error": err.Error()})
			return;
		}

		c.JSON(http.StatusOK, allOrdersItems)
	}
}

func UpdateOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func ItemsByOrder(id string) (OrderItems []primitive.M, err error) {

}

func GetOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func CreateOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}