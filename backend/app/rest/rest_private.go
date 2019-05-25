package rest

import (
	"fmt"
	"net/http"
)

func checkInitRest(s *Rest, w http.ResponseWriter, r *http.Request) bool {
	if s == nil || s.mongo == nil {
		SendErrorJSON(w, r, 500, "не удалось зарегистрировать пользователя",
			fmt.Errorf("не инициализированный рест сервер"))
		return false
	}
	return true
}
