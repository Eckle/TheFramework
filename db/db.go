package db

import (
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

var Database *sqlx.DB = nil

type DatabaseDriver string

const (
	Postgres DatabaseDriver = "postgres"
	MySQL    DatabaseDriver = "mysql"
	Turso    DatabaseDriver = "libsql"
)

func InitDB(driver DatabaseDriver) {
	db, err := sqlx.Open(string(driver), os.Getenv("DATABASE_URL"))

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
