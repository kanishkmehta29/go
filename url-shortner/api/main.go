package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"github.com/kanishkmehta29/url-shortner/routes"
)

func main(){
	err := godotenv.Load(".env")
	if err != nil{
		log.Printf("Error loading main .env %v\n",err)
	}
	app := fiber.New()

	app.Use(logger.New())

	app.Get("/:url",routes.ResolveUrl)
	app.Post("api/v1",routes.ShortenUrl)

	log.Fatal(app.Listen(os.Getenv("APP_PORT")))


}