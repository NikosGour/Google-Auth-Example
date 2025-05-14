package api

import (
	log "github.com/NikosGour/logging/src"

	"github.com/NikosGour/google-oauth-example/app/auth"
	"github.com/NikosGour/google-oauth-example/app/handlers"
	"github.com/NikosGour/google-oauth-example/build"
	"github.com/NikosGour/google-oauth-example/storage"
	"github.com/NikosGour/google-oauth-example/types"

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

	app := SetupFiberApp()

	err := app.Listen(server.listening_addr)
	if err != nil {
		log.Fatal("%s", err)
	}

}

func SetupFiberApp() *fiber.App {

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

	return app
}
