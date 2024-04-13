package controllers

import (
	"context"
	"log"
	"net/http"
	"resturant-backend/database"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)


var orderCollection *mongo.Collection = database.OpenCollection(database.Client, "order")

func GetOrders() gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)  

		result, err := orderCollection.Find(context.TODO(), bson.M{})

		defer cancel();

		if err != nil {
			c.JSON(504, gin.H{"error": err.Error()})
		}

		var allOrders []bson.M;
		if err = result.All(ctx, &allOrders); err!=nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allOrders)
	}
}

func GetOrder() gin.HandlerFunc{
	return func(c *gin.Context){
		
	}
}

func CreateOrder() gin.HandlerFunc{
	return func(c *gin.Context){
		
	}
}

func UpdateOrder() gin.HandlerFunc{
	return func(c *gin.Context){
		
	}
}
