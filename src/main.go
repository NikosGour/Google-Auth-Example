package main

import (
	"fmt"

	"github.com/NikosGour/date_management_API/src/build"
	log "github.com/NikosGour/logging/src"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	dotenv, err := godotenv.Read(".env")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("DEBUG_MODE = %t\n", build.DEBUG_MODE)
	fmt.Printf("env = %s\n", dotenv["MYSQL_ROOT_PASSWORD"])

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Nikos")
	})

	err = app.Listen(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
