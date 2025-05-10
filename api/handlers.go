package api

import (
	"context"
	"encoding/json"
	"io"
	"time"

	log "github.com/NikosGour/logging/src"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
)

func authGoogleHandle(c *fiber.Ctx) error {
	url := OAuth_config.AuthCodeURL("nikos", oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	return c.Redirect(url)
}
func authRedirectHandle(c *fiber.Ctx) error {
	code := c.Query("code")
	if code == "" {
		return c.Status(fiber.StatusBadRequest).SendString("")
	}

	token, err := OAuth_config.Exchange(context.Background(), code) //get token
	log.Debug("AccessToken=`%v`", token.AccessToken)
	log.Debug("RefreshToken=`%v`", token.RefreshToken)
	log.Debug("Expiry=`%v`", token.Expiry)
	log.Debug("IsValid=`%v`", token.Valid())
	log.Debug("TokenType=`%v`", token.TokenType)
	log.Debug("Type=`%v`", token.Type())
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

	client := OAuth_config.Client(context.Background(), token) //set client for getting user info like email, name, etc.

	response, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo") //get user info
	if err != nil {
		return User{}, c.Status(fiber.StatusBadRequest).SendString("Failed to get user info: " + err.Error())
	}

	defer response.Body.Close()

	var user User                           //user variable
	bytes, err := io.ReadAll(response.Body) //reading response body from client
	if err != nil {
		return User{}, c.Status(fiber.StatusBadRequest).SendString("Error reading response body: " + err.Error())
	}

	err = json.Unmarshal(bytes, &user) //unmarshal user info
	if err != nil {
		return User{}, c.Status(fiber.StatusBadRequest).SendString("Error unmarshal json body " + err.Error())
	}

	return user, nil
}
