package infra

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net/smtp"
	"os"
	"text/template"
)

type MailGateway struct {
	Host     string
	Port     string
	User     string
	Password string
}

func NewMailGateway() *MailGateway {
	return &MailGateway{
		Host:     os.Getenv("SMTP_HOST"),
		Port:     os.Getenv("SMTP_PORT"),
		User:     os.Getenv("SMTP_USER"),
		Password: os.Getenv("SMTP_PASSWORD"),
	}
}

func (gw *MailGateway) SendEmail(to, subject, templateName string, data map[string]interface{}) error {
	// Şablonu oku ve işle
	tmpl, err := template.ParseFiles("templates/" + templateName)
	if err != nil {
		return fmt.Errorf("template parse error: %v", err)
	}

	var bodyBuffer bytes.Buffer
	if err := tmpl.Execute(&bodyBuffer, data); err != nil {
		return fmt.Errorf("template execute error: %v", err)
	}

	// Mail mesajı oluştur
	msg := bytes.Buffer{}
	msg.WriteString(fmt.Sprintf("From: Carwise <%s>\r\n", gw.User))
	msg.WriteString(fmt.Sprintf("To: %s\r\n", to))
	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	msg.WriteString("MIME-Version: 1.0\r\n")
	msg.WriteString("Content-Type: text/html; charset=\"UTF-8\"\r\n")
	msg.WriteString("\r\n")
	msg.Write(bodyBuffer.Bytes())

	// SMTP kimlik doğrulaması
	auth := smtp.PlainAuth("", gw.User, gw.Password, gw.Host)

	// TLS yapılandırması
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         gw.Host,
	}

	// SMTP sunucusuna bağlan
	c, err := smtp.Dial(gw.Host + ":" + gw.Port)
	if err != nil {
		return fmt.Errorf("SMTP dial error: %v", err)
	}
	defer c.Close()

	if err = c.StartTLS(tlsconfig); err != nil {
		return fmt.Errorf("StartTLS error: %v", err)
	}

	if err = c.Auth(auth); err != nil {
		return fmt.Errorf("SMTP auth error: %v", err)
	}

	if err = c.Mail(gw.User); err != nil {
		return fmt.Errorf("MAIL FROM error: %v", err)
	}

	if err = c.Rcpt(to); err != nil {
		return fmt.Errorf("RCPT TO error: %v", err)
	}

	w, err := c.Data()
	if err != nil {
		return fmt.Errorf("DATA command error: %v", err)
	}

	if _, err = w.Write(msg.Bytes()); err != nil {
		return fmt.Errorf("message write error: %v", err)
	}

	if err = w.Close(); err != nil {
		return fmt.Errorf("close write error: %v", err)
	}

	return nil
}
