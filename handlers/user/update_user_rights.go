package user

import (
	"encoding/json"
	"github.com/dmidokov/rv2/lib"
	e "github.com/dmidokov/rv2/lib/entitie"
	resp "github.com/dmidokov/rv2/response"
	"github.com/sirupsen/logrus"
	"math"
	"net/http"
)

type UpdateUserRightsRequest struct {
	UserId int  `json:"userId,omitempty"`
	Value  int  `json:"value,omitempty"`
	Set    bool `json:"set,omitempty"`
}

type rightsSetter interface {
	CheckUserRight(user *e.User, right int) bool
}

type userRightsUpdater interface {
	GetById(userId int) (*e.User, error)
	GetUserIdFromSession(r *http.Request) int
	UpdateUser(user *e.User) (*e.User, error)
}

type navigationUpdater interface {
	Set(userId int, navigationId int, groupId int) (*e.NavigationAvailable, error)
	Delete(userId int, navigationId int, groupId int) error
}

func (s *Service) Update(
	userProvider userRightsUpdater,
	rightsProvider rightsSetter,
	navigationProvider navigationUpdater,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			return
		}

		log := s.Logger
		method := "api.user.UpdateUserRights"

		response := resp.Service{Writer: &w, Logger: s.Logger, Operation: method}

		request := UpdateUserRightsRequest{}
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			response.JsonDecodeError()
			return
		}

		query := r.URL.Query()
		field := query.Get("field")
		if field == "" {
			response.WrongParameter()
			return
		}

		s.Logger.Info(field)

		switch field {
		case "rights":
			updateRights(response, request, r, userProvider, rightsProvider, log)
		case "navigation":
			updateNavigation(response, request, r, userProvider, rightsProvider, navigationProvider, log)
		default:
			response.WrongParameter()
		}

	}
}

func updateRights(
	response resp.Service,
	request UpdateUserRightsRequest,
	r *http.Request,
	userProvider userRightsUpdater,
	rightsProvider rightsSetter,
	logger *logrus.Logger,
) {
	currentUserId := userProvider.GetUserIdFromSession(r)
	if currentUserId == 0 {
		logger.Info("unauthorized")
		response.Unauthorized()
		return
	}

	currentUser, err := userProvider.GetById(currentUserId)
	if err != nil {
		logger.Error("internal server error")
		response.InternalServerError()
	}

	userToUpdate, err := userProvider.GetById(request.UserId)
	if err != nil {
		logger.Info("user not found")
		response.UserNotFound()
		return
	}

	if !rightsProvider.CheckUserRight(currentUser, lib.EditUserRights) {
		logger.Info("not allowed")
		response.NotAllowed()
		return
	}

	if request.Set {
		userToUpdate.Rights = request.Value | userToUpdate.Rights
	} else {
		userToUpdate.Rights = (math.MaxInt ^ request.Value) & userToUpdate.Rights
	}

	_, err = userProvider.UpdateUser(userToUpdate)
	if err != nil {
		logger.Error("internal server error")
		response.InternalServerError()
		return
	}

	response.OK()
}
func updateNavigation(
	response resp.Service,
	request UpdateUserRightsRequest,
	r *http.Request,
	userProvider userRightsUpdater,
	rightsProvider rightsSetter,
	navigationProvider navigationUpdater,
	logger *logrus.Logger,
) {
	currentUserId := userProvider.GetUserIdFromSession(r)
	if currentUserId == 0 {
		response.Unauthorized()
		return
	}

	currentUser, err := userProvider.GetById(currentUserId)
	if err != nil {
		response.InternalServerError()
	}

	userToUpdate, err := userProvider.GetById(request.UserId)
	if err != nil {
		response.UserNotFound()
		return
	}

	// TODO: добавить право на редактирование навигации и проверить что
	// пользователю можно дать эту навигацию, то есть что тот кто устанавливает
	// имеет у себя такое поле навишации
	if !rightsProvider.CheckUserRight(currentUser, lib.EditUserNavigation) {
		logger.Info("not allowed")
		response.NotAllowed()
		return
	}

	if request.Set {
		_, err = navigationProvider.Set(userToUpdate.Id, request.Value, 1)
	} else {
		err = navigationProvider.Delete(userToUpdate.Id, request.Value, 1)
	}

	if err != nil {
		response.InternalServerError()
		return
	}

	response.OK()
}
