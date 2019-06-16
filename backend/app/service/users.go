package service

import (
	"fmt"
	"github.com/cartmanis/english_dictonary/backend/app/cmd/crypto"
	"github.com/cartmanis/english_dictonary/backend/app/provider_db"
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
	Activate bool
}

type NewUser struct {
	Login    string
	Password string
	Email    string
	Phone    string `bson:",omitempty"`
	Activate bool
}

func InsertUser(login, pass, email, phone string, m *provider_db.MongoClient) (interface{}, error) {
	if m == nil {
		return nil, fmt.Errorf("не иницилизированная база данных mongoDb")
	}
	type filter struct {
		Login string
	}
	countLogin, err := m.Count(&filter{login}, users)
	if err != nil {
		return nil, err
	}
	if countLogin >= 1 {
		return nil, fmt.Errorf("Пользователь с именем %v уже существует. Используйте другое имя для регистрации.", login)
	}
	if email != "" {
		countEmail, err := m.Count(map[string]string{"email": email}, users)
		if err != nil {
			return nil, err
		}
		if countEmail >= 1 {
			return nil, fmt.Errorf("Пользователь с email %v уже зарегистрирован в системе.", email)
		}
	}

	passHash, err := crypto.GetHashPassword(pass)
	if err != nil {
		return nil, err
	}
	newUser := &NewUser{login, passHash, email, phone, false}
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

func AuthUser(login, pass string, m *provider_db.MongoClient) (bool, int, string, error) {
	if m == nil {
		return false, 500, "", fmt.Errorf("не иницилизированная база данных mongoDb")
	}
	type filter struct {
		Login string
	}
	listUser := make([]*User, 0)
	if err := m.Find(&filter{login}, &listUser, users); err != nil {
		return false, 500, "", err
	}
	if len(listUser) == 0 {
		return false, 401, "", nil
	}
	if len(listUser) > 1 {
		return false, 500, "", fmt.Errorf("не возможно однозначно идентифицировать. "+
			"Найдено более одного пользователя с login: %v", login)
	}
	if !crypto.CompareHashPassword(listUser[0].Password, pass) {
		return false, 401, "", nil
	}
	if !listUser[0].Activate {
		return false, 403, "", nil
	}
	return true, 200, listUser[0].Id.(primitive.ObjectID).Hex(), nil
}
