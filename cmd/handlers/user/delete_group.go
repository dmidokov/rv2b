package user

import (
	"github.com/dmidokov/rv2/lib"
	"github.com/dmidokov/rv2/lib/request"
	resp "github.com/dmidokov/rv2/response"
	"net/http"
)

func (s *Service) DeleteGroup(
	userProvider userSwitchSetter,
	rightsProvider rightsSetter,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method == http.MethodOptions {
			return
		}

		method := "api.user.DeleteUserGroup"
		s.Logger.Info("Start method: ", method)

		response := resp.Service{Writer: &w, Logger: s.Logger, Operation: method}

		request := request.New().Parse(r)
		s.Logger.Infof("Request: %+v", request)

		requestUserId, err := request.GetInt("userId")
		if err != nil {
			s.Logger.Error("Invalid params")
			response.WrongParameter()
			return
		}
		requestGroupId, err := request.GetInt("groupId")
		if err != nil {
			s.Logger.Error("Invalid params")
			response.WrongParameter()
			return
		}

		currentUserId := userProvider.GetUserIdFromSession(r)
		if currentUserId == 0 {
			s.Logger.Error("Unauthorized")
			response.Unauthorized()
			return
		}

		currentUser, err := userProvider.GetById(currentUserId)
		if err != nil {
			s.Logger.Errorf("Error: %s", err.Error())
			response.InternalServerError()
			return
		}

		if rightsProvider.CheckUserRight(currentUser, lib.UnassignUserGroup) {
			err := rightsProvider.UnassignUserGroup(requestUserId, requestGroupId)
			if err != nil {
				s.Logger.Errorf("Can't unassign group to user %s", err.Error())
				response.InternalServerError()
				return
			}
		}

		response.OK()

	}
}
