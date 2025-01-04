package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/kanishkmehta29/hrms/database"
	"github.com/kanishkmehta29/hrms/router"

)

func main(){

	app := fiber.New()
    database.Global_client = database.Connect()
	router.ManageRoutes(app)
	fmt.Println("Starting server on port 8080")
	app.Listen(":8080")
	
}