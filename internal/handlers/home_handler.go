package handlers

import "github.com/gofiber/fiber/v2"

func HomeHandler(c *fiber.Ctx) error {
	return c.SendString("Hola, bienvenido a go auth")
}
