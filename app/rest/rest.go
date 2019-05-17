package rest

import (
	"english_dictonary/app/cmd/lg"
	"english_dictonary/app/provider_db"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	R "github.com/go-pkgz/rest"
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

func SendJSON(w http.ResponseWriter, r *http.Request, status int, i interface{}) {
	render.Status(r, status)
	render.JSON(w, r, i)
}

func SendErrorJSON(w http.ResponseWriter, r *http.Request, httpStatusCode int, details string, err error) {
	render.Status(r, httpStatusCode)
	render.JSON(w, r, R.JSON{"error": err.Error(), "details": details})
}
