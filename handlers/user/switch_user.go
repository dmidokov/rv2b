package user

import (
	"github.com/dmidokov/rv2/lib/entitie"
	resp "github.com/dmidokov/rv2/response"
	"net/http"
	"strconv"
)

type userSwitcher interface {
	GetUserIdFromSession(r *http.Request) int
	GetById(userId int) (*entitie.User, error)
	GetUsersToSwitch(userId int) ([]*entitie.UserSwitcher, error)
	CanUserSwitchToId(from int, to int) bool
}

func (s *Service) SwitchUser(
	userProvider userSwitcher,
	rightsProvider rightsSetter,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method == http.MethodOptions {
			return
		}

		log := s.Logger
		method := "api.user.switcher.switch"

		response := resp.Service{Writer: &w, Logger: s.Logger, Operation: method}

		query := r.URL.Query()
		field, err := strconv.Atoi(query.Get("id"))
		if err != nil {
			response.WrongParameter()
			return
		}

		currentUserId := userProvider.GetUserIdFromSession(r)
		if currentUserId == 0 {
			log.Error("Unauthorized")
			response.Unauthorized()
			return
		}

		userCanSwith := userProvider.CanUserSwitchToId(currentUserId, field)
		// todo: проверяем что этот пользователь может переключиться на указанного

		// todo: меняем айдишник юзера в сессии на указаный

		// todo: сохраняем в сессии старый айди чтоб потом вернуться

		//// TODO: тут еще надо проверить, что пользователи потомки одного родителя ???
		//userSwitcher, err := userProvider.GetUsersToSwitch(currentUserId)
		//if err != nil {
		//	log.Errorf("Error: %s", err.Error())
		//	response.InternalServerError()
		//	return
		//}

		//response.OKWithData(userSwitcher)

	}
}
