package rest

import (
	"encoding/json"
	"english_dictonary/app/db"
	"fmt"
	"net/http"
)

const (
	english = "english"
)

func (s *Rest) newUser(w http.ResponseWriter, r *http.Request) {
	if s == nil || s.mongo == nil {
		SendErrorJSON(w, r, 500, "не удалось зарегистрировать пользователя",
			fmt.Errorf("не инициализированный рест сервер"))
		return
	}
	login := r.PostFormValue("login")
	password := r.PostFormValue("password")
	email := r.PostFormValue("email")
	phone := r.PostFormValue("phone")
	id, err := db.InsertUser(login, password, email, phone, s.mongo)
	if err != nil {
		SendErrorJSON(w, r, 400, "не удалось зарегистрировать пользователя", err)
		return
	}
	fmt.Println(id)
	SendJSON(w, r, 200, map[string]bool{"result": true})
}

func (s *Rest) newWord(w http.ResponseWriter, r *http.Request) {
	if s == nil || s.mongo == nil {
		SendErrorJSON(w, r, 500, "не удалось зарегистрировать пользователя",
			fmt.Errorf("не инициализированный рест сервер"))
		return
	}
	//стоит ли делать всегда авторизацию???
	ok, userId := s.isAuthSession(w, r)
	if !ok {
		SendJSON(w, r, 403, map[string]bool{"result": false})
		return
	}
	words := map[string]string{"userId": userId}
	if err := json.NewDecoder(r.Body).Decode(&words); err != nil {
		SendErrorJSON(w, r, 500, "не удалось добавить слово", err)
		return
	}
	_, err := s.mongo.InsertOne(words, english)
	if err != nil {
		SendErrorJSON(w, r, 500, "не удалось добавить слово", err)
		return
	}
	SendJSON(w, r, 200, map[string]bool{"result": true})
}
