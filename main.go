package main

import (
	api "github.com/NikosGour/date_management_API/src"
	"github.com/NikosGour/date_management_API/src/storage"
	log "github.com/NikosGour/logging/src"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func main() {
	// Init env variables
	dotenv, err := godotenv.Read(".env")
	if err != nil {
		log.Fatal(err)
	}

	// Init Database
	mysql_db := storage.NewMySQL_Storage(dotenv["MYSQL_ROOT_PASSWORD"])

	// Init OAuth
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

	// Init Server
	listening_addr := dotenv["HOST_ADDRESS"] + ":" + dotenv["PORT"]
	log.Debug("%s", listening_addr)
	api := api.NewAPIServer(mysql_db, listening_addr, oauthConfig, dotenv)

	// Run!
	api.Start()
}
