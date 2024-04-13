package controllers

import (
	"context"
	"fmt"
	"math"
	"resturant-backend/database"
	"resturant-backend/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")
 var validate = validator.New()

func GetFoods() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)  

		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))

		if err!= nil || recordPerPage <1 {
          recordPerPage = 10;
		}

		page, err := strconv.Atoi(c.Query("page"));

		if err!= nil  || page < 1 {
			page = 1;
		}
		startIndex := (page-1) *recordPerPage;
startIndex, err = strconv.Atoi(c.Query("startIndex"))

		matchStage := bson.D{{"$match", bson.D{{}}}}
		groupStage := bson.D{{"$group", bson.D{
		    {"_id",bson.D{{"_id","null"}}}, 
			{"total_count", bson.D{{"$sum", 1}}}, 
			{"data", bson.D{{"$push", "$$ROOT"}}}
		}}}
		projectStage := bson.D{
			{"$project", bson.D{
				{"_id", 0},
				{"total_count", 1},
				{"foot_items",bson.D{{"$slice", []interface{}{"$data", startIndex, recordPerPage}}}, },
		}}}

		result, err := foodCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage, groupStage, projectStage,
		})

		defer cancel()

		// check err!\
		if err!= nil {
			c.JSON(504, gin.H{"error": err.Error()})
			
		}

		var allFoods []bson.M;
		if err = result.All(ctx, &allFoods); err!=nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allFoods[0])


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
		err := menuCollection.FindOne(ctx, bson.M{"menu_id": food.Menu_id}).Decode(&menu)
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
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var food models.Food
		var menu models.Menu

		foodId := c.Param("food_id")

		if err := c.BindJSON(&menu); err != nil {
			c.JSON(504, gin.H{"error": err.Error()})
			return
		}

		var updateObj primitive.D;

		if food.Name != nil {
updateObj = append(updateObj, bson.E{"name", name.Name})
		}

		if food.Price != nil{
			updateObj = append(updateObj, bson.E{"price", name.Price})

		}

		if food.Food_Image != nil {
			updateObj = append(updateObj, bson.E{"food_image", name.Food_Image})

		}

		if food.Menu_id != nil {
			err := menuColection.FindOne(ctx, bson.M{"menu_id": food.Menu_id}).Decode(&menu)

			defer cancel();

			if err!= nil {
				msg := fmt.Sprintf("message: menu was not found");
				c.JSOn(504, gin.H{"error":msg})
				return;
			}
			updateObj = append(updateObj, bson.E{"menu", food.Price})

		}

		menu.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{"updated_at", food.Updated_at})

		upsert := true;
		filter := bson.M{"food_id":foodId}

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

		if err != nil {
			msg := fmt.Sprintf("food item was not updated")
			c.JSON(504, gin.H{"error": msg})
			return
		}

		defer cancel()

		c.JSON(200, result)

	}
}

func round(num float64) int {
  return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
output := math.Pow(10, float64(precision));
return float64(round(num*output))/output
}
