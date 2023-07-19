package models

import (
	"database/sql"

	_ "github.com/lib/pq"
)

var DbConnection *sql.DB

func init() {

}
