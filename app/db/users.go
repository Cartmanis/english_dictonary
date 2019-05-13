package db

import (
	"english_dictonary/app/cmd/crypto"
	"english_dictonary/app/provider_db"
	"fmt"
)

const (
	users           = "users"
	noDocumentMongo = "mongo: no documents in result"
)

type User struct {
	Id       interface{} `bson:"_id"`
	Login    string
	Password string
}

type NewUser struct {
	Login    string
	Password string
}

type Filter struct {
	Login string
}

func InsertUser(login, pass string, m *provider_db.MongoClient) (interface{}, error) {
	if m == nil {
		return nil, fmt.Errorf("не иницилизированная база данных mongoDb")
	}
	filter := &Filter{login}
	count, err := m.Count(filter, users)
	if err != nil {
		return nil, err
	}
	if count >= 1 {
		return nil, fmt.Errorf("пользователь с именем %v уже существует. Используйте другое имя для регистрации.", login)
	}
	passHash, err := crypto.GetHashPassword(pass)
	if err != nil {
		return nil, err
	}
	newUser := &NewUser{login, passHash}
	return m.InsertOne(newUser, users)
}

func MigratorIndex(m *provider_db.MongoClient) error {
	if m == nil {
		return fmt.Errorf("не иницилизированная база данных mongoDb")
	}
	listIdx, err := m.ListIndexes(users)
	if err != nil {
		return err
	}
	for _, v := range listIdx {
		if v == "idx_login" {
			return nil
		}
	}
	_, err = m.CreateIndex("login", users, true)
	if err != nil {
		return err
	}
	return nil
}

func AuthUser(login, pass string, m *provider_db.MongoClient) (bool, interface{}, error) {
	if m == nil {
		return false, nil, fmt.Errorf("не иницилизированная база данных mongoDb")
	}
	filter := &Filter{login}
	findUser := &User{}
	if err := m.FindOne(filter, findUser, users); err != nil {
		if err.Error() == noDocumentMongo {
			return false, nil, nil
		}
		return false, nil, err
	}
	if !crypto.CompareHashPassword(findUser.Password, pass) {
		return false, nil, nil
	}
	return true, findUser.Id, nil
}
