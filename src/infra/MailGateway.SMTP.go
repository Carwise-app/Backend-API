package infra

import (
	"crypto/tls"
	"net/smtp"
	"os"
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

func (GW *MailGateway) Send(To string, Body []byte) error {
	auth := smtp.PlainAuth("", GW.User, GW.Password, GW.Host)

	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         GW.Host,
	}

	c, err := smtp.Dial(GW.Host + ":" + GW.Port)
	if err != nil {
		return err
	}

	c.StartTLS(tlsconfig)

	if err = c.Auth(auth); err != nil {
		return err
	}

	if err = c.Mail(GW.User); err != nil {
		return err
	}

	if err = c.Rcpt(To); err != nil {
		return err
	}

	w, err := c.Data()
	if err != nil {
		return err
	}

	_, err = w.Write(Body)
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	return nil
}
