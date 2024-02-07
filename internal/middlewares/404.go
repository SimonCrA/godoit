package middlewares

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// NotFoundHandler is a middleware to catch all 404 errors
func NotFoundHandler(c *fiber.Ctx) error {
	fmt.Println(c.Response())

	if c.Response().StatusCode() == fiber.StatusNotFound {
		return c.Status(fiber.StatusNotFound).SendString("Ups, estas perdido?  parece que esta ruta no existe!")
	}

	return c.Next()
}
