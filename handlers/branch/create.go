package branch

import (
	"encoding/json"
	"github.com/dmidokov/rv2/lib"
	e "github.com/dmidokov/rv2/lib/entitie"
	resp "github.com/dmidokov/rv2/response"
	"github.com/dmidokov/rv2/storage/postgres/rights"
	"net/http"
	"strings"
)

type branchCreator interface {
	Create(org *e.Branch, userId int) (*e.Branch, error)
}

type CreateBranchRequest struct {
	Name     string `json:"name"`
	Address  string `json:"address"`
	Phone    string `json:"phone"`
	WorkTime string `json:"workTime"`
}

func (s *Service) Create(branchCreator branchCreator, userProvider userProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			return
		}

		log := s.Logger
		method := "api.branch.create"
		response := resp.Service{Writer: &w, Logger: s.Logger, Operation: method}

		var orgData CreateBranchRequest
		err := json.NewDecoder(r.Body).Decode(&orgData)

		if err != nil {
			errorText := "При декодировании данных авторизации произошла ошибка"
			log.Errorf(errorText+": %s", err.Error())
			http.Error(w, errorText, http.StatusInternalServerError)
			return
		}

		rightsService := rights.New(s.DB, s.Logger)

		currentUserId := userProvider.GetUserIdFromSession(r)
		if currentUserId == 0 {
			response.Unauthorized()
			return
		}

		currentUserOrganizationId := userProvider.GetOrganizationIdFromSession(r)
		if currentUserOrganizationId == 0 {
			response.Unauthorized()
			return
		}

		newBranch := &e.Branch{
			OrgId:    currentUserOrganizationId,
			Name:     strings.Trim(orgData.Name, " "),
			Address:  strings.Trim(orgData.Address, " "),
			Phone:    strings.Trim(orgData.Phone, " "),
			WorkTime: strings.Trim(orgData.WorkTime, " "),
		}

		currentUser, err := userProvider.GetById(currentUserId)
		if err != nil {
			response.InternalServerError()
			return
		}

		if rightsService.CheckUserRight(currentUser, lib.AddBranch) {

			_, err := branchCreator.Create(newBranch, currentUserId)
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
