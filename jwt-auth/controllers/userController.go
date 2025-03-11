package controllers

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/kanishkmehta29/jwt-auth/database"
	"github.com/kanishkmehta29/jwt-auth/helper"
	"github.com/kanishkmehta29/jwt-auth/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var validate = validator.New()

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func VerifyPassword(userPassword string, providedPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(providedPassword))
	if err != nil {
		return false, errors.New("incorrect password")
	}
	return true, nil
}

func Signup(ctx *gin.Context) {
	ctx2, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User

	err := ctx.BindJSON(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Problem in parsing the request",
			"error":   err.Error(),
		})
		return
	}

	err = validate.Struct(user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Wrong format of input fields",
			"error":   err.Error(),
		})
		return
	}

	count, err := database.UserCollection.CountDocuments(ctx2, bson.M{"email": user.Email})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "some error occured while validating email",
			"error":   err.Error(),
		})
		return
	}
	if count > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "this email already exists, please choose a different email",
		})
		return
	}

	count, err = database.UserCollection.CountDocuments(ctx2, bson.M{"phone": user.Phone})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "some error occured while validating phone",
			"error":   err.Error(),
		})
		return
	}
	if count > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "this phone already exists, please choose a different phone",
		})
		return
	}

	ist, _ := time.LoadLocation("Asia/Kolkata")
	now := time.Now().In(ist)
	user.Created_at = now
	user.Updated_at = now
	user.Id = primitive.NewObjectID()
	user.User_id = user.Id.Hex()
	user.Password, _ = HashPassword(user.Password)
	token, refreshToken, err := helper.GenerateTokens(user.User_id, user.Email, user.First_name, user.User_type)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "some error occured while generating tokens",
			"error":   err.Error(),
		})
		return
	}
	user.Token = token
	user.Refresh_token = refreshToken

	res, err := database.UserCollection.InsertOne(ctx2, user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "some error occured while inserting entry to database",
			"error":   err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, res)

}

func Signin(ctx *gin.Context) {
	ctx2, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	var foundUser models.User

	err := ctx.BindJSON(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Problem in parsing the request",
			"error":   err.Error(),
		})
		return
	}

	err = database.UserCollection.FindOne(ctx2, bson.M{"email": user.Email}).Decode(&foundUser)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "User not found",
			"error":   err.Error(),
		})
		return
	}

	passwordIsValid, err := VerifyPassword(foundUser.Password, user.Password)
	if !passwordIsValid {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	token, refreshToken, err := helper.GenerateTokens(foundUser.User_id, foundUser.Email, foundUser.First_name, foundUser.User_type)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "some error occured while generating tokens",
			"error":   err.Error(),
		})
		return
	}

	err = helper.UpdateTokens(token, refreshToken, foundUser.User_id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "some error occured while updating tokens",
			"error":   err.Error(),
		})
		return
	}
	err = database.UserCollection.FindOne(ctx2, bson.M{"user_id": foundUser.User_id}).Decode(&foundUser)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "some error occured while searching for updated user",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, foundUser)

}

func GetUsers(ctx *gin.Context) {
	ctx2, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	uType := ctx.GetString("user_type")
	if uType != "admin" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Unauthorized access",
		})
		return
	}

	var users []models.User

	cursor, err := database.UserCollection.Find(ctx2, bson.M{})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error in fetching data from database",
			"error":   err.Error()})
	}

	err = cursor.All(ctx2, &users)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error in decoding data from cursor",
			"error":   err.Error()})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"users": users,
	})

}

func GetUserById(ctx *gin.Context) {
	user_id := ctx.Params.ByName("user_id")
	ctx2, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := helper.MatchUserTypeToUid(ctx, user_id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Unauthorized access",
			"error":   err.Error(),
		})
		return
	}

	var user models.User

	err = database.UserCollection.FindOne(ctx2, bson.M{"user_id": user_id}).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusOK, user)

}
