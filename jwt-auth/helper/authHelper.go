package helper

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func MatchUserTypeToUid(ctx *gin.Context,user_id string)(err error){
   userType := ctx.GetString("user_type")
   uId := ctx.GetString("user_id")
   err = nil

   if userType == "user" && uId != user_id{
	err = errors.New("unauthorized to access this information")
   }
   
   return err
}