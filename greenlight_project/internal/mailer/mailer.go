package mailer

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"
	"time"

	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
)

// var templateFS embed.FS

type Dialer struct {
	auth   sasl.Client
	server string
}

type Mailer struct {
	dialer Dialer
	sender string
}

func New(host string, port int, username, password, sender string) Mailer {
	auth := sasl.NewPlainClient("", username, password)
	server := fmt.Sprintf("%s:%d", host, port)
	return Mailer{
		dialer: Dialer{
			auth:   auth,
			server: server,
		},
		sender: sender,
	}
}

func (m Mailer) Send(recipient, templateFile string, data interface{}) error {
	tmpl, err := template.ParseFiles(templateFile)
	// tmpl, err := template.New("email").ParseFS(templateFS, "templates/"+templateFile)
	if err != nil {
		return err
	}
	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}

	// plainBody := new(bytes.Buffer)
	// err = tmpl.ExecuteTemplate(plainBody, "plainBody", data)
	// if err != nil {
	// 	return err
	// }
	//
	htmlBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(htmlBody, "htmlBody", data)
	if err != nil {
		return err
	}

	to := []string{recipient}
	from := m.sender
	msg := strings.NewReader(fmt.Sprintf("To: %s\r\nSubject: %s\r\nContent-Type: text/html; charset=UTF-8\r\n%s\r\n",
		to, subject.String(), htmlBody.String()))

	for i := 1; i <= 3; i++ {
		err = smtp.SendMail(m.dialer.server, m.dialer.auth, from, to, msg)
		if err == nil {
			return nil
		}
		time.Sleep(500 * time.Millisecond)
	}

	return err
}
