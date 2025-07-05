package main

import (
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/smtp"
	"os"
	"fmt"
)

// HashPassword hashes a plain password.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), err
}

// CheckPasswordHash compares a plain password with a hash.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// SendEmail sends an email using SMTP.
func SendEmail(to, subject, body string) error {
	from := os.Getenv("SMTP_USER")
	pass := os.Getenv("SMTP_PASSWORD")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n\n" +
		body

	auth := smtp.PlainAuth("", from, pass, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(msg))
	if err != nil {
		log.Println("Failed to send email:", err)
	}
	return err
}

// Helper to format order details for email
func OrderDetails(order Order, coffeeName string) string {
	return fmt.Sprintf("Name: %s\nCoffee: %s\nQuantity: %d\nNotes: %s\nEmail: %s",
		order.Name, coffeeName, order.Quantity, order.Notes, order.Email)
}