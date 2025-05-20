package auth_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"

	log "github.com/NikosGour/logging/src"

	api "github.com/NikosGour/google-oauth-example/app"
	"github.com/NikosGour/google-oauth-example/app/auth"

	"github.com/gofiber/fiber/v2"

	"github.com/stretchr/testify/assert"

	"github.com/joho/godotenv"
)

const (
	env_file = "../../testing.env"
)

var (
	dotenv map[string]string
	app    *fiber.App
)

func TestMain(m *testing.M) {
	_dotenv, err := godotenv.Read(env_file)
	dotenv = _dotenv
	if err != nil {
		log.Fatal("%s", err)
	}
	auth.InitOAuthConfig(dotenv)
	app = api.SetupFiberApp()

	exit_code := m.Run()
	os.Exit(exit_code)

}
func TestAuthenticateUser(t *testing.T) {
	// Init
	assert := assert.New(t)
	refresh_token, ok := dotenv["GOOGLE_REFRESH_TOKEN"]
	if !ok {
		t.Fatalf("`GOOGLE_REFRESH_TOKEN` not found in env file: `%s`", env_file)
	}

	success_response := "Success"

	test_app := fiber.New()
	test_app.Get("/", auth.AuthenticateUser, func(c *fiber.Ctx) error {
		return c.SendString(success_response)
	})

	date_now_plus_one_day_json := time.Now().Add(time.Hour * 24)

	// Test Cases
	test_cases := []struct {
		Name                     string
		CookieJSON               any
		StatusCode               int
		ShouldGetSuccessResponse bool
	}{
		{
			Name: "Will refresh but still passes",
			CookieJSON: map[string]any{
				"access_token":  "aksjdflkajsdfklalsdfjkajkldfklasdjkfjaksljdfklajkfljskld",
				"token_type":    "Bearer",
				"refresh_token": refresh_token,
				"expiry":        date_now_plus_one_day_json,
			},
			StatusCode:               fiber.StatusOK,
			ShouldGetSuccessResponse: true,
		},
		{
			Name: "No access token",
			CookieJSON: map[string]any{
				"token_type":    "Bearer",
				"refresh_token": refresh_token,
				"expiry":        "2025-05-14T22:51:37.536435719+03:00",
			},
			StatusCode:               fiber.StatusOK,
			ShouldGetSuccessResponse: true,
		},
		{
			Name: "No refresh token",
			CookieJSON: map[string]any{
				"token_type": "Bearer",
				"expiry":     "2025-05-14T22:51:37.536435719+03:00",
			},
			StatusCode:               fiber.StatusFound,
			ShouldGetSuccessResponse: false,
		},
		{
			Name:                     "Empty cookie",
			CookieJSON:               "",
			StatusCode:               fiber.StatusUnauthorized,
			ShouldGetSuccessResponse: false,
		},
		{
			Name:                     "Valid json mangled cookie",
			CookieJSON:               `{"not a":"valid token cookie"}`,
			StatusCode:               fiber.StatusBadRequest,
			ShouldGetSuccessResponse: false,
		},
		{
			Name:                     "Invalid json cookie",
			CookieJSON:               `{"not valid token cookie"}`,
			StatusCode:               fiber.StatusBadRequest,
			ShouldGetSuccessResponse: false,
		},
	}

	// Subtests
	for _, test := range test_cases {
		t.Run(test.Name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)

			log.Debug("test.CookieJSON=%s", test.CookieJSON)
			switch v := test.CookieJSON.(type) {
			case string:
				if v != "" {
					str := url.QueryEscape(v)
					cookie := &http.Cookie{Name: "token", Value: str, HttpOnly: true}
					req.AddCookie(cookie)
				}

			case map[string]any:
				json, err := json.Marshal(v)
				assert.NoError(err, test.Name)

				str := url.QueryEscape(string(json))
				cookie := &http.Cookie{Name: "token", Value: str, HttpOnly: true}
				req.AddCookie(cookie)
			}

			res, err := test_app.Test(req, -1)
			assert.NoError(err, test.Name)

			assert.Equal(test.StatusCode, res.StatusCode, test.Name)

			if test.ShouldGetSuccessResponse {
				body, err := io.ReadAll(res.Body)
				assert.NoError(err, test.Name)
				assert.Equal(success_response, string(body), test.Name)
			}
		})
	}
}
