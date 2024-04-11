package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/labstack/echo/middleware"
)

var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")
func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
		fmt.Println("INFO: No PORT environment variable detected, defaulting to " + port)
	}

	router := gin.New()

	router.Use(gin.Logger())

	 routes.UserRouter(router);
	 routes.Use(middleware.Authenticate());

	 routes.FoodRoutes(router);
	 routes.MenuRoutes(router);
	 routes.TableRoutes(router);
	 routes.OrderRoutes(router);
	 routes.OrderItemRoutes(router);
	 routes.InvoiceRoutes(router);


	 router.Run(":"+port)

}
