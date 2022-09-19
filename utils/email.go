package utils

import (
	"bytes"
	"crypto/tls"
	"html/template"
	"log"

	"github.com/amrilsyaifa/go_mongodb_rest_api/config"
	"github.com/amrilsyaifa/go_mongodb_rest_api/models"
	"github.com/k3a/html2text"
	"gopkg.in/gomail.v2"
)

type EmailData struct {
	URL			string
	FirstName	string
	Subject		string
}

func SendEmail(user *models.DBResponse, data *EmailData, template *template.Template, templateName string) error {
	config, err := config.LoadConfig(".")

	if err != nil {
		log.Fatal("could not load config", err)
	}

	// Sender data.
	from := config.EmailFrom
	smtpPass := config.SMTPPass
	smtpUser := config.SMTPUser
	to := user.Email
	smtpHost := config.SMTPHost
	smtpPort := config.SMTPPort

	var body bytes.Buffer

	if err := template.ExecuteTemplate(&body, templateName, &data); err != nil {
		log.Fatal("Could not execute template", err)
	}

	message := gomail.NewMessage()

	message.SetHeader("From", from)
	message.SetHeader("To", to)
	message.SetHeader("Subject", data.Subject)
	message.SetHeader("text/html",body.String())
	message.AddAlternative("text/plain", html2text.HTML2Text(body.String()))

	dialer := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Send Email
	if err := dialer.DialAndSend(message); err != nil {
		return err
	}

	return nil
}