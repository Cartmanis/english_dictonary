package main

import (
	"english_dictonary/app/cmd/crypto"
	"english_dictonary/app/db"
	"english_dictonary/app/provider_db"
	"fmt"
	"os"
)

const (
	url     = "mongodb://localhost:27017"
	dbName  = "books"
	english = "english"
)

type Word struct {
	En            string
	Ru            string
	Transcription string
	UserId        interface{}
}

type Filter struct {
	Id interface{} `bson:"_id"`
}

type FilterLogin struct {
	Login string
}

type Update struct {
	//Ru string
	User string
}

var (
	login    = os.Getenv("MONGO_USERNAME")
	password = os.Getenv("MONGO_PASSWORD")
)

func main() {
	m, err := provider_db.NewStoreContext(url, dbName, login, password, 15)
	if err != nil {
		fmt.Println("[ERROR] не удалось подключиться к базе данных mongodb. Ошибка:", err)
		return
	}
	defer m.Close()

	token, err := crypto.CreateToken()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(token)
	RunRest(m)
	if err := db.MigratorIndex(m); err != nil {
		fmt.Println("[ERROR] не удалось проверить или добавить уникальный индекс для поля login. Ошибка:", err)
		return
	}
	//id, err := db.InsertUser("admin", "5eu7ve", m)
	//if err != nil {
	//	fmt.Println("[ERROR] не удалось добавить нового пользователя. Ошибка:",err)
	//}
	//fmt.Println(id)
	isAuth, id, err := db.AuthUser("mila", "5eu7ve", m)
	if err != nil {
		fmt.Println("[ERROR] не удалось произвести авторизацию. Ошибка:", err)
		return
	}
	if isAuth {
		fmt.Println(id)
	}
	var s string
	fmt.Scanln(&s)

	//words := make([]interface{}, 0)
	//words = append(words, &Word{En: "most", Ru:"наибольшее количество",Transcription: "мост"})
	//words = append(words, &Word{"land", "Земля", "Лаанд"})
	//word:= &Word{"to say", "рассказ", "ту сэй"}
}
