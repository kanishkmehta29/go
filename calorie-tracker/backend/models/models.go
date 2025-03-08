package models

import( "go.mongodb.org/mongo-driver/bson/primitive")

type Entry struct {
	Id          primitive.ObjectID `json:"id" bson:"_id"`
	Dish        *string           `json:"dish" validate:"required"`
	Ingredients *string           `json:"ingredients" validate:"required"`
	Calories    *float64          `json:"calories" validate:"required,gte=0"`
	Fat         *float64          `json:"fat" validate:"required,gte=0"`
}

type Ingredient struct{
	Ingredients *string           `json:"ingredients" validate:"required"`
}