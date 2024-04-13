package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"resturant-backend/database"
	"resturant-backend/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


var orderCollection *mongo.Collection = database.OpenCollection(database.Client, "order")
var tableCollection *mongo.Collection = database.OpenCollection(database.Client, "table")

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
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)  
		orderId := c.Param("order_id")
		var order models.Order

		err := foodCollection.FindOne(ctx, bson.M{"order_id": orderId}).Decode(&order)
		defer cancel()
		if err != nil {
			c.JSON(504, gin.H{"error": err.Error()})
		}
		c.JSON(200, order)
	}
}

func CreateOrder() gin.HandlerFunc{
	return func(c *gin.Context){
		var table models.Table;
		var order models.Order;
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		if err := c.BindJSON(&order); err!=nil {
			c.JSON(504, gin.H{"error": err.Error()})
			return;
		}

		validationErr := validate.Struct(order);
		if validationErr != nil {
			c.JSON(504, gin.H{"error": validationErr.Error()})
			return;
		}

		if order.Table_id !=nil{
			err := tableCollection.FindOne(ctx, bson.M{"table_id":order.Table_id}).Decode(&table);
			defer cancel();

			if err!= nil {
				msg := fmt.Sprintf("message table was not found!");
				c.JSON(504, gin.H{"error": msg})
				return;
			}
		}

		order.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		order.Updated_at, _ =time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		order.ID = primitive.NewObjectID()
		order.Order_id = order.ID.Hex()
		result, insertErr := foodCollection.InsertOne(ctx, order)
		if insertErr != nil {
		    msg := fmt.Sprintf("order item was not created");
			c.JSOn(504, gin.H{"error":msg})
			return;
		}
		defer cancel()
		c.JSON(200, result)

	}
}

func UpdateOrder() gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var table models.Table
		var order models.Order
		var food models.Food


		var updateObj primitive.D;
		orderId := c.Param("order_id")

		if err := c.BindJSON(&order); err != nil {
			c.JSON(504, gin.H{"error": err.Error()})
			return
		}

		if order.Table_id != nil{
			err := menuCollection.FindOne(ctx, bson.M{"table_id":food.Table_id}).Decode(&table);
			defer cancel();

			if err!= nil {
				msg := fmt.Sprintf("message menu was not found!");
				c.JSON(504, gin.H{"error": msg})
				return;
			}
			updateObj = append(updateObj, bson.E{"menu", order.Table_id})
		}

		order.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{"updated_at", food.Updated_at})

		upsert := true;
		filter := bson.M{"order_id":orderId}

		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		result, err := foodCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{"$set", updateObj},
			},
			&opt,
		)

		if err!= nil {
			msg :=  fmt.Sprintf("order was not updated!");
			c.JSON(504, gin.H{"error": msg})
			return;
		}

		defer cancel();
		c.JSON(200, result)

	}
}


func OrderItemOrderCreator(order models.Order) string{
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	
	order.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		order.Updated_at, _ =time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		order.ID = primitive.NewObjectID();
		order.Order_id = order.ID.Hex();

		orderCollection.InsertOne(ctx, order);
		defer cancel();

		return order.Order_id


}