package rest

import (
	"fmt"
	"github.com/cartmanis/english_dictonary/backend/app/cmd/lg"
	"github.com/cartmanis/english_dictonary/backend/app/provider_db"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"net/http"
	"sync"
	"time"
)

type Rest struct {
	port  int
	mongo *provider_db.MongoClient
	lock  sync.Mutex
}

func NewRestService(port int, mongo *provider_db.MongoClient) *Rest {
	return &Rest{
		port:  port,
		mongo: mongo,
	}
}

func (s *Rest) Run() error {
	lg.Info("рест сервер запускается на порту:", s.port)
	router := s.routes()
	s.lock.Lock()
	httpServer := &http.Server{
		Addr:              fmt.Sprintf(":%d", s.port),
		Handler:           router,
		ReadHeaderTimeout: 30 * time.Second,
		WriteTimeout:      0 * time.Second,
		IdleTimeout:       0 * time.Second,
		MaxHeaderBytes:    100000000,
	}
	s.lock.Unlock()

	return httpServer.ListenAndServe()
}

func (s *Rest) routes() chi.Router {
	r := chi.NewRouter()

	var corsMiddleware = cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", " X-Requested-With", "Content-Type", "X-XSRF-Token", "X-JWT"},
		ExposedHeaders:   []string{"Authorization"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	r.Use(corsMiddleware.Handler)

	//r.Use(middleware.RealIP)
	//r.Use(middleware.Throttle(10000), middleware.Timeout(10*time.Minute))

	r.Route("/api/v1", s.groupRouteApiV1)
	return r
}

func (s *Rest) groupRouteApiV1(groupRoute chi.Router) {
	groupRoute.Group(func(r chi.Router) {
		r.Get("/activate", s.activate)
		r.Get("/confirm_new_password", s.confimNewPassword)
		r.Get("/run", s.run)
		r.Post("/auth", s.autharization)
		r.Post("/login", s.login)
		r.Post("/logout", s.logout)
		r.Post("/confirm_phone", s.confirmPhone)
		r.Post("/add_user", s.newUser)
		r.Post("/recovery_password", s.recoveryPassword)
		r.Post("/add_word", s.newWord)
		r.Delete("/delete_word", s.deleteWord)
		r.Put("/forget_word", s.forgetWord)
	})
}

func SendJSON(w http.ResponseWriter, r *http.Request, status int, i interface{}) {
	render.Status(r, status)
	render.JSON(w, r, i)
}

func SendErrorJSON(w http.ResponseWriter, r *http.Request, httpStatusCode int, details string, err error) {
	render.Status(r, httpStatusCode)
	lg.Errorf("%v. Ошибка:%v", details, err)
	if err == nil {
		render.JSON(w, r, map[string]string{"error": "", "details": details})
		return
	}
	render.JSON(w, r, map[string]string{"error": err.Error(), "details": details})
}
