package utils

import (
	"log"
	"net/smtp"
)

// SendEmail sends
func SendEmail(body string, subject string) {
	from := GetEnvVariable("EMAIL_USER")
	pass := GetEnvVariable("APP_PASSWORD")
	to := GetEnvVariable("EMAIL_USER")

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n\n" +
		body

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return
	}

	log.Print("Email sent")
}
