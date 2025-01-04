package main

import (
    "context"
    "fmt"
    "log"
    "net/http"

    "github.com/julienschmidt/httprouter"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "github.com/kanishkmehta29/mongo-golang/controllers" // Adjust the import path as necessary
)

const uri = "mongodb://localhost:27017"

func main() {
    r := httprouter.New()
    uc := controllers.NewUserController(getSession())
    r.GET("/user/:id", uc.GetUser)
    r.POST("/user", uc.CreateUser)
    r.DELETE("/user/:id", uc.DeleteUser)
    fmt.Println("Starting server on port 8080")
    http.ListenAndServe(":8080", r)
}

func getSession() *mongo.Client {
    serverAPI := options.ServerAPI(options.ServerAPIVersion1)
    opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

    client, err := mongo.Connect(context.TODO(), opts)
    if err != nil {
        log.Fatalf("Failed to connect to MongoDB: %v", err)
    }

    // Send a ping to confirm a successful connection
    var result bson.M
    if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Decode(&result); err != nil {
        log.Fatalf("Failed to ping MongoDB: %v", err)
    }
    fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

    return client
}