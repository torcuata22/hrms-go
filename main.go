package main

import (
	"go.mongodb.org/mongo-driver/mongo"
)

//Mongo DB struct

type MongoInstance struct {
	client *mongo.Client
	DB     *mongo.Database
}

var mg MongoInstance

const dbName = "hrms"
const mongoUri = ""
