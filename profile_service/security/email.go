package security

import (
	"crypto/tls"
	"gopkg.in/gomail.v2"
	"log"
	"profile_service/model"
)

type EmailData struct {
	URL       string
	FirstName string
	Subject   string
}

// 👇 Email template parser

//func ParseTemplateDir(dir string) (*template.Template, error) {
//	var paths []string
//	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
//		if err != nil {
//			return err
//		}
//		if !info.IsDir() {
//			paths = append(paths, path)
//		}
//		return nil
//	})
//
//	if err != nil {
//		return nil, err
//	}
//
//	return template.ParseFiles(paths...)
//}

func SendEmail(user *model.RegularProfile, data *EmailData) bool {

	// Sender data.
	from := "ibsit2022@gmail.com"
	smtpPass := "xeuloaiprwagrouh"
	smtpUser := "ibsit2022@gmail.com"
	to := user.Email
	smtpHost := "smtp.gmail.com"
	smtpPort := 587

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
		log.Fatal("Could not send email: ", err)
		return false
	}

	return true
}

func SendEmailBusiness(user *model.BusinessProfile, data *EmailData) bool {

	// Sender data.
	from := "ibsit2022@gmail.com"
	smtpPass := "xeuloaiprwagrouh"
	smtpUser := "ibsit2022@gmail.com"
	to := user.Email
	smtpHost := "smtp.gmail.com"
	smtpPort := 587

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
		log.Fatal("Could not send email: ", err)
		return false
	}

	return true
}

//func SendMailSMTP(code string, name string) bool {
//
//	from := "ibsit2022@gmail.com"
//	password := "xeuloaiprwagrouh"
//
//	toEmailAddress := "ibsit2022@gmail.com"
//	to := []string{toEmailAddress}
//
//	host := "smtp.gmail.com"
//	port := "587"
//	address := host + ":" + port
//
//	subject := "Account verification for " + name + "\n"
//	body := "Your account verification code is " + code
//	message := []byte(subject + body)
//
//	auth := smtp.PlainAuth("", from, password, host)
//
//	err := smtp.SendMail(address, auth, from, to, message)
//	if err != nil {
//		panic(err)
//		return false
//	}
//
//	return true
//}
