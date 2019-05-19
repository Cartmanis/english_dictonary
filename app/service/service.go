package service

import (
	"english_dictonary/app/provider_db"
	"fmt"
)

type word struct {
	En            string
	Ru            string
	Transcription string
}

type Service struct {
	mongo    *provider_db.MongoClient
	idUser   interface{}
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

func NewService(idUser string, interval int, m *provider_db.MongoClient) (*Service, error) {
	idObject, err := provider_db.GetObjectId(idUser)
	if err != nil {
		return nil, err
	}
	return &Service{m, idObject, interval}, nil
}

func (s *Service) GetRandomWord() (*word, error) {
	if err := checkService(s); err != nil {
		return nil, err
	}
	listWord := map[string]string{}
	err := s.mongo.Find(map[string]interface{}{"id_user": s.idUser}, listWord, "english")
	if err != nil {
		return nil, err
	}
	return nil, nil
}
