package organization

import (
	resp "github.com/dmidokov/rv2/response"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type DeleteProvider interface {
	Delete(orgId int) error
}

func (s *Service) DeleteOrganization(provider DeleteProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			return
		}

		log := *s.Logger
		method := "api.organizations.delete"
		response := resp.Service{Writer: &w, Logger: s.Logger, Operation: method}

		vars := mux.Vars(r)
		varsId, ok := vars["id"]

		if !ok {
			response.EmptyData()
			return
		}

		orgId, err := strconv.Atoi(varsId)
		if err != nil {
			response.InternalServerError()
			return
		}

		//organizationsService := organizations.Service{DB: s.DB, CookieStore: s.CookieStore, Log: s.Logger}

		err = provider.Delete(orgId)
		if err != nil {
			log.Errorf("Error: %s", err.Error())

			response.InternalServerError()
			return
		}

		response.OK()

	}
}
