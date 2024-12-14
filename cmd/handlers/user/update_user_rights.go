package user

import (
	"encoding/json"
	"github.com/dmidokov/rv2/lib"
	"github.com/dmidokov/rv2/lib/entitie"
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
	CheckUserRight(user *entitie.User, right int) bool
	AssignUserGroup(userId int, groupId int) error
	UnassignUserGroup(userId int, groupId int) error
}

type userRightsUpdater interface {
	GetById(userId int) (*entitie.User, error)
	GetUserIdFromSession(r *http.Request) int
	UpdateUser(user *entitie.User) (*entitie.User, error)
}

type navigationUpdater interface {
	Set(userId int, navigationId int, groupId int) (*entitie.NavigationAvailable, error)
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

		s.Logger.Info("Field value is: " + field + ".")

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

	logger.Info("start rights update")

	currentUserId := userProvider.GetUserIdFromSession(r)
	if currentUserId == 0 {
		logger.Info("unauthorized")
		response.Unauthorized()
		return
	}

	logger.Info("get current user")
	currentUser, err := userProvider.GetById(currentUserId)
	if err != nil {
		logger.Error("internal server error")
		response.InternalServerError()
	}

	logger.Info("get user to update")
	userToUpdate, err := userProvider.GetById(request.UserId)
	if err != nil {
		logger.Info("user not found")
		response.UserNotFound()
		return
	}

	logger.Infof("check user rights\n user: %d, checked rights: %d", currentUser.Rights, lib.EditUserRights)
	if !rightsProvider.CheckUserRight(currentUser, lib.EditUserRights) {
		logger.Info("not allowed")
		response.NotAllowed()
		return
	}

	logger.Infof("Before user rights: %d", userToUpdate.Rights)
	if request.Set {
		logger.Info("Is request set")
		userToUpdate.Rights = request.Value | userToUpdate.Rights
	} else {
		logger.Info("Is request unset")
		userToUpdate.Rights = (math.MaxInt ^ request.Value) & userToUpdate.Rights
	}
	logger.Infof("After user rights: %d", userToUpdate.Rights)

	logger.Info("Update user")
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
