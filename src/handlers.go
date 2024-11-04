package api

import (
	"context"
	"encoding/json"
	"io"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func (this *APIServer) authGoogleHandle(c *fiber.Ctx) error {
	url := this.OAuth_config.AuthCodeURL("state")
	return c.Redirect(url)
}
func (this *APIServer) authRedirectHandle(c *fiber.Ctx) error {
	code := c.Query("code") //get code from query params for generating token
	if code == "" {
		return c.Status(fiber.StatusBadRequest).SendString("")
	}
	token, err := this.OAuth_config.Exchange(context.Background(), code) //get token
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Failed to exchange token: " + err.Error())
	}
	client := this.OAuth_config.Client(context.Background(), token)              //set client for getting user info like email, name, etc.
	response, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo") //get user info
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Failed to get user info: " + err.Error())
	}

	defer response.Body.Close()
	var user User                           //user variable
	bytes, err := io.ReadAll(response.Body) //reading response body from client
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Error reading response body: " + err.Error())
	}
	err = json.Unmarshal(bytes, &user) //unmarshal user info
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Error unmarshal json body " + err.Error())
	}
	log.Debug("%+v", user)
	// return c.Status(fiber.StatusOK).JSON(user) //return user info
	return c.Redirect(user.Picture)

}
