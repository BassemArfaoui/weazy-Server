package main

import (
	"fmt"
	"os"

	"github.com/BassemArfaoui/Weazy-Server/config"
	"github.com/BassemArfaoui/Weazy-Server/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {

	//database and cloud 
	config.Connect()
	config.InitCloud()


	//port
	if err := godotenv.Load(); err != nil {
		fmt.Println("error loading env vars")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3333"
		fmt.Println("PORT not set. Using default port:", port)
	}
	
	//app
	app := fiber.New()


	//middlewares
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", 
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	app.Use(logger.New())


	//routes
	routes.Setup(app)


	//run
	app.Listen(":" + port)
}
