package rest

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"strconv"
	"time"
)

var mySiginKey = os.Getenv("SECRET_KEY")

func Run(port int) error {
	fmt.Println("[INFO] рест сервер запускается на порту:", port)
	r := mux.NewRouter()

	r.Handle("/", http.FileServer(http.Dir("./views")))
	r.Handle("/status", NotImplemented).Methods("GET")
	r.Handle("/products", NotImplemented).Methods("GET")
	r.Handle("/products/{slug}/feedback", NotImplemented).Methods("POST")

	//авторизация
	r.Handle("/get-token", GetTokenHandler).Methods("GET")

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static",
		http.FileServer(http.Dir("./static"))))

	return http.ListenAndServe(":"+strconv.Itoa(port), handlers.LoggingHandler(os.Stdout, r))
}

var NotImplemented = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Not Implemented"))
})

var GetTokenHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	token := jwt.New(jwt.SigningMethodHS256)
	token.Header["admin"] = true
	token.Header["name"] = "Logiiin"
	token.Header["exp"] = time.Now().Add(time.Hour * 24).Unix()
	tokenString, err := token.SignedString(mySiginKey)

	if err != nil {
		fmt.Println("[ERROR]: ", err)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte(tokenString))
})
