package main

import (
	"english_dictonary/app/provider_db"
	"english_dictonary/app/rest"
	"fmt"
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
