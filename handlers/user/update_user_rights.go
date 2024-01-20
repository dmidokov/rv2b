package user

import (
	"encoding/json"
	"fmt"
	"github.com/dmidokov/rv2/lib"
	e "github.com/dmidokov/rv2/lib/entitie"
	resp "github.com/dmidokov/rv2/response"
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
	//GetOrganizationIdFromSession(r *http.Request) int
	GetUserIdFromSession(r *http.Request) int
	UpdateUser(user *e.User) (*e.User, error)
	//GetByOrganizationId(orgId, userId int) ([]*e.UserShort, error)
	//Create(user *e.User) (int, error)
	//SetUserCreateRelations(creatorId int, createdId int) error
}

func (s *Service) UpdateUserRights(userProvider userRightsUpdater, rightsProvider rightsSetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			return
		}

		_ = s.Logger
		method := "api.user.UpdateUserRights"

		response := resp.Service{Writer: &w, Logger: s.Logger, Operation: method}

		request := UpdateUserRightsRequest{}
		err := json.NewDecoder(r.Body).Decode(&request)
		fmt.Println(r.Body)
		if err != nil {
			response.JsonDecodeError()
			return
		}

		currentUserId := userProvider.GetUserIdFromSession(r)
		if currentUserId == 0 {
			response.Unauthorized()
			return
		}

		currentUser, err := userProvider.GetById(currentUserId)
		if err != nil {
			response.InternalServerError()
		}

		if !rightsProvider.CheckUserRight(currentUser, lib.EditUserRights) {
			response.NotAllowed()
			return
		}

		userToUpdate, err := userProvider.GetById(request.UserId)
		if err != nil {
			response.UserNotFound()
			return
		}

		if request.Set {
			userToUpdate.Rights = request.Value | userToUpdate.Rights
		} else {
			userToUpdate.Rights = (math.MaxInt ^ request.Value) & userToUpdate.Rights
		}

		_, err = userProvider.UpdateUser(userToUpdate)
		if err != nil {
			response.InternalServerError()
			return
		}

		response.OK()

	}
}
