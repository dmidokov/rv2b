package organization

import (
	e "github.com/dmidokov/rv2/entitie"
	"github.com/dmidokov/rv2/response"
	"github.com/dmidokov/rv2/rights"
	"github.com/sirupsen/logrus"
	"net/http"
)

type CreateOrganizationRequest struct {
	Name     string `json:"name"`
	Host     string `json:"host"`
	UserName string `json:"user-name"`
	UserPass string `json:"user-pass"`
}

type ShortResponse struct {
	Name string `json:"name"`
	Host string `json:"host"`
}

type OrgGetter interface {
	GetAll() ([]*e.Organization, error)
}

type userProvider interface {
	GetById(userId int) (*e.User, error)
	GetOrganizationIdFromSession(r *http.Request) int
	GetUserIdFromSession(r *http.Request) int
}

// Get возвращает список организаций с проверкой права по просмотр организаций
func (s *Service) Get(orgProvider OrgGetter, userProvider userProvider) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			return
		}

		fn := "api.organizations.get"
		contextLogger := s.Logger.WithFields(logrus.Fields{
			"fn": fn,
		})

		response := resp.Service{Writer: &w, Logger: s.Logger, Operation: fn}

		currentUserId := userProvider.GetUserIdFromSession(r)
		if currentUserId == 0 {
			contextLogger.Warning("Пользователь не найден в сессии")
			response.Unauthorized()
			return
		}

		currentUser, err := userProvider.GetById(currentUserId)
		if err != nil {
			contextLogger.Warning("Пользователь не найден")
			response.Unauthorized()
			return
		}

		if !rights.New().CheckUserRight(currentUser.Rights, rights.ViewOrganization) {
			contextLogger.Warning("Недостаточно прав")
			response.NotAllowed()
			return
		}

		items, err := orgProvider.GetAll()
		if err != nil {
			contextLogger.Errorf("Не удалось получить список организаций: Error: %s", err.Error())

			response.InternalServerError()
			return
		}

		response.OKWithData(items)
	}
}
