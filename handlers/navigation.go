package handlers

import (
	"github.com/dmidokov/rv2/navigation"
	"github.com/dmidokov/rv2/response"
	"github.com/dmidokov/rv2/storage/postgres/user"
	"net/http"
)

func (hm *Service) GetNavigation(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		return
	}

	log := hm.Logger
	method := "api.navigation.get"
	response := resp.Service{Writer: &w, Logger: hm.Logger, Operation: method}

	navigationService := navigation.Service{DB: hm.DB, CookieStore: hm.CookieStore, Log: hm.Logger}
	userService := user.Service{DB: hm.DB, CookieStore: hm.CookieStore, Log: hm.Logger}

	userId := userService.GetUserIdFromSession(r)
	if userId == 0 {
		response.Unauthorized()
		return
	}

	navigationItems, err := navigationService.Get(userId)
	if err != nil {
		log.Errorf("Error: %s", err.Error())

		response.InternalServerError()
		return
	}

	response.OKWithData(navigationItems)

}
