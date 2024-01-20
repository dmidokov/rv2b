package user

import (
	"github.com/dmidokov/rv2/lib"
	e "github.com/dmidokov/rv2/lib/entitie"
	resp "github.com/dmidokov/rv2/response"
	"github.com/dmidokov/rv2/storage/postgres/rights"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type userRemover interface {
	GetById(userId int) (*e.User, error)
	GetUserIdFromSession(r *http.Request) int
	Delete(userId int) error
}

func (s *Service) DeleteUser(userRemover userRemover) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			return
		}

		log := s.Logger
		method := "api.user.delete"
		response := resp.Service{Writer: &w, Logger: s.Logger, Operation: method}

		vars := mux.Vars(r)
		var userId int
		varsId, ok := vars["id"]

		if !ok {
			response.EmptyData()
			return
		}

		userId, err := strconv.Atoi(varsId)
		if err != nil {
			response.InternalServerError()
			return
		}

		rightsService := rights.New(s.DB, s.Logger)

		currentUserId := userRemover.GetUserIdFromSession(r)
		if currentUserId == 0 {
			response.Unauthorized()
			return
		}

		currentUser, err := userRemover.GetById(currentUserId)
		if err != nil {
			response.InternalServerError()
			return
		}

		if rightsService.CheckUserRight(currentUser, lib.DeleteUser) {
			err := userRemover.Delete(userId)
			if err != nil {
				log.Errorf("Error: %s", err.Error())
				response.InternalServerError()
				return
			}
		} else {
			log.Warningf("Method now allowed for user %d", currentUserId)
			response.NotAllowed()
			return
		}

		response.OK()
	}
}
