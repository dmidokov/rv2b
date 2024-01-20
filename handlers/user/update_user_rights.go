package user

import (
	"encoding/json"
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

func (s *Service) Update(userProvider userRightsUpdater, rightsProvider rightsSetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			return
		}

		_ = s.Logger
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

		switch field {
		case "rights":
			updateRights(response, request, r, userProvider, rightsProvider)
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

	if !rightsProvider.CheckUserRight(currentUser, lib.EditUserRights) {
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
		response.InternalServerError()
		return
	}

	response.OK()
}
