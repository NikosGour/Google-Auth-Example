package handlers

import "github.com/gofiber/fiber/v2"

func RootHandle(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"Nikos": 10})
}

func TestingHandle(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"Skase": 69})
}
