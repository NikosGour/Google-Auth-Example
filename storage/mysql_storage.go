package storage

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
	conn_string := fmt.Sprintf("root:%s@(127.0.0.1:3306)/main_db", dbpass)

	var err error
	this.db, err = sql.Open("mysql", conn_string)
	if err != nil {
		log.Fatal("%s", err)
	}

	if err = this.db.Ping(); err != nil {
		log.Fatal("%s", err)
	}

	err = this.CreateTables()
	if err != nil {
		log.Fatal("%s", err)
	}
}

func (this *MySQL_Storage) CreateTables() error {
	err := this.createDatesTable()
	if err != nil {
		return err
	}
	log.Info("Database: all tables created")
	return nil
}

func (this *MySQL_Storage) createDatesTable() error {
	query := `create table if not exists dates
			(
    			id         int auto_increment primary key,
    			start_date timestamp not null,
    			end_date   timestamp not null,
    			CONSTRAINT start_date_end_date CHECK (start_date < end_date)
			);`
	_, err := this.db.Exec(query)
	if err != nil {
		return fmt.Errorf("createDatesTable: %s", err)
	}

	log.Info("Database: dates table created")
	return nil
}
