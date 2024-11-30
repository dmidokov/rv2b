package user

import (
	"github.com/dmidokov/rv2/lib/entitie"
	resp "github.com/dmidokov/rv2/response"
	"net/http"
)

type userSwitchGetter interface {
	GetUserIdFromSession(r *http.Request) int
	GetById(userId int) (*entitie.User, error)
	GetUsersToSwitch(userId int) ([]*entitie.UserSwitcher, error)
}

func (s *Service) GetSwitcher(
	userProvider userSwitchGetter,
	rightsProvider rightsSetter,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method == http.MethodOptions {
			return
		}

		log := s.Logger
		method := "api.user.UpdateUserRights"
		log.Info(method)

		response := resp.Service{Writer: &w, Logger: s.Logger, Operation: method}

		currentUserId := userProvider.GetUserIdFromSession(r)
		if currentUserId == 0 {
			log.Error("Unauthorized")
			response.Unauthorized()
			return
		}

		// TODO: тут еще надо проверить, что пользователи потомки одного родителя ???
		userSwitcher, err := userProvider.GetUsersToSwitch(currentUserId)
		if err != nil {
			log.Errorf("Error: %s", err.Error())
			response.InternalServerError()
			return
		}

		response.OKWithData(userSwitcher)

	}
}
