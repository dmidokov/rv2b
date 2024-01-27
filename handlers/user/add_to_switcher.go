package user

import (
	"encoding/json"
	"github.com/dmidokov/rv2/lib"
	"github.com/dmidokov/rv2/lib/entitie"
	resp "github.com/dmidokov/rv2/response"
	"net/http"
)

type AddToSwitcherRequest struct {
	FromId int `json:"fromId,omitempty"`
	ToId   int `json:"toId,omitempty"`
}

type userSwitchSetter interface {
	GetUserIdFromSession(r *http.Request) int
	GetById(userId int) (*entitie.User, error)
	SetHotSwitchRelation(fromId, toId int) error
}

func (s *Service) AddToSwitcher(
	userProvider userSwitchSetter,
	rightsProvider rightsSetter,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method == http.MethodOptions {
			return
		}

		log := s.Logger
		method := "api.user.UpdateUserRights"

		response := resp.Service{Writer: &w, Logger: s.Logger, Operation: method}

		log.Info(r.Body)

		request := AddToSwitcherRequest{}
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			log.Errorf("Error: %s", err.Error())
			response.JsonDecodeError()
			return
		}

		currentUserId := userProvider.GetUserIdFromSession(r)
		if currentUserId == 0 {
			log.Error("Unauthorized")
			response.Unauthorized()
			return
		}

		currentUser, err := userProvider.GetById(currentUserId)
		if err != nil {
			log.Errorf("Error: %s", err.Error())
			response.InternalServerError()
		}

		if !rightsProvider.CheckUserRight(currentUser, lib.EditUserHotSwitch) {
			log.Errorf("Method not allowed")
			response.NotAllowed()
			return
		}

		// TODO: тут еще надо проверить, что пользователи потомки одного родителя
		err = userProvider.SetHotSwitchRelation(request.FromId, request.ToId)
		if err != nil {
			log.Errorf("Error: %s", err.Error())
			response.InternalServerError()
			return
		}

	}
}
