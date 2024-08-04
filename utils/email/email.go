package email

import (
	"fmt"
	"net/smtp"
	"os"
)

var auth smtp.Auth
var server string

func Init() {
	server = fmt.Sprintf("%s:%s", os.Getenv("SMTP_SERVER"), os.Getenv("SMTP_PORT"))
	auth = smtp.PlainAuth("", os.Getenv("SMTP_USER"), os.Getenv("SMTP_PASS"), os.Getenv("SMTP_SERVER"))
}

func SendEmail(to string, subject string, msg string) error {
	msg = fmt.Sprintf("\r\n\r\n%s", msg)
	msg = addHeader("To", to, msg)
	msg = addHeader("Subject", subject, msg)
	msg = addHeader("Content-Type", "text/html; charset=UTF-8", msg)

	err := smtp.SendMail(server, auth, os.Getenv("SMTP_USER"), []string{to}, []byte(msg))
	if err != nil {
		return err
	}
	return nil
}

func addHeader(key string, value string, request string) string {
	return fmt.Sprintf("%s: %s\r\n%s", key, value, request)
}
