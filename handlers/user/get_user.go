package user

import (
	resp "github.com/dmidokov/rv2/response"
	"github.com/dmidokov/rv2/storage/postgres/rights"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func (s *Service) GetUser(userProvider userProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			return
		}

		log := s.Logger
		method := "api.users.get"
		response := resp.Service{Writer: &w, Logger: s.Logger, Operation: method}
		log.Info(method)

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

		/**
		TODO:
			- тут надо проверить что у пользователя есть право на просмотр конкретного пользователя (как?)
			- еще надо подготовить правильный ответ без пароля и прочей информация, которую отдавать нельзя
		*/
		//item, err := userProvider.GetById(userId)
		item, err := userProvider.GetInfo(userId, 1)
		if err != nil {
			log.Errorf("Error: %s", err.Error())

			response.InternalServerError()
			return
		}

		currentUserId := userProvider.GetUserIdFromSession(r)
		currentUser, err := userProvider.GetById(currentUserId)

		if err != nil {
			log.Errorf("Error: %s", err.Error())

			response.InternalServerError()
			return
		}

		rightsService := rights.New(s.DB, s.Logger)
		currentUserRightsWithDescriptions, err := rightsService.GetByUserRights(currentUser.Rights)

		if err != nil {
			log.Errorf("Error: %s", err.Error())

			response.InternalServerError()
			return
		}

		item.UserRightsWithDescriptions = *currentUserRightsWithDescriptions
		response.OKWithData(item)
	}
}
