package db

import (
	"english_dictonary/app/provider_db"
	"fmt"
)

func InsertWord(en, ru, trans, idUser string, m *provider_db.MongoClient) (interface{}, error) {
	if m == nil {
		return nil, fmt.Errorf("не иницилизированная база данных mongoDb")
	}
	var word = map[string]string{
		"en":            en,
		"ru":            ru,
		"transcription": trans,
		"id_user":       idUser,
	}
	return m.InsertOne(word, users)
}
