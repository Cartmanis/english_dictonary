package rest

import (
	"english_dictonary/app/db"
	"fmt"
	R "github.com/go-pkgz/rest"
	"github.com/gorilla/sessions"
	"net/http"
	"os"
)

var store = sessions.NewCookieStore([]byte(os.Getenv("SECRET_KEY")))

func (s *Rest) autharization(w http.ResponseWriter, r *http.Request) {
	ok, id := s.isAuthSession(w, r)
	if !ok {
		SendJSON(w, r, 403, R.JSON{"result": false})
		return
	}
	SendJSON(w, r, 200, R.JSON{"result": true})
	fmt.Println(id) //используем полученного user далее в коде
}

func (s *Rest) isAuthSession(w http.ResponseWriter, r *http.Request) (bool, string) {
	session, err := store.Get(r, "user_session")
	if err != nil {
		return false, ""
	}
	id := session.Values["user_id"]
	if id == nil {
		return false, ""
	}
	_, err = db.FindUserByIdUser(id.(string), s.mongo)
	if err != nil {
		return false, ""
	}
	return true, id.(string)
}

func (s *Rest) login(w http.ResponseWriter, r *http.Request) {
	login, password, _ := r.BasicAuth()
	auth, userId, err := db.AuthUser(login, password, s.mongo)
	if err != nil {
		SendErrorJSON(w, r, 500, "не удалось произвести авторизацию", err)
		return
	}
	if !auth {
		SendJSON(w, r, 403, R.JSON{"result": false})
		return
	}
	session, err := store.Get(r, "user_session")
	if err != nil {
		SendErrorJSON(w, r, 500, "не удалось получить сессию", err)
		return
	}
	store.Options.HttpOnly = true
	store.Options.MaxAge = 0
	session.Values["user_id"] = userId
	if err := session.Save(r, w); err != nil {
		SendErrorJSON(w, r, 500, "не удалось сохранить сессию", err)
		return
	}
	SendJSON(w, r, 200, R.JSON{"result": true})
	//fmt.Println(id)//используем полученый id далее в коде
}

func (s *Rest) logout(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "user_session")
	if err != nil {
		SendErrorJSON(w, r, 500, "не удалось удалить сессию", err)
		return
	}
	//сбрасываем cookie
	session.Options.MaxAge = -1
	if err := session.Save(r, w); err != nil {
		SendErrorJSON(w, r, 500, "не удалось удалить сессию", err)
		return
	}
	//http.Redirect(w,r, "localhost:" + strconv.Itoa(s.port) + "/login", 303)
}
