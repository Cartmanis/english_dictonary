package mailer

import (
	"bytes"
	"crypto/tls"
	"html/template"
	"net"
	"net/smtp"
)

const (
	tplActivate = "activate.msg"
	retrievetpl = "retrieve.msg"
)

type config struct {
	smtphost    string
	user        string
	from        string
	serviceName string
	serverPort  string
	password    string
}

type Message struct {
	From    string
	To      string
	Subject string
	Body    string
	tplName string
}

var cnf config
var tpl *template.Template

func init() {
	const (
		smptpYandex = "smtp.yandex.ru:465"
		userYandex  = "englishDictonary@yandex.ru"
		passYandex  = "5eu7ve"

		smptpGoogle = "smtp.gmail.com:465"
		userGoogle  = "VShmelcer@gmail.com"
		passGoogle  = "5eu7ve5eu7ve"
	)
	cnf = config{
		smtphost:    smptpYandex,
		user:        userYandex,
		password:    passYandex,
		from:        userYandex,
		serviceName: "Dictonary",
		serverPort:  "27111",
	}
	tpl = template.Must(template.New("").ParseGlob("cmd/mailer/templates/*.msg"))

}

func (m *Message) getMailMessage() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := tpl.ExecuteTemplate(buf, m.tplName, m); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func getSMTPClient() (*smtp.Client, error) {
	host, _, _ := net.SplitHostPort(cnf.smtphost)

	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	conn, err := tls.Dial("tcp", cnf.smtphost, tlsconfig)
	if err != nil {
		return nil, err
	}
	client, err := smtp.NewClient(conn, cnf.smtphost)
	if err != nil {
		return nil, err
	}

	auth := smtp.PlainAuth("", cnf.user, cnf.password, cnf.smtphost)

	if err := client.Auth(auth); err != nil {
		return nil, err
	}
	return client, nil
}

func messageLoop(m Message) error {
	client, err := getSMTPClient()
	if err != nil {
		return err
	}
	defer client.Close()
	if err := client.Mail(cnf.user); err != nil {
		return err
	}
	if err := client.Rcpt(m.To); err != nil {
		return err
	}
	writeCloser, err := client.Data()
	if err != nil {
		return err
	}
	data, err := m.getMailMessage()
	if err != nil {
		return err
	}
	if _, err = writeCloser.Write(data); err != nil {
		return err
	}
	if err := writeCloser.Close(); err != nil {
		return err
	}
	return nil
}
