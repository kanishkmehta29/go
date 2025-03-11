package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kanishkmehta29/jwt-auth/controllers"
	"github.com/kanishkmehta29/jwt-auth/middleware"
)

func UserRoutes(incomingRoutes *gin.Engine){
	incomingRoutes.Use(middleware.Authenticate)
	incomingRoutes.GET("/users",controllers.GetUsers)
	incomingRoutes.GET("/users/:user_id",controllers.GetUserById)
}