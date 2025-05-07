package routes

import (
	ctr "github.com/BassemArfaoui/Weazy-Server/controllers"
	"github.com/gofiber/fiber/v2"
)


func Setup(app *fiber.App) {
	// Home route
	app.Get("/", ctr.Home)

	// Upload route
	app.Post("/upload", ctr.Upload)

	// Chat routes
	app.Get("/chats/:userId", ctr.GetChatsByUserId)
	app.Put("/edit-chat/:chatId", ctr.EditChat)
	app.Delete("/delete-chat/:chatId", ctr.DeleteChat)
	app.Post("/create-chat", ctr.CreateChat)
	app.Get("/chat/:chatId/:userId", ctr.GetChatById)

	//message routes
	app.Post("/save-payload", ctr.SavePayload)
	app.Post("/save-response" , ctr.SaveResponse)

	//wishlist routes
	app.Get("/wishlist/:userId", ctr.GetWishlistByUserId)
	app.Post("/add-to-wishlist", ctr.SaveToWishlist)
	app.Delete("/delete-wishlist-item/:userId/:productId" , ctr.DeleteWishlistItem)
}
