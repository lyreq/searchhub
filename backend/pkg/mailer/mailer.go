package mailer

import (
	"crypto/tls"
	"lytemp/config"

	"gopkg.in/gomail.v2"
)

type GoMailer struct {
	dialer *gomail.Dialer
	from   string
	name   string
}

func NewGoMailer(cfg config.Mail) *GoMailer {
	dialer := gomail.NewDialer(cfg.Host, cfg.Port, cfg.Username, cfg.Password)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: !cfg.SSL}

	return &GoMailer{
		dialer: dialer,
		from:   cfg.FromEmail,
		name:   cfg.FromName,
	}
}

func (g *GoMailer) Send(to []string, subject string, body string, isHTML bool) error {
	m := gomail.NewMessage()
	m.SetAddressHeader("From", g.from, g.name)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)

	if isHTML {
		m.SetBody("text/html", body)
	} else {
		m.SetBody("text/plain", body)
	}

	return g.dialer.DialAndSend(m)
}
