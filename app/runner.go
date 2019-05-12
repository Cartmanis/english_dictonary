package main

import (
	"english_dictonary/app/rest"
	"fmt"
)

const (
	portRest = 27333
)

func RunRest() {
	if err := rest.Run(portRest); err != nil {
		fmt.Println("[ERROR] не удалось запустить Rest server на порту %v. Ошибка: %v", portRest, err)
		return
	}
}
