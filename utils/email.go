package utils

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"

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

// 👇 Email template parser
func ParseTemplateDir(dir string) (*template.Template, error) {
	var paths []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			paths = append(paths, path)
		}
		return nil
	})

	fmt.Println("Am parsing templates...")

	if err != nil {
		return nil, err
	}

	return template.ParseFiles(paths...)
}

func SendEmail(user *models.DBResponse, data *EmailData,  templateName string) error {
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

	template, err := ParseTemplateDir("templates")
	if err != nil {
		log.Fatal("Could not parse template", err)
	}

	template = template.Lookup(templateName)
	template.Execute(&body, &data)
	fmt.Println(template.Name())
	

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

