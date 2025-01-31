package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/torcuata22/hrms-mongo/database"
	model "github.com/torcuata22/hrms-mongo/models"
	"github.com/torcuata22/hrms-mongo/routes"
	"go.mongodb.org/mongo-driver/bson"
)

func initdatabase(t *testing.T) {
	err := database.Connect()
	if err != nil {
		t.Fatalf("Error connecting to database: %v", err)
	}
}

func setupapp() *fiber.App {
	app := fiber.New()
	routes.SetupRoutes(app)
	return app
}

// func insertTestEmployee(t *testing.T, employee model.Employee) string {
// 	// Insert the employee into the database
// 	inserted, err := database.Mg.DB.Collection("employees").InsertOne(context.TODO(), employee)
// 	if err != nil {
// 		t.Fatalf("Error inserting test employee: %v", err)
// 	}

// Convert the inserted ID from ObjectID to string
// insertedID := inserted.InsertedID.(primitive.ObjectID).Hex() // .Hex() converts ObjectID to a string

// 	return insertedID
// }

func makeRequest(t *testing.T, app *fiber.App, method, url string) *http.Response {
	//request
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		t.Fatalf("Error making test request: %v", err)
	}
	//response
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Error making test request: %v", err)
	}
	return resp
}

func TestGetEmployees(t *testing.T) {
	//initialize db
	initdatabase(t)
	// Create Fiber app and setup routes
	app := setupapp()

	// Create test request
	resp := makeRequest(t, app, "GET", "/employees")

	// Check response status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// Decode response body
	var employees []model.Employee
	err := json.NewDecoder(resp.Body).Decode(&employees)
	if err != nil {
		t.Fatalf("Error decoding response body: %v", err)
	}

	// Check that at least one employee is returned
	if len(employees) == 0 {
		t.Errorf("Expected at least one employee, got none")
	}
}

func TestGetEmployee(t *testing.T) {
	//initialize db
	initdatabase(t)
	// Create Fiber app and setup routes
	app := setupapp()
	resp := makeRequest(t, app, "GET", "/employee/1") //use existing db data because I don't have test db set up

	//check status code
	if resp.StatusCode != http.StatusOK {
		t.Error("Expected status code 200, got", resp.StatusCode)
	}

	//Decode resp body
	var employee model.Employee
	err := json.NewDecoder(resp.Body).Decode(&employee)
	if err != nil {
		t.Fatalf("Error decoding response body: %v", err)
	}

	//check employee
	if employee.ID != "1" {
		t.Error("Expected employee ID to be 1, got", employee.ID)
	}
	if employee.Name != "Loki Rabito" {
		t.Errorf("Expected employee name %s, got %s", "Loki Rabito", employee.Name)
	}
	if employee.Salary != 499.99 {
		t.Errorf("Expected employee salary %f, got %f", 499.99, employee.Salary)
	}
	if employee.Age != 4 {
		t.Errorf("Expected employee age %d, got %d", 4, employee.Age)
	}
}

func TestCreateEmployee(t *testing.T) {
	//initialize db
	initdatabase(t)
	// Create Fiber app and setup routes
	app := setupapp()

	//create test employee
	testEmployee := model.Employee{
		ID:     uuid.New().String(),
		Name:   "John Doe",
		Salary: 50000,
		Age:    30,
	}
	//convert to JSON
	jsonEmployee, err := json.Marshal(testEmployee)
	if err != nil {
		t.Fatalf("Error marshalling test employee: %v", err)
	}

	//build the request
	req, err := http.NewRequest("POST", "/employees", bytes.NewBuffer(jsonEmployee))
	if err != nil {
		t.Fatalf("Error making test request: %v", err)
	}

	//set the content type//set headers:
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Error making test request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, resp.StatusCode)
	}

	//check response body
	var createdEmployee model.Employee
	err = json.NewDecoder(resp.Body).Decode(&createdEmployee)
	if err != nil {
		t.Fatalf("Error decoding response body: %v", err)
	}
	if createdEmployee.Name != testEmployee.Name {
		t.Errorf("Expected employee name %s, got %s", testEmployee.Name, createdEmployee.Name)
	}
	if createdEmployee.Salary != testEmployee.Salary {
		t.Errorf("Expected employee salary %f, got %f", testEmployee.Salary, createdEmployee.Salary)
	}
	if createdEmployee.Age != testEmployee.Age {
		t.Errorf("Expected employee age %d, got %d", testEmployee.Age, createdEmployee.Age)
	}
}

func TestUpdateEmployee(t *testing.T) {
	//initialize db
	initdatabase(t)
	// Create Fiber app and setup routes
	app := setupapp()

	//create test employee
	testEmployee := model.Employee{
		ID:     "ee5c432f-ba8a-4b69-b512-09d64e18b788",
		Name:   "John Doe",
		Salary: 50000.0,
		Age:    30,
	}
	//convert to JSON
	jsonEmployee, err := json.Marshal(testEmployee)
	if err != nil {
		t.Fatalf("Error marshalling test employee: %v", err)
	}

	//build the request
	req, err := http.NewRequest("GET", "/employee/"+testEmployee.ID, bytes.NewBuffer(jsonEmployee))
	if err != nil {
		t.Fatalf("Error making test request: %v", err)
	}

	//set the content type//set headers:
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Error making test request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	//update the employee information
	testEmployee.Name = "Jane Doe"
	testEmployee.Salary = 60000.0
	testEmployee.Age = 30

	//convert to JSON
	jsonEmployee, err = json.Marshal(testEmployee)
	if err != nil {
		t.Fatalf("Error marshalling test employee: %v", err)
	}

	//build the request
	req, err = http.NewRequest("PUT", "/employee/"+testEmployee.ID, bytes.NewBuffer(jsonEmployee))
	if err != nil {
		t.Fatalf("Error making test request: %v", err)
	}

	//set the content type//set headers:
	req.Header.Set("Content-Type", "application/json")

	resp, err = app.Test(req)
	if err != nil {
		t.Fatalf("Error making test request: %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("Expected status code %d, got %d", http.StatusNoContent, resp.StatusCode)
	}

	//retrieve updated employee from the database
	var updatedEmployee model.Employee
	collection := database.Mg.DB.Collection("employees")
	filter := bson.M{"id": testEmployee.ID}
	err = collection.FindOne(context.Background(), filter).Decode(&updatedEmployee)
	if err != nil {
		t.Fatalf("Error retrieving updated employee: %v", err)
	}

	// Verify that the employee's data was updated correctly
	if updatedEmployee.Name != "Jane Doe" {
		t.Errorf("Expected employee name %s, got %s", "Jane Doe", updatedEmployee.Name)
	}
	if updatedEmployee.Salary != 60000.0 {
		t.Errorf("Expected employee salary %f, got %f", 60000.0, updatedEmployee.Salary)
	}
	if updatedEmployee.Age != 30 {
		t.Errorf("Expected employee age %d, got %d", 30, updatedEmployee.Age)
	}
}
