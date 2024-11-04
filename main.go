package main

import (
	api "github.com/NikosGour/date_management_API/src"
	"github.com/NikosGour/date_management_API/src/types"
	log "github.com/NikosGour/logging/src"
	"github.com/joho/godotenv"
)

func main() {
	dotenv, err := godotenv.Read(".env")
	if err != nil {
		log.Fatal(err)
	}

	mysql_db := types.NewMySQL_Storage(dotenv["MYSQL_ROOT_PASSWORD"])

	listening_addr := dotenv["HOST_ADDRESS"] + ":" + dotenv["PORT"]
	log.Debug("%s", listening_addr)
	api := api.NewAPIServer(mysql_db, listening_addr, dotenv)
	api.Start()
}
