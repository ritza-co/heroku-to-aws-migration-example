package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var db = make(map[string]string)
var employeesCollection *mongo.Collection
var client *mongo.Client
var ctx = context.TODO()

// Replace these values with the URLs of your Kafka, and MongoDB instances
const MONGODBURL = "YOUR_URL_HERE"
const KAFKAURL = "YOUR_URL_HERE"

const topic = "employee_updates"

type Employee struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	FirstName string             `bson:"firstname,omitempty"`
	LastName  string             `bson:"lastname,omitempty"`
	Email     string             `bson:"email,omitempty"`
	HourlyPay string             `bson:"hourlypay,omitempty"`
}

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(CORSMiddleware())
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// As a requirement for Elastic Beanstalk, the root path must return a 200
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello World!")
	})

	employees := r.Group("/employees")
	{
		employees.GET("/", GetEmployees)
		employees.POST("/", CreateEmployee)
		employees.DELETE("/:id", DeleteEmployee)
	}

	return r
}

func GetEmployees(c *gin.Context) {

	// Get all employees from MongoDB
	cursor, err := employeesCollection.Find(ctx, bson.D{})
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to execute MongoDB query"})
		return
	}
	defer cursor.Close(ctx)

	// Collect documents into a slice
	var results []bson.M
	for cursor.Next(ctx) {
		var result bson.M
		err := cursor.Decode(&result)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to decode MongoDB document"})
			return
		}
		results = append(results, result)
	}

	// Check if there were any errors during iteration
	if err := cursor.Err(); err != nil {
		c.JSON(500, gin.H{"error": "Error during MongoDB cursor iteration"})
		return
	}

	// if no results are found, return an empty slice
	if results == nil {
		results = []bson.M{}
	}

	// Return the documents as a JSON array
	c.JSON(200, results)
}

func CreateEmployee(c *gin.Context) {
	// Parse the JSON request body into an Employee struct
	var employee Employee
	if err := c.ShouldBindJSON(&employee); err != nil {
		fmt.Println("error", err)
		c.JSON(400, gin.H{"error": "Invalid request payload"})
		return
	}
	// Insert the new employee document
	result, err := employeesCollection.InsertOne(ctx, employee)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to insert employee document into MongoDB"})
		return
	}

	conn, err := kafka.DialLeader(context.Background(), "tcp", KAFKAURL, topic, 0)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

	// Publish a message to the Kafka topic
	msg := kafka.Message{
		Value: []byte(fmt.Sprintf(`{"id": "%s"}`, result.InsertedID.(primitive.ObjectID).Hex())),
	}
	_, err = conn.WriteMessages(
		msg,
	)

	if err != nil {
		log.Fatal("failed to write messages:", err)
	}

	if err := conn.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}
	// Return the ID of the newly created employee
	c.JSON(201, gin.H{"id": result.InsertedID})
}

func DeleteEmployee(c *gin.Context) {
	// Get the employee ID from the URL path
	idParam := c.Param("id")

	employeeId, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid employee ID"})
		return
	}
	// Delete the employee document with the specified ID
	result, err := employeesCollection.DeleteOne(ctx, bson.M{"_id": employeeId})
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete employee document from MongoDB"})
		return
	}

	// Check if any document was deleted
	if result.DeletedCount == 0 {
		c.JSON(404, gin.H{"error": "Employee not found"})
		return
	}

	// Return success message
	c.JSON(200, gin.H{"message": "Employee deleted successfully"})
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "*")
		c.Header("Access-Control-Allow-Methods", "*")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// Establish connection with MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(MONGODBURL))
	if err != nil {
		panic(err)
	}

	// Check the connection
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}

	// Set the collection
	employeesCollection = client.Database("test").Collection("employees")

	r := setupRouter()

	r.Run(":8080")
}
