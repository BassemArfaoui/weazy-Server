package main

import (
	"os"
	db "github.com/BassemArfaoui/Weazy-Server/config"
	routes "github.com/BassemArfaoui/Weazy-Server/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {

	//database
	db.Connect()


	//port
	godotenv.Load()
	port := os.Getenv("APP_PORT")


	
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



	app.Listen(":" + port)
}
