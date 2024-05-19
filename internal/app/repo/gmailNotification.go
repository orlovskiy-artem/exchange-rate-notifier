package repo

import (
	"context"
	"fmt"
	"log"
	"net/smtp"

	"genesis_test_task/internal/app/model"
)

type GmailNotificationRepo struct {
	adressorEmail    model.Email
	adressorPassword string
	smtpHost         string
	smtpPort         string
	log              *log.Logger
}

func NewGmailNotificationRepo(
	adressorEmail model.Email,
	adressorPassword string,
	smtpHost string,
	smtpPort string,
	logger *log.Logger) *GmailNotificationRepo {
	return &GmailNotificationRepo{
		adressorEmail:    adressorEmail,
		adressorPassword: adressorPassword,
		smtpHost:         smtpHost,
		smtpPort:         smtpPort,
		log:              logger}
}

func (nr GmailNotificationRepo) SendEmailMessage(addressee model.Email, text string) error {
	message := []byte(text)
	// Authentication.
	auth := smtp.PlainAuth("",
		nr.adressorEmail.String(),
		nr.adressorPassword,
		nr.smtpHost)
	// Sending email.
	err := smtp.SendMail(
		nr.smtpHost+":"+nr.smtpPort,
		auth,
		nr.adressorEmail.String(),
		[]string{addressee.String()},
		message)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("Email Sent Successfully!")
	return nil
}

func (nr GmailNotificationRepo) SendEmailExchangeRate(ctx context.Context, addressee model.Email, exchangeRate model.ExchangeRate) error {
	message := fmt.Sprintf("To: %s\r\nSubject: Exchange Rate UAH-USD\r\n\r\nExchange Rate: %f", addressee.String(), exchangeRate.Value)
	nr.SendEmailMessage(addressee, message)
	return nil
}
