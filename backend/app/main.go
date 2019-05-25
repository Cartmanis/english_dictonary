package main

import (
	"github.com/cartmanis/english_dictonary/backend/app/cmd/lg"
	"github.com/cartmanis/english_dictonary/backend/app/db"
	"github.com/cartmanis/english_dictonary/backend/app/provider_db"
	"os"
	"time"
)

const (
	url     = "mongodb://localhost:27017"
	dbName  = "dictionary"
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
		lg.Error("не удалось подключиться к базе данных mongodb. Ошибка:", err)
		return
	}
	defer m.Close()
	if err := db.MigratorIndex(m); err != nil {
		lg.Error("не удалось проверить или добавить уникальный индекс для поля login. Ошибка:", err)
		return
	}

	go RunRest(m)
	for {
		time.Sleep(10 * time.Minute)
	}
}
