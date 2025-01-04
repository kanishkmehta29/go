package main

import (
	"fmt"
	"log"
	"os"
    "gorm.io/gorm"
	"github.com/gofiber/fiber"
	"github.com/joho/godotenv"
	"github.com/kanishkmehta29/crm-golang/database"
	"github.com/kanishkmehta29/crm-golang/lead"
	"gorm.io/driver/postgres"
)

func setupRoutes(app *fiber.App){
	app.Get("api/v1/lead/:id",lead.GetLead)
	app.Get("api/v1/lead",lead.GetLeads)
	app.Post("api/v1/lead",lead.NewLead)
	app.Delete("api/v1/lead/:id",lead.DeleteLead)
}

func initDatabase(){
	err := godotenv.Load(".env")
	if err != nil{
		log.Fatalln("Error loading .env")
	}
	url := os.Getenv("DATABASE_URL")
	database.DB,err = gorm.Open(postgres.Open(url), &gorm.Config{})
    if err != nil{
		log.Fatalf("Error connecting to database %v\n",err)
	}
	database.DB.AutoMigrate(&lead.Lead{})
	fmt.Println("Sucessfully connected to the database")
}

func main(){

	app := fiber.New()
	initDatabase()
	setupRoutes(app)
	fmt.Println("Starting server on port 8080")
    app.Listen(":8080")
}