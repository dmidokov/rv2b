package group

import (
	"github.com/dmidokov/rv2/lib"
	"github.com/dmidokov/rv2/lib/request"
	resp "github.com/dmidokov/rv2/response"
	"net/http"
)

type CreateGroupRequest struct {
	Name   string
	Rights int64
}

func (s *Service) AddGroup(userProvider userProvider, rightsProvider rightsProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			return
		}

		method := "api.AddGroups"
		response := resp.Service{Writer: &w, Logger: s.Logger, Operation: method}

		request := request.New().Parse(r)

		requestRights, err := request.GetInt64("rights")
		requestName := request.GetString("name")
		if err != nil {
			s.Logger.Warningf("Can't parse request: %s\n", err.Error())
			response.WrongParameter()
			return
		}

		s.Logger.Info("ADD groups")

		currentUserId := userProvider.GetUserIdFromSession(r)
		if currentUserId == 0 {
			s.Logger.Warning("User not authorized")
			response.Unauthorized()
			return
		}

		currentUser, err := userProvider.GetById(currentUserId)
		if err != nil {
			s.Logger.Warning("User not authorized")
			response.InternalServerError()
			return
		}

		if !rightsProvider.CheckUserRight(currentUser, lib.CreateUserGroup) {
			s.Logger.Warning("User have no rights to create group")
			response.NotAllowed()
			return
		}

		err = rightsProvider.CreateGroup(requestName, requestRights, currentUser.OrganizationId, currentUserId)
		if err != nil {
			s.Logger.Warningf("can't get groups list: %s", err.Error())
			response.InternalServerError()
			return
		}

		group, err := rightsProvider.GetGroupByName(requestName, currentUser.OrganizationId)
		if err != nil {
			s.Logger.Warningf("can't get groups list: %s", err.Error())
			response.InternalServerError()
			return
		}

		response.OKWithData(group)
	}
}
