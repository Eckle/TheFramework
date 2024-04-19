package db

import (
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

var Database *sqlx.DB

func InitDB() {
	db, err := sqlx.Open("libsql", os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(1)
	Database = db
}

func Cleanup() {
	Database.Close()
}
