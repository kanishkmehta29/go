package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	helper "github.com/kanishkmehta29/jwt-auth/helper"
)

func Authenticate(ctx *gin.Context) {
	clientToken := ctx.Request.Header.Get("Authorization")
	if clientToken == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "No authorization token provided",
		})
		return
	}
	claims, err := helper.ValidateToken(clientToken)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
	}

	ctx.Set("email", claims.Email)
	ctx.Set("first_name", claims.Name)
	ctx.Set("user_id", claims.UserId)
	ctx.Set("user_type", claims.UserType)
	ctx.Next()
}
