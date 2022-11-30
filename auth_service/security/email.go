package security

import (
	"auth-service/model"
	"crypto/tls"
	"html/template"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/gomail.v2"
)

type EmailData struct {
	URL       string
	FirstName string
	Subject   string
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

	if err != nil {
		return nil, err
	}

	return template.ParseFiles(paths...)
}

func SendEmail(user *model.RegularProfile, data *EmailData) {

	// Sender data.
	from := "ibsit2022@gmail.com"
	smtpPass := "xeuloaiprwagrouh"
	smtpUser := "ibsit2022@gmail.com"
	to := user.Email
	smtpHost := "smtp.gmail.com"
	smtpPort := 587

	//var body bytes.Buffer

	//template, err := ParseTemplateDir("templates")
	//if err != nil {
	//	log.Fatal("Could not parse template", err)
	//}
	//
	//template.ExecuteTemplate(&body, "verificationCode.html", &data)

	//parsedTemplate, err := template.ParseFiles("auth_service/templates/verificationCode.html")
	//parsedTemplate.Execute(&body, &data)
	//if err != nil {
	//	log.Println("Error executing template :", err)
	//	return
	//}

	m := gomail.NewMessage()

	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", data.Subject)
	m.SetBody("text/html", data.URL)
	m.AddAlternative("text/plain", data.URL)

	d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Send Email
	if err := d.DialAndSend(m); err != nil {
		log.Println("---- GRESKA ---- ", err)
		log.Fatal("Could not send email: ", err)
	}

}
