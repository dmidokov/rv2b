package branch

import (
	resp "github.com/dmidokov/rv2/response"
	"github.com/gorilla/mux"
	"net/http"
)

type SessionStorage interface {
	Save(r *http.Request, w http.ResponseWriter, data map[string]interface{}) bool
	Get(r *http.Request, key string) (interface{}, bool)
	SetMaxAge(maxAge int)
}

func (s *Service) SetActiveBranch(userProvider userProvider, cookieStorage SessionStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			return
		}

		log := s.Logger
		method := "api.branch.setActive"
		response := resp.Service{Writer: &w, Logger: s.Logger, Operation: method}

		vars := mux.Vars(r)
		branchId, ok := vars["branchId"]

		if !ok {
			log.Error("Empty data: branchId")
			response.EmptyData()
			return
		}

		if !userProvider.IsAuthorized(r) {
			response.Unauthorized()
		}

		var savingParams = make(map[string]interface{}, 1)
		savingParams["selected_branch"] = branchId
		if !cookieStorage.Save(r, w, savingParams) {
			response.InternalServerError()
		}

	}
}
