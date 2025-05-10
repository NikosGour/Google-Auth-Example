package auth

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	log "github.com/NikosGour/logging/src"

	"github.com/gofiber/fiber/v2"

	"golang.org/x/oauth2"
)

var (
	OAuth_config *oauth2.Config
)

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
