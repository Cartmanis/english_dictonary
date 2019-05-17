package main

import (
	"english_dictonary/app/db"
	"english_dictonary/app/provider_db"
	"fmt"
	"os"
	"time"
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
	go RunRest(m)
	if err := db.MigratorIndex(m); err != nil {
		fmt.Println("[ERROR] не удалось проверить или добавить уникальный индекс для поля login. Ошибка:", err)
		return
	}
	for {
		time.Sleep(10 * time.Minute)
	}

	//words := make([]interface{}, 0)
	//words = append(words, &Word{En: "most", Ru:"наибольшее количество",Transcription: "мост"})
	//words = append(words, &Word{"land", "Земля", "Лаанд"})
	//word:= &Word{"to say", "рассказ", "ту сэй"}
}
