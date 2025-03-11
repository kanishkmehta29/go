package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kanishkmehta29/jwt-auth/controllers"
)

func AuthRoutes(incomingRoutes *gin.Engine){
	incomingRoutes.POST("users/signup",controllers.Signup)
	incomingRoutes.POST("users/signin",controllers.Signin)
	
}