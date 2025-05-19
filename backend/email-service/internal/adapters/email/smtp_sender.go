package email

import (
	"fmt"
	"net/smtp"
	"strings"
)

type SMTPSender struct {
	host string
	auth smtp.Auth
	from string
}

type NewSMTPSenderFunc func(host, user, pass, from string) *SMTPSender

func NewSMTPSender(host, user, pass, from string) *SMTPSender {
	auth := smtp.PlainAuth("", user, pass, host[:strings.Index(host, ":")])
	return &SMTPSender{host: host, auth: auth, from: from}
}

type SMTPEmailSender interface {
	SendConfirmation(to, link string) error
	SendNotification(to, city, temp, cond, unsubURL string) error
}

func (s *SMTPSender) SendConfirmation(to, link string) error {
	subj := "Confirm your weather subscription"
	body := fmt.Sprintf("%s", link)
	return s.sendMail(to, subj, body)
}

func (s *SMTPSender) SendNotification(to, city, temp, cond, unsubURL string) error {
	subj := fmt.Sprintf("Weather update for %s", city)
	body := fmt.Sprintf("%s: %s°C, %s. Unsubscribe: %s", city, temp, cond, unsubURL)
	return s.sendMail(to, subj, body)
}

func (s *SMTPSender) sendMail(to, subj, body string) error {
	msg := []byte(fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", s.from, to, subj, body))
	return smtp.SendMail(s.host, s.auth, s.from, []string{to}, msg)
}
