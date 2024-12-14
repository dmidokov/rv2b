package group

import (
	"github.com/dmidokov/rv2/lib"
	resp "github.com/dmidokov/rv2/response"
	"net/http"
	"strconv"
)

type DeleteGroupRequest struct {
	groupId int `json:"group_id"`
}

func (s *Service) DeleteGroup(userProvider userProvider, rightsProvider rightsProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			return
		}

		method := "api.DeleteGroups"
		response := resp.Service{Writer: &w, Logger: s.Logger, Operation: method}

		err := r.ParseForm()
		if err != nil {
			s.Logger.Warning("Invalid params")
			response.WrongParameter()
			return
		}

		var val string
		if val = r.Form.Get("group_id"); val == "" {
			s.Logger.Warning("Empty delete data, can't find group_id value")
			response.EmptyData()
			return
		}

		groupId, err := strconv.Atoi(val)
		if err != nil {
			s.Logger.Warning("Invalid group_id value")
			response.WrongParameter()
			return
		}

		currentUserId := userProvider.GetUserIdFromSession(r)
		if currentUserId == 0 {
			response.Unauthorized()
			return
		}

		currentUser, err := userProvider.GetById(currentUserId)
		if err != nil {
			s.Logger.Warning("User not authorized")
			response.InternalServerError()
			return
		}

		if !rightsProvider.CheckUserRight(currentUser, lib.DeleteUserGroups) {
			s.Logger.Warning("User have no rights to create group")
			response.NotAllowed()
			return
		}

		err = rightsProvider.DeleteGroup(groupId)
		if err != nil {
			s.Logger.Warningf("can't get groups list: %s", err.Error())
			response.InternalServerError()
			return
		}

		response.OK()
	}
}
