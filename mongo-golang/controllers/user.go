package controllers

import (
    "context"
    "encoding/json"
    "net/http"

    "github.com/julienschmidt/httprouter"
    "github.com/kanishkmehta29/mongo-golang/models"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
)

type UserController struct {
    client *mongo.Client
}

func NewUserController(client *mongo.Client) *UserController {
    return &UserController{client}
}

func (uc UserController) GetUser(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
    id := p.ByName("id")

    oid, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    var user models.User
    collection := uc.client.Database("your_database").Collection("users")
    err = collection.FindOne(context.TODO(), bson.M{"_id": oid}).Decode(&user)
    if err != nil {
        w.WriteHeader(http.StatusNotFound)
        return
    }

    userJSON, err := json.Marshal(user)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(userJSON)
}

func (uc UserController) CreateUser(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
    var user models.User
    json.NewDecoder(req.Body).Decode(&user)

    user.Id = primitive.NewObjectID()
    collection := uc.client.Database("your_database").Collection("users")
    _, err := collection.InsertOne(context.TODO(), user)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    userJSON, err := json.Marshal(user)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    w.Write(userJSON)
}

func (uc UserController) DeleteUser(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
    id := p.ByName("id")

    oid, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    collection := uc.client.Database("your_database").Collection("users")
    _, err = collection.DeleteOne(context.TODO(), bson.M{"_id": oid})
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
}