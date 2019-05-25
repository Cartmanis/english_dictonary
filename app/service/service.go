package service

import (
	"fmt"
	"github.com/cartmanis/english_dictonary/app/provider_db"
	"math/rand"
	"time"
)

type word struct {
	En            string
	Ru            string
	Transcription string
	Date          time.Time
}

type Service struct {
	mongo    *provider_db.MongoClient
	idUser   string
	interval int
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
	const location = "Local"
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
	loc, err := time.LoadLocation(location)
	if err != nil {
		return nil, err
	}
	rand.Seed(time.Now().Unix())
	index := rand.Intn(len(listWord))
	listWord[index].Date = listWord[index].Date.In(loc)
	return listWord[index], nil
}
