package branch

import (
	e "github.com/dmidokov/rv2/entitie"
	resp "github.com/dmidokov/rv2/response"
	"github.com/dmidokov/rv2/rights"
	"github.com/sirupsen/logrus"
	"net/http"
)

type branchGetter interface {
	GetAll() ([]*e.Branch, error)
}

// Get возвращает список филиалов с проверкой права по просмотр
func (s *Service) Get(branchGetter branchGetter, userProvider userProvider) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			return
		}

		fn := "api.branches.get"
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

		if !rights.New().CheckUserRight(currentUser.Rights, rights.ViewBranchList) {
			contextLogger.Warning("Недостаточно прав")
			response.NotAllowed()
			return
		}

		items, err := branchGetter.GetAll()

		if err != nil {
			contextLogger.Errorf("Не удалось получить список филиалов: Error: %s", err.Error())

			response.InternalServerError()
			return
		}

		response.OKWithData(items)
	}
}
