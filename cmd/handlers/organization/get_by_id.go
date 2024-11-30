package organization

import (
	e "github.com/dmidokov/rv2/lib/entitie"
	resp "github.com/dmidokov/rv2/response"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type GetByIdProvider interface {
	GetById(id int) (*e.Organization, error)
	GetByHostName(hostName string) (*e.Organization, error)
}

type AuthCheckProvider interface {
	IsAuthorized(r *http.Request) bool
}

func (s *Service) GetById(org GetByIdProvider, authCheckProvider AuthCheckProvider) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			return
		}

		method := "api.organizations.get"

		response := resp.Service{Writer: &w, Logger: s.Logger, Operation: method}

		vars := mux.Vars(r)
		varsId, ok := vars["id"]

		if !ok {
			response.EmptyData()
			return
		}

		var item *e.Organization
		var err error

		switch varsId {
		case "current":
			item, err = org.GetByHostName(r.Host)
			if err != nil {
				response.InternalServerError()
				return
			}
		default:
			id, ok := strconv.Atoi(varsId)
			if ok != nil {
				response.WrongParameter()
				return
			}
			item, err = org.GetById(id)
			if err != nil {
				response.InternalServerError()
				return
			}
		}

		if !authCheckProvider.IsAuthorized(r) {
			shortResponse := ShortResponse{Name: item.Name, Host: item.Host}
			response.OKWithData(shortResponse)
			return
		}

		response.OKWithData(item)

	}
}
