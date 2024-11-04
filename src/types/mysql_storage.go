package types

import (
	"database/sql"
	"fmt"

	log "github.com/NikosGour/logging/src"
	_ "github.com/go-sql-driver/mysql"
)

type MySQL_Storage struct {
	db *sql.DB
}

func NewMySQL_Storage(dbpass string) *MySQL_Storage {
	this := &MySQL_Storage{}
	this.init_database(dbpass)
	return this
}

func (this *MySQL_Storage) init_database(dbpass string) {
	fmt.Printf("dbpass = %s\n", dbpass)
	conn_string := fmt.Sprintf("root:%s@(127.0.0.1:3306)/main_db", dbpass)

	var err error
	this.db, err = sql.Open("mysql", conn_string)
	if err != nil {
		log.Fatal(err)
	}

	if err = this.db.Ping(); err != nil {
		log.Fatal(err)
	}

	rows, err := this.db.Query("SHOW TABLES;")
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			log.Fatal(err)
		}

		log.Debug("Table = %s\n", table)
	}

}
