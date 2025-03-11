package helper

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/kanishkmehta29/jwt-auth/database"
	"go.mongodb.org/mongo-driver/bson"
)

// TokenClaims represents the claims in our JWT
type TokenClaims struct {
	UserId   string
	Email    string
	Name     string
	UserType string
	jwt.RegisteredClaims
}

func GenerateTokens(userId string, email string, fname string, user_type string) (token string, refreshToken string, err error) {

	// Create claims with our custom fields and standard claims for expiration time
	claims := &TokenClaims{
		UserId:   userId,
		Email:    email,
		Name:     fname,
		UserType: user_type,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}

	refreshClaims := &TokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		},
	}

	godotenv.Load(".env")
	key := os.Getenv("JWT_SECRET")
	token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(key))
	if err != nil {
		return "", "", err
	}
	refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(key))
	if err != nil {
		return "", "", err
	}
	return token, refreshToken, nil

}

func UpdateTokens(token string, refreshToken string, user_id string) (err error) {
	ctx2, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ist, _ := time.LoadLocation("Asia/Kolkata")

	_, err = database.UserCollection.UpdateOne(
		ctx2,
		bson.M{"user_id": user_id},
		bson.M{
			"$set": bson.M{
				"token":        token,
				"refreshToken": refreshToken,
				"updated_at":   time.Now().In(ist),
			},
		},
	)
	return err
}

func ValidateToken(signedToken string) (claims *TokenClaims, msg error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&TokenClaims{},
		func(token *jwt.Token) (interface{}, error) {
			godotenv.Load(".env")
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		},
	)

	if err != nil {
		return nil, err
	}

	// Validate token claims
	claims, ok := token.Claims.(*TokenClaims)
	if ok && token.Valid {
		// Check expiration time
		ist, _ := time.LoadLocation("Asia/Kolkata")
		if time.Now().In(ist).After(claims.ExpiresAt.Time) {
			return nil, fmt.Errorf("token has expired")
		}
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
