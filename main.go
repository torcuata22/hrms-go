package main

import (
	"context"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson" //binary json
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
	mg.DB = db
	mg.client = client
	log.Println("Connected to Mongo!")
	return nil
}
func GetEmployees(c *fiber.Ctx) error {
	var employees []Employee
	collection := mg.DB.Collection("employees")
	filter := bson.D{}
	cur, err := collection.Find(context.Background(), filter)
	if err != nil {
		return err
	}
	err = cur.All(context.Background(), &employees)
	if err != nil {
		return err
	}
	return c.JSON(employees)
}

func main() {
	if err := Connect(); err != nil {
		log.Fatal(err)
	}
	app := fiber.New()
	app.Get("/employees", GetEmployees)
	// app.Get("/employee/:id", GetEmployee)
	// app.Post("/employees", PostEmployee)
	// app.Put("/employee/:id", PutEmployee)
	// app.Delete("/employee/:id", DeleteEmployee)
}
