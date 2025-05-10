package api

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	log "github.com/NikosGour/logging/src"
	"golang.org/x/oauth2"

	"github.com/NikosGour/date_management_API/build"
	"github.com/NikosGour/date_management_API/storage"
	"github.com/NikosGour/date_management_API/types"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type Storage storage.Storage
type User types.User

var (
	OAuth_config *oauth2.Config
)

type APIServer struct {
	storage        Storage
	listening_addr string
	env_variables  map[string]string
}

func NewAPIServer(storage Storage, listening_addr string, oauth_config *oauth2.Config, dotenv map[string]string) *APIServer {
	OAuth_config = oauth_config
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

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"Nikos": 10})
	})

	app.Get("/oauth/google", authGoogleHandle)

	//redirect endpoint
	app.Get("/oauth/redirect", authRedirectHandle)

	with_auth := app.Group("/api", AuthenticateUser)

	with_auth.Get("/logout", func(c *fiber.Ctx) error {
		c.ClearCookie("token")
		return c.Redirect("/")
	})

	with_auth.Get("/testing", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"Skase": 69})
	})

	err := app.Listen(server.listening_addr)
	if err != nil {
		log.Fatal("%s", err)
	}

}

func AuthenticateUser(c *fiber.Ctx) error {
	log.Debug("Path: %s", c.Path())
	cookie := c.Cookies("token")
	log.Debug("cookie=%#v", cookie)

	if cookie == "" {
		log.Error("`token` cookie is empty")
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	// TODO: Check if the cookie is not empty but doesn't match the token struct
	token := &oauth2.Token{}
	err := json.Unmarshal([]byte(cookie), token)
	if err != nil {
		log.Error("%s. Got mangled token cookie: `%s`", err, cookie)

		return c.SendStatus(fiber.StatusBadRequest)
	}

	if !token.Valid() {
		// TODO: Check if token is empty/valid
		err := RefreshToken(c, token)
		if err != nil {
			log.Error("Got error: `%s`", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		json_token, err := json.Marshal(token)
		if err != nil {
			return err
		}

		cookie := &fiber.Cookie{Name: "token", Value: string(json_token), HTTPOnly: true, Expires: time.Now().Add(24 * time.Hour)}
		c.Cookie(cookie)
	}

	// log.Debug("AccessToken=`%v`", token.AccessToken)
	// log.Debug("RefreshToken=`%v`", token.RefreshToken)
	// log.Debug("Expiry=`%v`", token.Expiry)
	// log.Debug("IsValid=`%v`", token.Valid())
	// log.Debug("TokenType=`%v`", token.TokenType)
	// log.Debug("Type=`%v`", token.Type())

	user, err := getUserProfile(c, token)
	if err != nil {
		log.Error("%s", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Google Auth don't work")
	}
	log.Debug("user=%#v", user)
	return c.Next()
}

func RefreshToken(c *fiber.Ctx, token *oauth2.Token) error {

	token_source := OAuth_config.TokenSource(context.Background(), token)

	newToken, err := token_source.Token()
	*token = *newToken

	var erro *oauth2.RetrieveError
	if errors.As(err, &erro) {
		return c.Redirect("/oauth/google")
	}

	return err
}
