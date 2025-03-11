package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Global_client *mongo.Client
var UserCollection *mongo.Collection

// Connect establishes a connection to MongoDB and returns the client
func Connect() *mongo.Client {
	// Load environment variables from .env file
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	url := os.Getenv("DATABASE_URL")

	// Create a context with timeout for the connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(url))
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	// Verify the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	fmt.Println("Successfully connected to the database")
	return client
}

func OpenConnection(client *mongo.Client, collection_name string) *mongo.Collection {
	return client.Database("auth-db").Collection(collection_name)
}
