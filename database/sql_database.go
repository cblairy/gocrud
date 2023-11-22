package database

import (
	"database/sql"
	//_ "github.com/go-sql-driver/mysql"
)

func DbConnection() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := "root"
	dbName := "test"
	dbHost := "127.0.0.1:3307"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@tcp("+dbHost+")/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}
