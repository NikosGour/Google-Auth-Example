package auth

import (
	"context"
	"encoding/json"
	"io"
	"time"

	"github.com/NikosGour/google-oauth-example/types"
	log "github.com/NikosGour/logging/src"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
)

type User = types.User

func GoogleHandle(c *fiber.Ctx) error {
	url := OAuth_config.AuthCodeURL("nikos", oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	return c.Redirect(url)
}
func RedirectHandle(c *fiber.Ctx) error {
	code := c.Query("code")
	if code == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	token, err := OAuth_config.Exchange(context.Background(), code)
	// log.Debug("AccessToken=`%v`", token.AccessToken)
	// log.Debug("RefreshToken=`%v`", token.RefreshToken)
	// log.Debug("Expiry=`%v`", token.Expiry)
	// log.Debug("IsValid=`%v`", token.Valid())
	// log.Debug("TokenType=`%v`", token.TokenType)
	// log.Debug("Type=`%v`", token.Type())
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Failed to exchange token: " + err.Error())
	}
	user, err := getUserProfile(c, token)
	if err != nil {
		return err
	}

	json_token, err := json.Marshal(token)
	if err != nil {
		return err
	}

	cookie := &fiber.Cookie{Name: "token", Value: string(json_token), HTTPOnly: true, Expires: time.Now().Add(24 * time.Hour)}
	c.Cookie(cookie)
	log.Debug("%+v\n", user)
	return c.Status(fiber.StatusOK).JSON(user)
	// return c.Redirect(user.Picture)
	// return c.Redirect("/")

}

func getUserProfile(c *fiber.Ctx, token *oauth2.Token) (User, error) {

	client := OAuth_config.Client(context.Background(), token)

	response, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return User{}, c.Status(fiber.StatusBadRequest).SendString("Failed to get user info: " + err.Error())
	}

	defer response.Body.Close()

	var user User
	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return User{}, c.Status(fiber.StatusBadRequest).SendString("Error reading response body: " + err.Error())
	}

	err = json.Unmarshal(bytes, &user)
	if err != nil {
		return User{}, c.Status(fiber.StatusBadRequest).SendString("Error unmarshal json body " + err.Error())
	}

	return user, nil
}

func LogoutHandle(c *fiber.Ctx) error {
	c.ClearCookie("token")
	return c.Redirect("/")
}
