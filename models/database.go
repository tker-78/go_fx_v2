package models

import (
	"database/sql"
	"fmt"
	"log"

	"example.com/tker-78/fx2/config"
	_ "github.com/lib/pq"
)

var DbConnection *sql.DB

func init() {
	var err error
	connectionStr := "user=takuyakinoshita dbname=exchange_2 sslmode=disable"
	DbConnection, err := sql.Open(config.Config.SQLDriver, connectionStr)
	if err != nil {
		log.Fatalln("Error occured while opening database file: ", err)
	}

	tableName := "USD_JPY_1d"
	cmd := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			time TIMESTAMP PRIMARY KEY NOT NULL,
			open FLOAT,
			high FLOAT,
			low FLOAT,
			close FLOAT,
			swap FLOAT)
	`, tableName)

	_, err = DbConnection.Exec(cmd)
	if err != nil {
		log.Fatalln("Error occuredd while creating database table: ", err)
	}
}
