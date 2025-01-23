package controller

import (
	"context"
	"log"

	"github.com/torcuata22/hrms-mongo/database"
	model "github.com/torcuata22/hrms-mongo/models"

	//"os"

	//"strconv"

	"github.com/gofiber/fiber/v2"
	//"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson" //binary json

	//"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	//"go.mongodb.org/mongo-driver/mongo/options"
)

func GetEmployees(c *fiber.Ctx) error {
	var employees []model.Employee
	collection := database.Mg.DB.Collection("employees")
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

	var employee model.Employee
	collection := database.Mg.DB.Collection("employees")
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
	var employee model.Employee
	if err := c.BodyParser(&employee); err != nil {
		return err
	}

	collection := database.Mg.DB.Collection("employees")
	_, err := collection.InsertOne(context.Background(), employee)
	if err != nil {
		return err
	}
	return c.JSON(employee)
}

func DeleteEmployee(c *fiber.Ctx) error {
	id := c.Params("id")
	collection := database.Mg.DB.Collection("employees")
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
	collection := database.Mg.DB.Collection("employees")
	filter := bson.M{"id": id}

	var employee model.Employee
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
