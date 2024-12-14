package user

import (
	"fmt"
	"github.com/dmidokov/rv2/lib"
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
	GetParentId(userId int) (int, error)
	GetChild(userId int) ([]*e.UserShort, error)
}

func (s *Service) GetUsers(userProvider userProvider, rightsProvider rightsSetter) http.HandlerFunc {
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
		log.Info(method)
		currentUserId := userProvider.GetUserIdFromSession(r)
		if currentUserId == 0 {
			response.Unauthorized()
			return
		}
		log.Info(method)
		currentUser, err := userProvider.GetById(currentUserId)
		if err != nil {
			response.InternalServerError()
			return
		}
		log.Info(method)

		fmt.Println(currentUser.Rights)
		fmt.Println(lib.ViewUsers)
		fmt.Println(rightsProvider.CheckUserRight(currentUser, lib.ViewUsers))

		if rightsProvider.CheckUserRight(currentUser, lib.ViewUsers) {
			items, err := userProvider.GetByOrganizationId(currentUserId)
			if err != nil {
				log.Errorf("Error: %s", err.Error())

				response.InternalServerError()
				return
			}

			response.OKWithData(items)
			return
		}

		response.NotAllowed()
	}
}
