package db

import (
	"fmt"
	"github.com/cartmanis/english_dictonary/app/cmd/crypto"
	"github.com/cartmanis/english_dictonary/app/provider_db"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	Email    string
	Phone    string `bson:",omitempty"`
}

func InsertUser(login, pass, email, phone string, m *provider_db.MongoClient) (interface{}, error) {
	if m == nil {
		return nil, fmt.Errorf("не иницилизированная база данных mongoDb")
	}
	type filter struct {
		Login string
	}
	count, err := m.Count(&filter{login}, users)
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
	newUser := &NewUser{login, passHash, email, phone}
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

func FindUserByIdUser(userId string, m *provider_db.MongoClient) (*User, error) {
	if m == nil {
		return nil, fmt.Errorf("не иницилизированная база данных mongoDb")
	}
	type filter struct {
		Id interface{} `bson:"_id"`
	}
	objectId, err := provider_db.GetObjectId(userId)
	if err != nil {
		return nil, fmt.Errorf("не возможно привести к типу objectId, так как не корректный userId: %v. "+
			"Дополнительно:%v", userId, err)
	}
	user := &User{}
	if err := m.FindOne(&filter{objectId}, user, users); err != nil {
		return nil, err
	}
	return user, nil
}

func AuthUser(login, pass string, m *provider_db.MongoClient) (bool, string, error) {
	if m == nil {
		return false, "", fmt.Errorf("не иницилизированная база данных mongoDb")
	}
	type filter struct {
		Login string
	}
	listUser := make([]*User, 0)
	if err := m.Find(&filter{login}, &listUser, users); err != nil {
		return false, "", err
	}
	if len(listUser) == 0 {
		return false, "", nil
	}
	if len(listUser) > 1 {
		return false, "", fmt.Errorf("не возможно однозначно идентифицировать. "+
			"Найдено более одного пользователя с login: %v", login)
	}
	if !crypto.CompareHashPassword(listUser[0].Password, pass) {
		return false, "", nil
	}
	return true, listUser[0].Id.(primitive.ObjectID).Hex(), nil
}
