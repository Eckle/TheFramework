package db

import (
	"os"

	_ "github.com/glebarez/go-sqlite"
	"github.com/jmoiron/sqlx"
)

var Database *sqlx.DB = nil

func InitDB() {
	db, err := sqlx.Open("sqlite", os.Getenv("DATABASE_URL"))

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
