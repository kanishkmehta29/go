package routes

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/kanishkmehta29/calorie-tracker/database"
	"github.com/kanishkmehta29/calorie-tracker/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Initialize validator
var validate = validator.New()

func AddEntry(ctx *gin.Context) {
	var newEntry models.Entry
	ctx2, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := database.Global_client.Database("calorie_database").Collection("calorie_collection")

	err := ctx.BindJSON(&newEntry)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "can't parse input properly", "error": err})
		return
	}

	// Perform explicit validation
	err = validate.Struct(newEntry)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "validation failed",
			"errors":  validationErrors.Error(),
		})
		return
	}

	newEntry.Id = primitive.NewObjectID()
	res, err := collection.InsertOne(ctx2, newEntry)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "error in inserting in database", "error": err})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Entry created sucessfully", "new Entry": res})

}

func GetEntries(ctx *gin.Context) {
	collection := database.Global_client.Database("calorie_database").Collection("calorie_collection")

	ctx2, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx2, bson.M{})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(ctx2)

	var entries []bson.M
	err = cursor.All(ctx2, &entries)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"entries": entries})
}

func GetEntriesByIngredient(ctx *gin.Context) {
	ingredient := ctx.Params.ByName("ingredient")
	ctx2, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var entries []bson.M
	collection := database.Global_client.Database("calorie_database").Collection("calorie_collection")

	cursor, err := collection.Find(ctx2, bson.M{"ingredients": ingredient})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	err = cursor.All(ctx2, &entries)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	ctx.JSON(http.StatusOK, entries)

}

func GetEntryById(ctx *gin.Context) {
	id := ctx.Params.ByName("id")
	docId, _ := primitive.ObjectIDFromHex(id)

	collection := database.Global_client.Database("calorie_database").Collection("calorie_collection")

	ctx2, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	var entry bson.M
	err := collection.FindOne(ctx2, bson.M{"_id": docId}).Decode(&entry)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		cancel()
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"Entry": entry})
	cancel()

}

func UpdateEntry(ctx *gin.Context) {
	id := ctx.Params.ByName("id")
	docId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid ID format", "error": err.Error()})
		return
	}

	collection := database.Global_client.Database("calorie_database").Collection("calorie_collection")
	ctx2, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var newEntry models.Entry
	err = ctx.BindJSON(&newEntry)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "can't parse input properly", "error": err.Error()})
		return
	}

	// Perform explicit validation
	err = validate.Struct(newEntry)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "validation failed",
			"errors":  validationErrors.Error(),
		})
		return
	}

	newEntry.Id = docId

	result, err := collection.UpdateOne(
		ctx2,
		bson.M{"_id": docId},
		bson.M{"$set": bson.M{
			"dish":        newEntry.Dish,
			"ingredients": newEntry.Ingredients,
			"calories":    newEntry.Calories,
			"fat":         newEntry.Fat,
		}},
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "error updating entry", "error": err.Error()})
		return
	}

	if result.MatchedCount == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "entry not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "entry updated successfully"})
}

func UpdateIngredient(ctx *gin.Context) {
	id := ctx.Params.ByName("id")
	docId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid ID format", "error": err.Error()})
		return
	}

	collection := database.Global_client.Database("calorie_database").Collection("calorie_collection")
	ctx2, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var newIngredient models.Ingredient
	err = ctx.BindJSON(&newIngredient)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "can't parse input properly", "error": err.Error()})
		return
	}

	// Perform explicit validation
	err = validate.Struct(newIngredient)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "validation failed",
			"errors":  validationErrors.Error(),
		})
		return
	}

	result, err := collection.UpdateOne(
		ctx2,
		bson.M{"_id": docId},
		bson.M{"$set": bson.M{
			"ingredients": newIngredient.Ingredients,
		}},
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "error updating entry", "error": err.Error()})
		return
	}

	if result.MatchedCount == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "entry not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "entry updated successfully"})
}

func DeleteEntry(ctx *gin.Context) {
	id := ctx.Params.ByName("id")
	docId, _ := primitive.ObjectIDFromHex(id)

	var ctx2, cancel = context.WithTimeout(context.Background(), 10*time.Second)

	collection := database.Global_client.Database("calorie_database").Collection("calorie_collection")

	result, err := collection.DeleteOne(ctx2, bson.M{"_id": docId})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		cancel()
		return
	}

	if result.DeletedCount == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Entry not found"})
		cancel()
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Entry deleted successfully"})
	cancel()
}
