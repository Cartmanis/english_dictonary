package rest

import (
	"fmt"
	"github.com/cartmanis/english_dictonary/app/cmd/lg"
	"github.com/cartmanis/english_dictonary/app/db"
	"github.com/cartmanis/english_dictonary/app/provider_db"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	R "github.com/go-pkgz/rest"
	"github.com/gorilla/sessions"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Rest struct {
	port  int
	mongo *provider_db.MongoClient
}

func NewRestService(port int, mongo *provider_db.MongoClient) *Rest {
	return &Rest{
		port:  port,
		mongo: mongo,
	}
}

var store = sessions.NewCookieStore([]byte(os.Getenv("SECRET_KEY")))

func (s *Rest) Run() error {
	lg.Info("рест сервер запускается на порту:", s.port)
	r := chi.NewRouter()

	//r.Handle("/", http.FileServer(http.Dir("./views")))
	r.Post("/auth", s.autharization)
	r.Post("/login", s.login)
	r.Post("/logout", s.logout)

	//авторизация
	//r.Get("/get-token", GetTokenHandler)

	//r.PathPrefix("/static/").Handler(http.StripPrefix("/static",
	//	http.FileServer(http.Dir("./static"))))

	return http.ListenAndServe(":"+strconv.Itoa(s.port), r)
}

var NotImplemented = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Not Implemented"))
})

var GetTokenHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	token := jwt.New(jwt.SigningMethodHS256)
	token.Header["admin"] = true
	token.Header["name"] = "Logiiin"
	token.Header["exp"] = time.Now().Add(time.Hour * 24).Unix()
	tokenString, err := token.SignedString(os.Getenv("SECRET_KEY"))

	if err != nil {
		fmt.Println("[ERROR]: ", err)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte(tokenString))
})

func (s *Rest) autharization(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "user_session")
	if err != nil {
		SendErrorJSON(w, r, 500, "не удалось получить сессию", err)
		return
	}
	id := session.Values["user_id"]
	if id == nil {
		SendErrorJSON(w, r, 403, "нет авторизации",
			fmt.Errorf("пустая сессия"))
		return
	}
	user, err := db.FindUserByIdUser(id.(string), s.mongo)
	if err != nil {
		SendErrorJSON(w, r, 403, "авторизация", err)
		return
	}
	fmt.Println(user)
}

func (s *Rest) login(w http.ResponseWriter, r *http.Request) {
	login, password, _ := r.BasicAuth()
	auth, id, err := db.AuthUser(login, password, s.mongo)
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
	//store.Options.HttpOnly = true
	//store.Options.Secure = true
	//store.Options.MaxAge = 0
	fmt.Println(id)
	session.Values["user_id"] = interface{}(true)
	if err := session.Save(r, w); err != nil {
		SendErrorJSON(w, r, 500, "не удалось сохранить сессию", err)
		return
	}
	//http.Redirect(w, r, host+":30500/api/v1/sc_started", 303)
	SendJSON(w, r, 200, R.JSON{"result": true})
}

func (s *Rest) logout(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "user_session")
	if err != nil {
		SendErrorJSON(w, r, 500, "не удалось удалить сессию", err)
		return
	}
	session.Values["user_id"] = ""
	if err := session.Save(r, w); err != nil {
		SendErrorJSON(w, r, 500, "не удалось удалить сессию", err)
		return
	}
	//http.Redirect(w,r, "localhost:" + strconv.Itoa(s.port) + "/login", 303)
}

func SendJSON(w http.ResponseWriter, r *http.Request, status int, i interface{}) {
	render.Status(r, status)
	render.JSON(w, r, i)
}

func SendErrorJSON(w http.ResponseWriter, r *http.Request, httpStatusCode int, details string, err error) {
	render.Status(r, httpStatusCode)
	render.JSON(w, r, R.JSON{"error": err.Error(), "details": details})
}
