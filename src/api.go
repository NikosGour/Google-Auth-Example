package api

import (
	"fmt"

	log "github.com/NikosGour/logging/src"

	"github.com/NikosGour/date_management_API/src/build"
	"github.com/NikosGour/date_management_API/src/types"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type Storage types.Storage

type APIServer struct {
	storage        Storage
	listening_addr string
	env_variables  map[string]string
}

func NewAPIServer(storage Storage, listening_addr string, dotenv map[string]string) *APIServer {
	this := &APIServer{storage: storage, listening_addr: listening_addr, env_variables: dotenv}
	return this
}

func (this *APIServer) Start() {

	fmt.Printf("DEBUG_MODE = %t\n", build.DEBUG_MODE)

	app := fiber.New()
	app.Use(logger.New(logger.Config{
		Format: "${time} | ${status} | ${latency} | ${ip} | ${method} | ${path} | Params: ${queryParams} | ReqBody: ${body} | ResBody: ${resBody} | ${error}\n",
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"Nikos": 10})
	})

	err := app.Listen(this.listening_addr)
	if err != nil {
		log.Fatal(err)
	}
}
