package branch

import (
	"github.com/dmidokov/rv2/lib"
	resp "github.com/dmidokov/rv2/response"
	"github.com/dmidokov/rv2/storage/postgres/rights"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type DeleteProvider interface {
	Delete(branchId, orgId int) error
}

func (s *Service) DeleteBranch(branchProvider DeleteProvider, userProvider userProvider) http.HandlerFunc {
	// todo: add logs, add tests positive and negative, user cannot delete branch if it not part of it
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			return
		}

		log := s.Logger
		method := "api.branch.delete"
		response := resp.Service{Writer: &w, Logger: s.Logger, Operation: method}

		vars := mux.Vars(r)
		varsId, ok := vars["id"]

		if !ok {
			response.EmptyData()
			return
		}

		branchId, err := strconv.Atoi(varsId)
		if err != nil {
			response.InternalServerError()
			return
		}

		userId := userProvider.GetUserIdFromSession(r)
		if userId == 0 {
			response.Unauthorized()
			return
		}

		currentUser, err := userProvider.GetById(userId)
		if err != nil {
			response.InternalServerError()
			return
		}

		rightsProvider := rights.New(s.DB, s.Logger)
		if rightsProvider.CheckUserRight(currentUser, lib.DeleteBranch) {
			err = branchProvider.Delete(branchId, currentUser.OrganizationId)
			if err != nil {
				log.Errorf("Error: %s", err.Error())

				response.InternalServerError()
				return
			}
			response.OK()
			return
		}
		response.NotAllowed()
	}
}
