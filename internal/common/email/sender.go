package email

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"path/filepath"
)

type Sender struct {
	from     string
	password string
	smtpHost string
	smtpPort string
}

func NewSender(from, password, smtpHost, smtpPort string) *Sender {
	return &Sender{
		from:     from,
		password: password,
		smtpHost: smtpHost,
		smtpPort: smtpPort,
	}
}

func (s *Sender) Send(to, subject, body string) error {
	auth := smtp.PlainAuth("", s.from, s.password, s.smtpHost)

	msg := fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n\r\n"+
		"%s\r\n", s.from, to, subject, body)

	err := smtp.SendMail(s.smtpHost+":"+s.smtpPort, auth, s.from, []string{to}, []byte(msg))
	if err != nil {
		return err
	}

	return nil
}

func (s *Sender) SendWithTemplate(to string, templateName string, data interface{}) error {
	templatePath := filepath.Join("internal", "common", "email", "templates", templateName)
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		return err
	}

	var body bytes.Buffer
	if err := t.Execute(&body, data); err != nil {
		return err
	}

	auth := smtp.PlainAuth("", s.from, s.password, s.smtpHost)

	msg := fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"%s\r\n", s.from, to, body.String())

	err = smtp.SendMail(s.smtpHost+":"+s.smtpPort, auth, s.from, []string{to}, []byte(msg))
	if err != nil {
		return err
	}

	return nil
}
