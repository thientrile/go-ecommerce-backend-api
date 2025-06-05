package sendto

import (
	"bytes"
	"fmt"
	"net/smtp"
	"strings"
	"text/template"

	"go-ecommerce-backend-api.com/global"
	"go.uber.org/zap"
)

const (
	SMTPHost = "smtp.gmail.com"
	SMTPPort = "587"
	SMTPUser = "thientrile2003@gmail.com"
	SMTPPass = "tiok zown both uqiq"
)

type EmailAddress struct {
	Address string `json:"address"`
	Name    string `json:"name,omitempty"`
}
type Mail struct {
	From    EmailAddress
	To      []string
	Subject string
	Body    string
}

func BuildMessage(mail Mail) string {
	msg := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\n"
	msg += fmt.Sprintf("From: %s\r\n", mail.From.Address)
	msg += fmt.Sprintf("To: %s\r\n", strings.Join(mail.To, ", "))
	msg += fmt.Sprintf("Subject: %s\r\n", mail.Subject)
	msg += fmt.Sprintf("\r\n%s\r\n", mail.Body)

	return msg
}

func SendTextEmailOtp(to []string, from string, otp string) error {
	contentEmail := Mail{
		From: EmailAddress{
			Address: from,
			Name:    "Ecommerce Backend API",
		},
		To:      to,
		Subject: "Your OTP Code",
		Body:    fmt.Sprintf("Your OTP code is: %s Please enter it to verify your account.", otp),
	}

	messageEmail := BuildMessage(contentEmail)
	// send smtp

	auth := smtp.PlainAuth("", SMTPUser, SMTPPass, SMTPHost)

	err := smtp.SendMail(fmt.Sprintf("%s:%s", SMTPHost, SMTPPort), auth, from, to, []byte(messageEmail))
	if err != nil {
		global.Logger.Error("Email send failed::", zap.Error(err))
		return err
	}
	return nil
}

func SendTemplateEmailOtp(to []string, from string, nameTemplate string, dataTemplate map[string]interface{}) error {

	htmlBody, err := getMailTemplate(nameTemplate, dataTemplate)
	if err != nil {
		return err
	}
	return send(
		"Your OTP Code", to, from, htmlBody)
}

func getMailTemplate(nameTemplate string, dataTemplate map[string]interface{}) (string, error) {
	htmlTemplate := new(bytes.Buffer)
	t := template.Must(template.New(nameTemplate).ParseFiles("templates-email/" + nameTemplate))
	err := t.Execute(htmlTemplate, dataTemplate)
	if err != nil {
		return "", err
	}
	return htmlTemplate.String(), nil
}

func send(Subject string, to []string, from string, htmlTemplate string) error {
	contentEmail := Mail{
		From: EmailAddress{
			Address: from,
			Name:    "Ecommerce Backend API",
		},
		To:      to,
		Subject: Subject,
		Body:    htmlTemplate,
	}

	messageEmail := BuildMessage(contentEmail)
	// send smtp

	auth := smtp.PlainAuth("", SMTPUser, SMTPPass, SMTPHost)

	err := smtp.SendMail(fmt.Sprintf("%s:%s", SMTPHost, SMTPPort), auth, from, to, []byte(messageEmail))
	if err != nil {
		global.Logger.Error("Email send failed::", zap.Error(err))
		return err
	}
	return nil
}
