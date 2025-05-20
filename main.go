package main

import (
	api "github.com/NikosGour/google-oauth-example/app"
	"github.com/NikosGour/google-oauth-example/app/auth"
	"github.com/NikosGour/google-oauth-example/storage"
	log "github.com/NikosGour/logging/src"
	"github.com/joho/godotenv"
)

func main() {
	// Init env variables
	dotenv := InitDotenv()

	// Init Database
	mysql_db := storage.NewMySQL_Storage(dotenv["MYSQL_ROOT_PASSWORD"])

	// Init OAuth
	auth.InitOAuthConfig(dotenv)

	// Init Server
	listening_addr := dotenv["HOST_ADDRESS"] + ":" + dotenv["PORT"]
	log.Debug("%s", listening_addr)
	api := api.NewAPIServer(mysql_db, listening_addr, dotenv)

	// Run!
	api.Start()
}

func InitDotenv() map[string]string {
	dotenv, err := godotenv.Read(".env")
	if err != nil {
		log.Fatal("%s", err)
	}
	return dotenv
}
