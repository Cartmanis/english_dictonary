package rest

import (
	"english_dictonary/app/cmd/lg"
	"english_dictonary/app/provider_db"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	R "github.com/go-pkgz/rest"
	"net/http"
	"strconv"
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

	r.Post("/auth", s.autharization)
	r.Post("/login", s.login)
	r.Post("/logout", s.logout)
	r.Post("/add_user", s.newUser)
	r.Post("/add_word", s.newWord)

	return http.ListenAndServe(":"+strconv.Itoa(s.port), r)
}

func SendJSON(w http.ResponseWriter, r *http.Request, status int, i interface{}) {
	render.Status(r, status)
	render.JSON(w, r, i)
}

func SendErrorJSON(w http.ResponseWriter, r *http.Request, httpStatusCode int, details string, err error) {
	render.Status(r, httpStatusCode)
	render.JSON(w, r, R.JSON{"error": err.Error(), "details": details})
}
