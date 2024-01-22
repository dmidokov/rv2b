package user

import (
	e "github.com/dmidokov/rv2/lib/entitie"
	resp "github.com/dmidokov/rv2/response"
	"github.com/dmidokov/rv2/storage/postgres/rights"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type navigationProvider interface {
	Get(userId int) (*[]e.Navigation, error)
}

type userGetter interface {
	GetById(userId int) (*e.User, error)
	GetOrganizationIdFromSession(r *http.Request) int
	GetUserIdFromSession(r *http.Request) int
	GetByOrganizationId(userId int) ([]*e.UserShort, error)
	GetInfo(userId int, infoLevel int) (*e.UserInfoFull, error)
	GetParentId(userId int) (int, error)
	GetChild(userId int) ([]*e.UserShort, error)
}

func (s *Service) GetUser(userProvider userGetter, navigationProvider navigationProvider) http.HandlerFunc {
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
			log.Errorf("Empty data")
			response.EmptyData()
			return
		}

		userId, err := strconv.Atoi(varsId)
		if err != nil {
			log.Errorf("Can't conver string to int %s", err.Error())
			response.InternalServerError()
			return
		}

		/**
		TODO:
			- тут надо проверить что у пользователя есть право на просмотр конкретного пользователя (как?)
			- еще надо подготовить правильный ответ без пароля и прочей информация, которую отдавать нельзя
		*/
		fullUserInfo, err := userProvider.GetInfo(userId, 1)
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

		fullUserInfo.UserRightsWithDescriptions = *currentUserRightsWithDescriptions
		currentUserNavigation, err := navigationProvider.Get(currentUserId)
		if err != nil {
			log.Errorf("Error: %s", err.Error())

			response.InternalServerError()
			return
		}

		searchUserAvailableNavigation, err := navigationProvider.Get(userId)
		if err != nil {
			log.Errorf("Error: %s", err.Error())

			response.InternalServerError()
			return
		}

		var shortNavigationInfo []e.NavigationInfoPage

		for _, v := range *currentUserNavigation {
			elem := e.NavigationInfoPage{
				Id:      v.Id,
				Title:   v.Title,
				Enabled: false,
			}
			for _, vv := range *searchUserAvailableNavigation {
				if v.Id == vv.Id {
					elem.Enabled = true
				}
			}
			shortNavigationInfo = append(shortNavigationInfo, elem)
		}

		fullUserInfo.Navigation = shortNavigationInfo

		parentId, err := userProvider.GetParentId(userId)
		if err != nil {
			log.Errorf("Error: %s", err.Error())
			response.InternalServerError()
			return
		}

		childUsers, err := userProvider.GetChild(parentId)
		if err != nil {
			log.Errorf("Error: %s", err.Error())
			response.InternalServerError()
			return
		}

		var childUsersAndLoginList []e.UserIdAndLogin

		for _, item := range childUsers {
			childUsersAndLoginList = append(childUsersAndLoginList, e.ConvertUserToUserLogin(*item))
		}

		fullUserInfo.Childs = childUsersAndLoginList

		response.OKWithData(fullUserInfo)
	}
}
