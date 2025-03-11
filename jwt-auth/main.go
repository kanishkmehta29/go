package main

import (
	"github.com/gin-gonic/gin"
	"github.com/kanishkmehta29/jwt-auth/database"
	"github.com/kanishkmehta29/jwt-auth/routes"
)

func main() {
	database.Global_client = database.Connect()
	database.UserCollection = database.OpenConnection(database.Global_client, "user")

	router := gin.Default()

	routes.AuthRoutes(router)
	routes.UserRoutes(router)

	router.Run(":8080")
}
