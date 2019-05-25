package main

import (
	"fmt"
	"github.com/cartmanis/english_dictonary/backend/app/provider_db"
	"github.com/cartmanis/english_dictonary/backend/app/rest"
)

const (
	portRest = 27333
)

func RunRest(m *provider_db.MongoClient) {
	r := rest.NewRestService(27333, m)
	if err := r.Run(); err != nil {
		fmt.Println("[ERROR] не удалось запустить Rest server на порту %v. Ошибка: %v", portRest, err)
		return
	}
}
