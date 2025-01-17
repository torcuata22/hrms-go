package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//Mongo DB struct

type MongoInstance struct {
	client *mongo.Client
	DB     *mongo.Database
}

type Employee struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Salary float64 `json:"salary"`
	Age    int     `json:"age"`
}

var mg MongoInstance

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	mongoURI := os.Getenv("MONGO_URI")
	dbName := os.Getenv("MONGO_DB")

	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Connect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	db := client.Database(dbName)
	mg.DB = db
	mg.client = client
	log.Println("Connected to Mongo!")

}
