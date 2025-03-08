package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kanishkmehta29/calorie-tracker/database"
	"github.com/kanishkmehta29/calorie-tracker/routes"
)

func main() {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Updated to your actual React port
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * 60 * 60,
	}))

	r.POST("/entry/create", routes.AddEntry)
	r.GET("/entry", routes.GetEntries)
	r.GET("/entry/:id", routes.GetEntryById)
	r.GET("/entry/ingredient/:ingredient", routes.GetEntriesByIngredient)
	r.PUT("/entry/update/:id", routes.UpdateEntry)
	r.PUT("/entry/update/ingredient/:id", routes.UpdateIngredient)
	r.DELETE("/entry/delete/:id", routes.DeleteEntry)

	database.Global_client = database.Connect()
	r.Run(":8080")

}
