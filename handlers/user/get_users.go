package user

import (
	e "github.com/dmidokov/rv2/lib/entitie"
	resp "github.com/dmidokov/rv2/response"
	"net/http"
)

type userProvider interface {
	GetById(userId int) (*e.User, error)
	GetOrganizationIdFromSession(r *http.Request) int
	GetUserIdFromSession(r *http.Request) int
	GetByOrganizationId(userId int) ([]*e.UserShort, error)
	GetInfo(userId int, infoLevel int) (*e.UserInfoFull, error)
}

func (s *Service) GetUsers(userProvider userProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			return
		}

		log := s.Logger
		method := "api.users.get"
		response := resp.Service{Writer: &w, Logger: s.Logger, Operation: method}
		log.Info(method)

		organizationId := userProvider.GetOrganizationIdFromSession(r)
		if organizationId == 0 {
			response.Unauthorized()
			return
		}

		currentUserId := userProvider.GetUserIdFromSession(r)
		if currentUserId == 0 {
			response.Unauthorized()
			return
		}

		items, err := userProvider.GetByOrganizationId(currentUserId)
		if err != nil {
			log.Errorf("Error: %s", err.Error())

			response.InternalServerError()
			return
		}

		response.OKWithData(items)
	}
}
