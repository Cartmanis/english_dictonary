package smtp

import "html/template"

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
	cnf := config{
		smtphost:    "smtp.yandex.ru:465",
		user:        "Cartmanis@yandex.ru",
		from:        "Cartmanis@yandex.ru",
		serviceName: "english_dictonary",
		serverPort:  "27330",
		password:    "5eu7ve",
	}
	tpl = template.Must(template.New("").ParseGlob("./templates/mail/*.msg"))

	queue = make(chan Message, 10)
}
