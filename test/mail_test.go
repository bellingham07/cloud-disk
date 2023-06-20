package test

import (
	"github.com/jordan-wright/email"
	"net/smtp"
	"testing"
)

func TestMail(t *testing.T) {
	e := email.NewEmail()
	e.From = "Jordan Wright <test@gmail.com>"
	e.To = []string{"test@example.com"}
	e.Bcc = []string{"test_bcc@example.com"}
	e.Cc = []string{"test_cc@example.com"}
	e.Subject = "Awesome Subject"
	e.Text = []byte("Text Body is, of course, supported!")
	e.HTML = []byte("<h1>Fancy HTML is supported, too!</h1>")
	e.Send("smtp.gmail.com:587", smtp.PlainAuth("", "test@gmail.com", "password123", "smtp.gmail.com"))
}
