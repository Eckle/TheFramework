package utils

import (
	"fmt"
	"runtime"
	"time"
)

const (
	LogTypeInfo    = 0
	LogTypeWarning = 1
	LogTypeError   = 2
)

func Log(message string, error_type int) {
	currentTime := time.Now()
	_, filename, line, _ := runtime.Caller(1)

	final_message := fmt.Sprintf("Type: %d\nTime: %s\nLocation: file '%s' line '%d'\nMessage: %s", error_type, currentTime.Format("2006-01-02 15:04:05.000000000"), filename, line, message)

	println(final_message)
}
