package main

import (
	"log"
	"net/smtp"
)

// SMTP server configuration.
const (
	smtpHost = "smtp.gmail.com"
	smtpPort = "587"
	smtpUser = "yourusername@gmail.com"
	smtpPass = "yourapppassword"
)

func sendEmail(to, subject, body string) error {
	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)

	message := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, smtpUser, []string{to}, message)
	if err != nil {
		return err
	}

	log.Println("Email sent successfully")
	return nil
}
