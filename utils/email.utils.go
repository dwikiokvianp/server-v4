package utils

import (
	"gopkg.in/gomail.v2"
	"os"
	"strconv"
)

func SendEmail(to, subject, body, attachment string) error {
	from := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	if attachment != "" {
		m.Attach(attachment)
	}

	d := gomail.NewDialer(smtpHost, smtpPort, from, password)

	err := d.DialAndSend(m)
	if err != nil {
		return err
	}
	return nil
}
