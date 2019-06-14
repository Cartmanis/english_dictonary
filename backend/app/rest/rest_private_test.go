package rest

import (
	"testing"
)

func TestGetUrlEmailUser(t *testing.T) {
	empty := getUrlUserEmail("")
	if empty != "" {
		t.Error("тест не пройден")
	}
	no := getUrlUserEmail("23214123dd")
	if no != "" {
		t.Error("тест не пройден")
	}
	yes := getUrlUserEmail("12ed2@kcmk")
	if yes != "" {
		t.Error("тест не пройден")
	}
	success := getUrlUserEmail("213e3d@mail.ru")
	if success != "https://e.mail.ru/" {
		t.Error("тест не пройден")
	}
	successCom := getUrlUserEmail("23ed2@gmail.com@12@1")
	if successCom != "https://mail.google.com/" {
		t.Error("тест не пройден")
	}
}
