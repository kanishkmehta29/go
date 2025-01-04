package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Employee struct{
	Id string `bson:"_id,omitempty" json:"id,omitempty"`
	Name string `json:"name"`
	Salary int64 `json:"salary"`
	Age int `json:"age"`
}

var Global_client *mongo.Client

func Connect() *mongo.Client{
   err := godotenv.Load(".env")
   if err != nil{
	log.Fatalf("Error loading the env file %v\n",err)
   }
   url := os.Getenv("DATABASE_URL")

   client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(url))
   if err != nil{
	log.Fatalf("Error connecting to database %v\n",err)
   }

   err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatalf("Error connecting to database %v\n",err)
	}

   fmt.Println("Sucessfully connected to the database")   
   return client

}