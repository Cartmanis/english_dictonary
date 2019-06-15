package mailer

import (
	"bytes"
	"fmt"
	"html/template"
)

const (
	passtpl     = "password.msg"
	retrievetpl = "retrieve.msg"
)

type Message struct {
	From    string
	To      string
	Body    string
	tplName string
}

func (m *Message) getMailMessage() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := tpl.ExecuteTemplate(buf, m.tplName, m); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

type config struct {
	smtphost    string
	user        string
	from        string
	serviceName string
	serverPort  string
	password    string
}

var cnf config
var tpl *template.Template
var queue chan Message

func init() {
	cnf = config{
		smtphost:    "smtp.yandex.ru:465",
		user:        "Cartmanis@yandex.ru",
		from:        "Cartmanis@yandex.ru",
		serviceName: "english_dictonary",
		serverPort:  "27111",
		password:    "5eu7ve",
	}
	fmt.Print(cnf)
	tpl = template.Must(template.New("").ParseGlob("./templates/mail/*.msg"))

	//буферезированный канал. Максимальная очередь из сообщений в количестве 10
	queue = make(chan Message, 10)
}
