package routes

import (
	ctr "github.com/BassemArfaoui/Weazy-Server/controllers"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {

	//home route
	app.Get("/", ctr.Home)

	//chat routes
	app.Get("/chats/:userId", ctr.GetChatsByUserId)

	app.Put("/edit-chat*/:chatId", ctr.EditChat)

}
