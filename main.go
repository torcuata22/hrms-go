package main

import (
	"context"
	"log"
	"os"

	//"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson" //binary json

	//"go.mongodb.org/mongo-driver/bson/primitive"
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
	filter := bson.D{} //slice representation of binary json, in this case empty
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

func GetEmployee(c *fiber.Ctx) error {
	id := c.Params("id")
	// objID, err := primitive.ObjectIDFromHex(id)
	// if err != nil {
	// 	return err
	// }
	var employee Employee
	collection := mg.DB.Collection("employees")
	filter := bson.M{"id": id} // fix: use ObjectId instead of string
	err := collection.FindOne(context.Background(), filter).Decode(&employee)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fiber.ErrNotFound
		}
		return err
	}
	return c.JSON(employee)
}

func CreateEmployee(c *fiber.Ctx) error {
	var employee Employee
	if err := c.BodyParser(&employee); err != nil {
		return err
	}

	collection := mg.DB.Collection("employees")
	_, err := collection.InsertOne(context.Background(), employee)
	if err != nil {
		return err
	}
	return c.JSON(employee)
}

func DeleteEmployee(c *fiber.Ctx) error {
	id := c.Params("id")
	collection := mg.DB.Collection("employees")
	filter := bson.M{"id": id}
	result, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return fiber.ErrNotFound
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func UpdateEmployee(c *fiber.Ctx) error {
	id := c.Params("id")
	collection := mg.DB.Collection("employees")
	filter := bson.M{"id": id}

	var employee Employee
	if err := c.BodyParser(&employee); err != nil {
		return err
	}

	log.Printf("Received employee data: %+v", employee)

	update := bson.M{
		"$set": bson.M{
			"name":   employee.Name,
			"salary": employee.Salary,
			"age":    employee.Age,
		},
	}
	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return fiber.ErrNotFound
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func main() {
	if err := Connect(); err != nil {
		log.Fatal(err)
	}
	app := fiber.New()
	app.Get("/employees", GetEmployees)
	app.Get("/employee/:id", GetEmployee)
	app.Post("/employees", CreateEmployee)
	app.Put("/employee/:id", UpdateEmployee)
	app.Delete("/employee/:id", DeleteEmployee)

	log.Fatal(app.Listen(":3000"))
}
