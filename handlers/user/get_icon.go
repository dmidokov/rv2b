package user

import (
	e "github.com/dmidokov/rv2/entitie"
	resp "github.com/dmidokov/rv2/response"
	"net/http"
)

type userIconProvider interface {
	GetIcon(int) *e.UserIcon
	GetUserIdFromSession(r *http.Request) int
}

func (s *Service) GetUserIcon(userIconProvider userIconProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			return
		}

		method := "api.users.getIcon"
		response := resp.Service{Writer: &w, Logger: s.Logger, Operation: method}

		userId := userIconProvider.GetUserIdFromSession(r)
		if userId == 0 {
			response.Unauthorized()
			return
		}

		userIcon := userIconProvider.GetIcon(userId)

		response.OKWithData(userIcon)
	}
}
