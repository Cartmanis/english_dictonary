package rest

import (
	"bytes"
	"fmt"
	"github.com/cartmanis/english_dictonary/backend/app/service"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/sessions"
	"net/http"
	"net/smtp"
	"os"
	"time"
)

var (
	secretKey = []byte(os.Getenv("SECRET_KEY"))
	store     = sessions.NewCookieStore(secretKey)
)

func (s *Rest) confirmEmail(w http.ResponseWriter, r *http.Request) {
	const email = "VShmelcer@gmail.com"
	const link = ""
	const details = "не удалось подтвердить электронную почту"
	// Connect to the remote SMTP server.
	c, err := smtp.Dial("127.0.0.1:2525")
	if err != nil {
		SendErrorJSON(w, r, 500, details, err)
		return
	}
	// Set the sender and recipient.
	if err := c.Mail("sender@example.org"); err != nil {
		SendErrorJSON(w, r, 500, details, err)
		return
	}
	if err := c.Rcpt(email); err != nil {
		SendErrorJSON(w, r, 500, details, err)
		return
	}
	// Send the email body.
	wc, err := c.Data()
	if err != nil {
		SendErrorJSON(w, r, 500, details, err)
		return
	}
	defer wc.Close()
	buf := bytes.NewBufferString("Please confirm this email " + link)
	if _, err = buf.WriteTo(wc); err != nil {
		SendErrorJSON(w, r, 500, details, err)
	}
	SendJSON(w, r, 200, map[string]bool{"result": true})
}

func (s *Rest) autharization(w http.ResponseWriter, r *http.Request) {
	ok, status, id, userName := s.isAuthSession(w, r)
	if !ok {
		SendJSON(w, r, status, map[string]bool{"result": false})
		return
	}
	if status == 403 {
		SendJSON(w, r, status, map[string]interface{}{"result": false, "user": userName})
		return
	}
	SendJSON(w, r, 200, map[string]interface{}{"result": true, "user": userName})
	fmt.Println(id) //используем полученного user далее в коде
}

func (s *Rest) isAuthSession(w http.ResponseWriter, r *http.Request) (bool, int, string, string) {
	c, err := r.Cookie("token")
	if err != nil {
		return false, 401, "", ""
	}
	id, err := verifyJwtToken(c.Value)
	if err != nil {
		return false, 401, "", ""
	}
	user, err := service.FindUserByIdUser(id, s.mongo)
	if err != nil {
		return false, 401, "", ""
	}
	if !user.Activate {
		return true, 403, "", user.Login
	}
	return true, 200, id, user.Login
}

func (s *Rest) login(w http.ResponseWriter, r *http.Request) {
	login, password, _ := r.BasicAuth()
	auth, status, userId, err := service.AuthUser(login, password, s.mongo)
	if err != nil {
		SendErrorJSON(w, r, 500, "не удалось произвести авторизацию", err)
		return
	}
	if !auth {
		SendJSON(w, r, status, map[string]bool{"result": false})
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

	SendJSON(w, r, 200, map[string]bool{"result": true})
}

func (s *Rest) logout(w http.ResponseWriter, r *http.Request) {
	//предыдущая реализация через session
	//session, err := store.Get(r, "user_session")
	//if err != nil {
	//	SendErrorJSON(w, r, 500, "не удалось удалить сессию", err)
	//	return
	//}
	////сбрасываем cookie
	//session.Options.MaxAge = -1
	//if err := session.Save(r, w); err != nil {
	//	SendErrorJSON(w, r, 500, "не удалось удалить сессию", err)
	//	return
	//}
	c := &http.Cookie{
		Name:    "token",
		Value:   "",
		Path:    "/",
		Expires: time.Unix(0, 0),

		HttpOnly: true,
	}
	http.SetCookie(w, c)
	SendJSON(w, r, 200, map[string]bool{"result": true})
	//http.Redirect(w,r, "localhost:" + strconv.Itoa(s.port) + "/login", 303)
}

func getJwtToken(id string) (string, error) {
	var token = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": id,
		"exp":     time.Now().Add(time.Hour * 3).Unix(),
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
