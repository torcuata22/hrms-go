package database

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv" //binary json
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//Mongo DB struct

type MongoInstance struct {
	client *mongo.Client
	DB     *mongo.Database
}

var Mg MongoInstance

func Connect() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}

	mongoURI := os.Getenv("MONGO_URI")
	dbName := os.Getenv("MONGO_DB")

	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		return err
	}

	err = client.Connect(context.TODO())
	if err != nil {
		return err
	}
	db := client.Database(dbName)
	Mg.DB = db
	Mg.client = client
	log.Println("Connected to Mongo!")
	return nil
}
