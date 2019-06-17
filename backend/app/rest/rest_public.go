package rest

import (
	"encoding/json"
	"fmt"
	"github.com/cartmanis/english_dictonary/backend/app/service"
	"net/http"
	"time"
)

const (
	english = "english"
)

func (s *Rest) activate(w http.ResponseWriter, r *http.Request) {
	if !checkInitRest(s, w, r) {
		return
	}
	values := r.URL.Query()
	userId := values.Get("id")
	if err := service.ActivateUser(userId, s.mongo); err != nil {
		SendErrorJSON(w, r, 500, "не удалось произвести активацию пользователя", err)
		return
	}

	token, err := getJwtToken(userId)
	if err != nil {
		SendErrorJSON(w, r, 500, "не удалось создать jwt токен. Ошибка:", err)
		return
	}
	c := &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		Secure:   false, //cookie отсылаются на сервер только если запрос выполняется по протоколу SSL и HTTPS
		HttpOnly: true,  //Куки HTTPonly не доступны из JavaScript через свойства Document.cookie API, что помогает избежать межсайтового скриптинга (XSS)
	}
	http.SetCookie(w, c)
	http.Redirect(w, r, "http://192.168.0.83:8080", 303)
}

func (s *Rest) run(w http.ResponseWriter, r *http.Request) {
	ok, _, id, _ := s.isAuthSession(w, r)
	if !ok {
		SendJSON(w, r, 401, map[string]bool{"result": false})
		return
	}
	sr := service.NewService(id, 1, s.mongo)
	word, err := sr.GetRandomWord()
	if err != nil {
		SendErrorJSON(w, r, 500, "не удалось найти слово", err)
		return
	}
	SendJSON(w, r, 200, map[string]interface{}{
		"En":            word.En,
		"Ru":            word.Ru,
		"Transcription": word.Transcription,
	})
}

func (s *Rest) newUser(w http.ResponseWriter, r *http.Request) {
	if !checkInitRest(s, w, r) {
		return
	}
	login := r.PostFormValue("login")
	password := r.PostFormValue("password")
	email := r.PostFormValue("email")
	phone := r.PostFormValue("phone")
	if login == "" || password == "" {
		SendErrorJSON(w, r, 400, "не заполнены имя пользователя или пароль", nil)
		return
	}
	if err := checkEmail(email); err != nil {
		SendErrorJSON(w, r, 400, "не корректно заполнен email", err)
		return
	}

	//objectId, err := service.InsertUser(login, password, email, phone, s.mongo)
	//if err != nil {
	//	SendErrorJSON(w, r, 200, "не удалось зарегистрировать пользователя", err)
	//	return
	//}
	//id, err := service.GetIdString(objectId)
	//if err != nil {
	//	SendErrorJSON(w, r, 500, "не удалось зарегистрировать пользователя", err)
	//	return
	//}

	//urlConfirm, err := getUrlConfirmEmail(id)
	//if err != nil {
	//	SendErrorJSON(w, r, 200, "не удалось зарегистрировать пользователя", err)
	//	return
	//}

	//if err := service.SendEmail(urlConfirm, email); err != nil {
	//	SendErrorJSON(w, r, 500, "не удалось отправить ссылку подтвержения на электронный адрес", err)
	//	return
	//}
	urlEmail := getUrlUserEmail(email)
	SendJSON(w, r, 200, map[string]interface{}{"result": true, "url": urlEmail})
}

func (s *Rest) newWord(w http.ResponseWriter, r *http.Request) {
	if !checkInitRest(s, w, r) {
		return
	}
	//стоит ли делать всегда авторизацию???
	ok, _, userId, _ := s.isAuthSession(w, r)
	if !ok {
		SendJSON(w, r, 401, map[string]bool{"result": false})
		return
	}
	ww := struct {
		En            string
		Ru            string
		Transcription string
		UserId        string `bson:"id_user"`
		Date          time.Time
	}{UserId: userId, Date: time.Now()}
	if err := json.NewDecoder(r.Body).Decode(&ww); err != nil {
		SendErrorJSON(w, r, 500, "не удалось добавить слово", err)
		return
	}
	_, err := s.mongo.InsertOne(ww, english)
	if err != nil {
		SendErrorJSON(w, r, 500, "не удалось добавить слово", err)
		return
	}
	SendJSON(w, r, 200, map[string]bool{"result": true})
}

func (s *Rest) deleteWord(w http.ResponseWriter, r *http.Request) {
	if !checkInitRest(s, w, r) {
		return
	}
	ok, _, id, _ := s.isAuthSession(w, r)
	if !ok {
		SendJSON(w, r, 401, map[string]bool{"result": false})
		return
	}
	idUser := r.PostFormValue("id_user")
	//убеждаемя, что это слово этого пользователя иначе возвращаем код 403
	if idUser != id {
		SendJSON(w, r, 403, map[string]bool{"result": false})
		return
	}
	idWord := r.PostFormValue("id_word")
	count, err := service.DeleteWord(idWord, s.mongo)
	if err != nil {
		SendErrorJSON(w, r, 500, "не удалось удалить слово", err)
		return
	}
	if count < 1 {
		SendErrorJSON(w, r, 400, "не удалось удалить слово",
			fmt.Errorf("не найденно не одного подходящего слова для удаления"))
		return
	}
	SendJSON(w, r, 200, map[string]int64{"result": count})
}

func (s *Rest) forgetWord(w http.ResponseWriter, r *http.Request) {
	if !checkInitRest(s, w, r) {
		return
	}
	ok, _, id, _ := s.isAuthSession(w, r)
	if !ok {
		SendJSON(w, r, 401, map[string]bool{"result": false})
		return
	}
	idUser := r.PostFormValue("id_user")
	//убеждаемя, что это слово этого пользователя иначе возвращаем код 403
	if idUser != id {
		SendJSON(w, r, 403, map[string]bool{"result": false})
		return
	}
	idWord := r.PostFormValue("id_word")
	count, err := service.ForgetOrRecall(true, idWord, s.mongo)
	if err != nil {
		SendErrorJSON(w, r, 500, "не удалось забыть слово", err)
		return
	}
	if count < 1 {
		SendErrorJSON(w, r, 400, "не удалось забвть слово",
			fmt.Errorf("не найденно не одного подходящего слова"))
		return
	}
	SendJSON(w, r, 200, map[string]int64{"result": count})
}
