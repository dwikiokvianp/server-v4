package utils

import (
	"encoding/base64"
	"net/smtp"
	"os"
)

func SendPDFEmail(email, subject, body string, pdfData []byte) error {
	from := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	auth := smtp.PlainAuth("", from, password, smtpHost)

	msg := "From: " + from + "\n" +
		"To: " + email + "\n" +
		"Subject: " + subject + "\n" +
		"MIME-Version: 1.0" + "\n" +
		"Content-Type: multipart/mixed; boundary=boundary" + "\n\n" +
		"--boundary" + "\n" +
		"Content-Type: text/html; charset=utf-8" + "\n" +
		"Content-Transfer-Encoding: 7bit" + "\n\n" +
		body + "\n\n" +
		"--boundary" + "\n" +
		"Content-Type: application/pdf" + "\n" +
		"Content-Disposition: attachment; filename=invoice.pdf" + "\n" +
		"Content-Transfer-Encoding: base64" + "\n\n" +
		base64.StdEncoding.EncodeToString(pdfData) + "\n\n" +
		"--boundary--"

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{email}, []byte(msg))
	if err != nil {
		return err
	}

	return nil
}
