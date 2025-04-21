package routes


import (
	"github.com/gofiber/fiber/v2"
	ctr "github.com/BassemArfaoui/Quinsat-Server-Side/controllers"
)


func Setup(app *fiber.App) {

	//home route
	app.Get("/" , ctr.Home)



	//auth




}
