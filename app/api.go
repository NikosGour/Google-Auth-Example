package api

import (
	log "github.com/NikosGour/logging/src"

	"github.com/NikosGour/date_management_API/app/auth"
	"github.com/NikosGour/date_management_API/app/handlers"
	"github.com/NikosGour/date_management_API/build"
	"github.com/NikosGour/date_management_API/storage"
	"github.com/NikosGour/date_management_API/types"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type Storage storage.Storage
type User types.User

type APIServer struct {
	storage        Storage
	listening_addr string
	env_variables  map[string]string
}

func NewAPIServer(storage Storage, listening_addr string, dotenv map[string]string) *APIServer {
	this := &APIServer{storage: storage, listening_addr: listening_addr, env_variables: dotenv}
	return this
}

func (server *APIServer) Start() {

	log.Debug("DEBUG_MODE = %t\n", build.DEBUG_MODE)

	app := fiber.New()

	app.Use(logger.New(logger.Config{
		Format: "${time} | ${status} | ${latency} | ${ip} | ${method} | ${path} | Params: ${queryParams} | ReqBody: ${body} | ResBody: ${resBody} | ${error}\n",
	}))
	app.Use(favicon.New())

	app.Get("/", handlers.RootHandle)

	app.Get("/oauth/google", auth.GoogleHandle)
	app.Get("/oauth/redirect", auth.RedirectHandle)

	with_auth := app.Group("/api", auth.AuthenticateUser)
	with_auth.Get("/logout", auth.LogoutHandle)
	with_auth.Get("/testing", handlers.TestingHandle)

	err := app.Listen(server.listening_addr)
	if err != nil {
		log.Fatal("%s", err)
	}

}
