package email

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"image"
	"image/png"
	"net/smtp"
	"os"
	"path/filepath"
	"time"

	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
)

type Sender struct {
	from     string
	password string
	smtpHost string
	smtpPort string
}

func NewSender(smtpHost, smtpPort, from, password string) *Sender {
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
		"Subject: %s\r\n"+
		"%s", s.from, to, subject, body)

	smtpAddr := fmt.Sprintf("%s:%s", s.smtpHost, s.smtpPort)
	err := smtp.SendMail(smtpAddr, auth, s.from, []string{to}, []byte(msg))
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}

func (s *Sender) SendWithTemplate(to, subject, templateName string, data interface{}) error {
	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %v", err)
	}

	// Construct the absolute path to the template file
	templatePath := filepath.Join(cwd, "templates", "email", templateName)

	// Parse the template file
	t, err := template.New(templateName).Funcs(template.FuncMap{
		"formatDate": func(date time.Time) string {
			return date.Format("2006-01-02")
		},
		"printf": fmt.Sprintf,
	}).ParseFiles(templatePath)
	if err != nil {
		return fmt.Errorf("failed to parse template: %v", err)
	}

	// Read Stori logo SVG
	logoPath := filepath.Join(cwd, "web", "static", "images", "stori-logo.svg")
	logoBytes, err := os.ReadFile(logoPath)
	if err != nil {
		return fmt.Errorf("failed to read Stori logo: %v", err)
	}

	// Convert SVG to PNG
	icon, _ := oksvg.ReadIconStream(bytes.NewReader(logoBytes))
	w, h := int(icon.ViewBox.W), int(icon.ViewBox.H)
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	scanner := rasterx.NewScannerGV(w, h, img, img.Bounds())
	raster := rasterx.NewDasher(w, h, scanner)
	icon.Draw(raster, 1.0)

	// Encode PNG to base64
	var buf bytes.Buffer
	png.Encode(&buf, img)
	logoBase64 := base64.StdEncoding.EncodeToString(buf.Bytes())

	// Prepare data for template
	templateData := struct {
		StoriLogo string
		Data      interface{}
	}{
		StoriLogo: fmt.Sprintf("data:image/png;base64,%s", logoBase64),
		Data:      data,
	}

	var body bytes.Buffer
	if err := t.Execute(&body, templateData); err != nil {
		return fmt.Errorf("failed to execute template: %v", err)
	}

	// Construct the email content
	emailContent := fmt.Sprintf("MIME-Version: 1.0\r\n"+
		"Content-Type: text/html; charset=\"utf-8\"\r\n"+
		"\r\n%s", body.String())

	// Use the Send function to send the email
	err = s.Send(to, subject, emailContent)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	// // Set up email headers
	// headers := make(map[string]string)
	// headers["From"] = s.from
	// headers["To"] = to
	// headers["Subject"] = subject
	// headers["MIME-Version"] = "1.0"
	// headers["Content-Type"] = "text/html; charset=\"utf-8\""

	// // Construct message
	// message := ""
	// for k, v := range headers {
	// 	message += fmt.Sprintf("%s: %s\r\n", k, v)
	// }
	// message += "\r\n" + body.String()

	// // Configurar la dirección del servidor SMTP
	// smtpAddr := fmt.Sprintf("%s:%s", s.smtpHost, s.smtpPort)

	// // Configurar la autenticación
	// auth := smtp.PlainAuth("", s.from, s.password, s.smtpHost)

	// // Enviar el email
	// err = smtp.SendMail(smtpAddr, auth, s.from, []string{to}, []byte(message))
	// if err != nil {
	// 	return fmt.Errorf("failed to send email: %v", err)
	// }
	return nil
}
