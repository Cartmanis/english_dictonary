package service

import (
	"fmt"
	"github.com/cartmanis/english_dictonary/backend/app/provider_db"
	"math/rand"
	"net/smtp"
	"time"
)

type word struct {
	En            string
	Ru            string
	Transcription string
}

type Service struct {
	mongo    *provider_db.MongoClient
	idUser   string
	interval int
}

func SendEmail(message, email string) error {
	const smtpServer = "smtp.yandex.ru"
	const mailSender = "cartmanis@yandex.ru"
	const password = "5eu7ve"
	//tlsconfig := &tls.Config {
	//	InsecureSkipVerify: true,
	//	ServerName: "smtp.gmail.com",
	//}
	//conn , err  := tls.Dial ("tcp",  smtpServer+":465", tlsconfig)
	//if err != nil {
	//	return err
	//}
	//defer conn.Close()
	//client, err := smtp.NewClient(conn, smtpServer)
	//if err != nil {
	//	return err
	//}
	//defer client.Close()
	//

	//if err := client.Mail(mailSender); err != nil {
	//	return err
	//}
	//if err := client.Rcpt(email); err != nil {
	//	return err
	//}
	//wc, err := client.Data()
	//if err != nil {
	//	return err
	//}
	//defer wc.Close()
	//buf := bytes.NewBufferString(msg)
	//if _, err = buf.WriteTo(wc); err != nil {
	//	return err
	//}
	//client.Quit()
	msg := "From: dictonaryyys\r\n " +
		"To: " + email + "\r\n" +
		"Subject: Активация аккаунта english_dictonary\r\n" +
		"\r\n" +
		"Регистрация почти закончена, осталось только подтвердить электронную почту, указанную при регистрации\r\n" +
		"Для подтверждения регистрации на dictonary_english необходимо перейти по ссылке:.\r\n" +
		"<button>Кнопка регистрации</button>" + "\r\n" +
		"http://127.0.0.1:27333/api/v1/activate?id=177891\r\n"
	auth := smtp.PlainAuth("", mailSender, password, smtpServer)

	if err := smtp.SendMail(smtpServer+":587", auth, mailSender,
		[]string{email}, []byte(msg)); err != nil {
		return err
	}
	return nil
}

func checkService(s *Service) error {
	if s == nil {
		return fmt.Errorf("не проиницилизированный Service")
	}
	if s.mongo == nil {
		return fmt.Errorf("не установлено подключение к базе данных")
	}
	return nil
}

func NewService(idUser string, interval int, m *provider_db.MongoClient) *Service {
	return &Service{m, idUser, interval}
}

func (s *Service) GetRandomWord() (*word, error) {
	if err := checkService(s); err != nil {
		return nil, err
	}
	listWord := make([]*word, 0)
	filter := struct {
		IdUser string `bson:"id_user"`
	}{s.idUser}
	err := s.mongo.Find(&filter, &listWord, "english")
	if err != nil {
		return nil, err
	}
	if len(listWord) == 0 {
		return nil, fmt.Errorf("в базе данных нет ни одного слова")
	}
	rand.Seed(time.Now().Unix())
	index := rand.Intn(len(listWord))
	return listWord[index], nil
}
