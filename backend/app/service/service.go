package service

import (
	"fmt"
	"github.com/cartmanis/english_dictonary/backend/app/provider_db"
	"math/rand"
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
