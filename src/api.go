package api

import (
	log "github.com/NikosGour/logging/src"
	"golang.org/x/oauth2"

	"github.com/NikosGour/date_management_API/src/build"
	"github.com/NikosGour/date_management_API/src/storage"
	"github.com/NikosGour/date_management_API/src/types"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type Storage storage.Storage
type User types.User

type APIServer struct {
	storage        Storage
	listening_addr string
	OAuth_config   *oauth2.Config
	env_variables  map[string]string
}

func NewAPIServer(storage Storage, listening_addr string, oauth_config *oauth2.Config, dotenv map[string]string) *APIServer {
	this := &APIServer{storage: storage, listening_addr: listening_addr, OAuth_config: oauth_config, env_variables: dotenv}
	return this
}

func (this *APIServer) Start() {

	log.Debug("DEBUG_MODE = %t\n", build.DEBUG_MODE)

	app := fiber.New()
	app.Use(logger.New(logger.Config{
		Format: "${time} | ${status} | ${latency} | ${ip} | ${method} | ${path} | Params: ${queryParams} | ReqBody: ${body} | ResBody: ${resBody} | ${error}\n",
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"Nikos": 10})
	})

	app.Get("/oauth/google", this.authGoogleHandle)

	//redirect endpoint
	app.Get("/oauth/redirect", this.authRedirectHandle)

	err := app.Listen(this.listening_addr)
	if err != nil {
		log.Fatal(err)
	}

}
