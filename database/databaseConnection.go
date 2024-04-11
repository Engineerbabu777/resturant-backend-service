package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBInit()  *mongo.Client{
	mongoDB := "string...";
	fmt.Println(mongoDB);

	client, err := mongo.NewClient(options.Client().ApplyURI(mongoDB));

	if err != nil {
		log.Fatal(err);
	}

	var ctx,cancel = context.WithTimeout(context.Background(), 10*time.Second);

	defer cancel();

	err = client.Connect(ctx);

	if err != nil {
		log.Fatal(err);
	}

	fmt.Println("Connected to database")

	return client;
}


var Client *mongo.Client = DBInit();

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
    return client.Database("resturant-db").Collection(collectionName)
}