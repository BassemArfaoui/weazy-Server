package middlewares

import "github.com/gofiber/fiber/v2"


func TokenCheck(c *fiber.Ctx) error{
	return c.Next()
}


func JwtCheck(c *fiber.Ctx) error{
	return c.Next()
}



