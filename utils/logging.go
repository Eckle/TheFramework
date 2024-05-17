package utils

import (
	"fmt"
	"runtime"
	"time"

	"github.com/Eckle/TheFramework/db"
)

const (
	LogTypeInfo    = 0
	LogTypeWarning = 1
	LogTypeError   = 2
)

func InitLogger() error {
	if db.Database == nil {
		Log("The database has not been initialized. Starting Logger in print only mode.", LogTypeInfo)
		return nil
	}

	migrations_table := `
		CREATE TABLE IF NOT EXISTS logs (
			id INTEGER PRIMARY KEY,
			message TEXT,
			filename TEXT,
			line INTEGER,
			time DATETIME DEFAULT CURRENT_TIMESTAMP
		);
	`

	_, err := db.Database.Exec(migrations_table)
	if err != nil {
		Log(fmt.Sprintf("Could not create logs table due to an error. Starting Logger in print only mode. Error: %v", err), LogTypeError)
		return err
	}

	return nil
}

func Log(message string, error_type int) {
	currentTime := time.Now()
	_, filename, line, _ := runtime.Caller(1)

	final_message := fmt.Sprintf("Type: %d\nTime: %s\nLocation: file '%s' line '%d'\nMessage: %s", error_type, currentTime.Format("2006-01-02 15:04:05.000000000"), filename, line, message)

	if db.Database != nil {
		_, err := db.Database.Exec(db.Database.Rebind("INSERT INTO logs (message, filename, line) VALUES (?, ?, ?)"), message, filename, line)
		if err != nil {
			println("Oops. Error happened in the Logger.")
		}
	}

	println(final_message)
}
