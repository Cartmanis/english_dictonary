package rest

import (
	"english_dictonary/app/db"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	R "github.com/go-pkgz/rest"
	"github.com/gorilla/sessions"
	"net/http"
	"os"
	"time"
)

var (
	secretKey = []byte(os.Getenv("SECRET_KEY"))
	store     = sessions.NewCookieStore(secretKey)
)

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
	token := session.Values["user_id"]
	if token == nil {
		return false, ""
	}
	id, err := verifyJwtToken(token.(string))
	if err != nil {
		return false, ""
	}
	_, err = db.FindUserByIdUser(id, s.mongo)
	if err != nil {
		return false, ""
	}
	return true, id
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
	token, err := getJwtToken(userId)
	if err != nil {
		SendErrorJSON(w, r, 500, "не удалось создать jwt токен. Ошибка:", err)
		return
	}
	session.Values["user_id"] = token //jwt токен
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

func getJwtToken(id string) (string, error) {
	var token = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": id,
		"exp":     time.Now().Add(time.Hour * 1).Unix(),
	})
	return token.SignedString(secretKey)
}

func verifyJwtToken(token string) (string, error) {
	tok, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неожиданный метод подписания %v", token.Header["alg"])
		}
		return secretKey, nil
	})
	if err != nil {
		return "", err
	}
	id := tok.Claims.(jwt.MapClaims)["user_id"].(string)
	return id, nil
}
