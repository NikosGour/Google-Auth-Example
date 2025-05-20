package auth

import (
	"context"
	"encoding/json"
	"errors"
	"net/url"
	"time"

	"github.com/NikosGour/google-oauth-example/common"
	log "github.com/NikosGour/logging/src"

	"github.com/gofiber/fiber/v2"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	OAuth_config *oauth2.Config
)

func InitOAuthConfig(dotenv map[string]string) {
	oauthConfig := &oauth2.Config{
		ClientID:     dotenv["GOOGLE_CLIENT_ID"],
		ClientSecret: dotenv["GOOGLE_CLIENT_SECRET"],
		RedirectURL:  dotenv["GOOGLE_REDIRECT_URL"],
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	OAuth_config = oauthConfig
}

func AuthenticateUser(c *fiber.Ctx) error {
	log.Debug("Path: %s", c.Path())
	cookie := c.Cookies("token")
	log.Debug("cookie=%#v", cookie)

	if cookie == "" {
		log.Error("`token` cookie is empty")
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	// TODO: URL encode every cookie
	decoded, err := url.QueryUnescape(cookie)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	log.Debug("decoded=%#v", decoded)
	// TODO: Check if the cookie is not empty but doesn't match the token struct
	token := &oauth2.Token{}
	err = json.Unmarshal([]byte(decoded), token)
	if err != nil || *token == (oauth2.Token{}) {
		log.Error("%s. Got mangled token cookie: `%s`", err, cookie)

		return c.SendStatus(fiber.StatusBadRequest)
	}
	log.Debug("token=%#v", token)
	log.Debug("Valid=%v", token.Valid())

	if !token.Valid() {
		// TODO: Check if token is empty/valid
		err := RefreshToken(c, token)
		if err != nil {
			if errors.Is(err, common.ErrRedirected) {
				return nil
			}

			log.Error("Got error: `%s`", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		json_token, err := json.Marshal(token)
		if err != nil {
			return err
		}
		encoded := url.QueryEscape(string(json_token))
		cookie := &fiber.Cookie{Name: "token", Value: encoded, HTTPOnly: true, Expires: time.Now().Add(30 * 24 * time.Hour)}
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
	log.Info("Token was expired, trying to refresh. Token: %#v", token)
	token_source := OAuth_config.TokenSource(context.Background(), token)

	newToken, err := token_source.Token()

	// Invalid/Expired refresh token
	var erro *oauth2.RetrieveError
	if errors.As(err, &erro) {
		_ = c.Redirect("/oauth/google")
		return common.ErrRedirected
	}
	if err != nil {
		// No refresh token
		if err.Error() == "oauth2: token expired and refresh token is not set" {
			_ = c.Redirect("/oauth/google")
			return common.ErrRedirected
		}
		return err
	}

	*token = *newToken
	return nil
}
