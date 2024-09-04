package user

import (
	"github.com/dmidokov/rv2/lib"
	"github.com/dmidokov/rv2/lib/entitie"
	resp "github.com/dmidokov/rv2/response"
	"github.com/dmidokov/rv2/session/cookie"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type userSwitcher interface {
	GetUserIdFromSession(r *http.Request) int
	GetById(userId int) (*entitie.User, error)
	GetUsersToSwitch(userId int) ([]*entitie.UserSwitcher, error)
	CanUserSwitchToId(from int, to int) bool
	GetParentId(userId int) (int, error)
}

type SessionStorage interface {
	Save(r *http.Request, w http.ResponseWriter, data map[string]interface{}) bool
	Get(r *http.Request) (map[interface{}]interface{}, error)
}

func (s *Service) SwitchUser(
	userProvider userSwitcher,
	rightsProvider rightsSetter,
	sessionProvider SessionStorage,
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
			log.Warning("Unauthorized")
			response.Unauthorized()
			return
		}

		if !isSwitchAllow(userProvider, currentUserId, field, log, rightsProvider) {
			response.NotAllowed()
			return
		}

		var savingParams = make(map[string]interface{}, 1)
		savingParams[cookie.SwitchedTo] = field

		sessionProvider.Save(r, w, savingParams)

		response.OK()

	}
}

func isSwitchAllow(userProvider userSwitcher, currentUserId int, field int, log *logrus.Logger, rightsProvider rightsSetter) bool {
	canSwitch := userProvider.CanUserSwitchToId(currentUserId, field)

	if !canSwitch {
		log.Warning("Hasn't rights to switch, no relations")
		return false
	}

	currentUser, _ := userProvider.GetById(currentUserId)

	if !rightsProvider.CheckUserRight(currentUser, lib.HotSwitchToAnotherUser) {
		log.Warning("Hasn't rights to switch, no right")
		return false
	}

	currentUserParent, _ := userProvider.GetParentId(currentUserId)
	switchUserParent, _ := userProvider.GetParentId(field)

	if currentUserParent != switchUserParent {
		log.Warning("Users has different parent")
		return false
	}

	return true
}
