package config

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func DBConnect() (db *sql.DB, err error) {

	dbDriver := "mysql"
	dbUser := "root"
	dbPass := ""
	dbName := "objecttostatistika"

	db, err = sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	return
}