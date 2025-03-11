package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id            primitive.ObjectID `bson:"_id" json:"id"`
	First_name    string             `json:"first_name" validate:"required,min=2,max=100"`
	Last_name     string             `json:"last_name" validate:"min=1,max=100"`
	Email         string             `json:"email" validate:"email,required"`
	Password      string             `json:"password" validate:"required,min=6"`
	Phone         string             `json:"phone" validate:"numeric,required,len=10"`
	Token         string             `json:"token"`
	User_type     string             `json:"user_type" validate:"required,oneof=admin user"`
	Refresh_token string             `json:"refresh_token"`
	Created_at    time.Time          `json:"created_at"`
	Updated_at    time.Time          `json:"updated_at"`
	User_id       string             `json:"user_id"`
}
