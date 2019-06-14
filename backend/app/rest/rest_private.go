package rest

import (
	"fmt"
	"net/http"
	"strings"
)

func checkInitRest(s *Rest, w http.ResponseWriter, r *http.Request) bool {
	if s == nil || s.mongo == nil {
		SendErrorJSON(w, r, 500, "не удалось зарегистрировать пользователя",
			fmt.Errorf("не инициализированный рест сервер"))
		return false
	}
	return true
}

func checkEmail(email string) error {
	if strings.Trim(email, " ") == "" {
		return fmt.Errorf("электронная почта не может быть пустой")
	}
	arr := strings.Split(email, "@")
	if len(arr) < 2 {
		return fmt.Errorf(`в электронной почте обязательно должен присутствовать символ '@'`)
	}
	if len(arr) > 2 {
		return fmt.Errorf(`в электронной почте может быть только один символ '@'`)
	}
	if !strings.Contains(arr[1], ".") {
		return fmt.Errorf("в доменне почты обязательно должна присутствовать точка")
	}
	return nil
}

func getUrlUserEmail(email string) string {
	mapEmails := map[string]string{
		"mail.ru":        "https://e.mail.ru/",
		"bk.ru":          "https://e.mail.ru/",
		"list.ru":        "https://e.mail.ru/",
		"inbox.ru":       "https://e.mail.ru/",
		"yandex.ru":      "https://mail.yandex.ru/",
		"ya.ru":          "https://mail.yandex.ru/",
		"yandex.ua":      "https://mail.yandex.ua/",
		"yandex.by":      "https://mail.yandex.by/",
		"yandex.kz":      "https://mail.yandex.kz/",
		"yandex.com":     "https://mail.yandex.com/",
		"gmail.com":      "https://mail.google.com/",
		"googlemail.com": "https://mail.google.com/",
		"outlook.com":    "https://mail.live.com/",
		"hotmail.com":    "https://mail.live.com/",
		"live.ru":        "https://mail.live.com/",
		"live.com":       "https://mail.live.com/",
		"me.com":         "https://www.icloud.com/",
		"icloud.com":     "https://www.icloud.com/",
		"rambler.ru":     "https://mail.rambler.ru/",
		"yahoo.com":      "https://mail.yahoo.com/",
		"ukr.net":        "https://mail.ukr.net/",
		"i.ua":           "http://mail.i.ua/",
		"bigmir.net":     "http://mail.bigmir.net/",
		"tut.by":         "https://mail.tut.by/",
		"inbox.lv":       "https://www.inbox.lv/",
		"mail.kz":        "http://mail.kz/",
	}
	arr := strings.Split(email, "@")
	if len(arr) < 2 {
		return ""
	}
	return mapEmails[arr[1]]
}
