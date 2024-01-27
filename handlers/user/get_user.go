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
	GetHotSwitch(userId int) ([]*e.UserShort, error)
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

		rightsService := rights.New(s.DB, s.Logger)
		currentUserId := userProvider.GetUserIdFromSession(r)
		fullUserInfo.UserRightsWithDescriptions, err = getCurrentUserRightsWithDescription(userProvider, currentUserId, rightsService)
		if err != nil {
			log.Errorf("Error: %s", err.Error())
			response.InternalServerError()
			return
		}

		fullUserInfo.Navigation, err = getUserNavigation(navigationProvider, currentUserId, userId)
		if err != nil {
			log.Errorf("Error: %s", err.Error())
			response.InternalServerError()
			return
		}

		fullUserInfo.Childs, err = getChildList(userProvider, userId)
		if err != nil {
			log.Errorf("Error: %s", err.Error())
			response.InternalServerError()
			return
		}

		fullUserInfo.HotSwitch, err = getHotSwitchFromUser(userProvider, userId)
		if err != nil {
			log.Errorf("Error: %s", err.Error())
			response.InternalServerError()
			return
		}

		response.OKWithData(fullUserInfo)
	}
}

func getHotSwitchFromUser(userProvider userGetter, userId int) ([]e.UserIdAndLogin, error) {
	hotSwitch, err := userProvider.GetHotSwitch(userId)
	if err != nil {
		return nil, err
	}

	var hotSwitchList []e.UserIdAndLogin

	for _, item := range hotSwitch {
		hotSwitchList = append(hotSwitchList, e.ConvertUserToUserLogin(*item))
	}
	return hotSwitchList, nil

}

func getCurrentUserRightsWithDescription(userProvider userGetter, currentUserId int, rightsService *rights.Service) ([]e.Right, error) {

	currentUser, err := userProvider.GetById(currentUserId)

	if err != nil {
		return nil, err
	}

	currentUserRightsWithDescriptions, err := rightsService.GetByUserRights(currentUser.Rights)

	if err != nil {
		return nil, err
	}
	return *currentUserRightsWithDescriptions, nil
}

func getUserNavigation(navigationProvider navigationProvider, currentUserId int, userId int) ([]e.NavigationInfoPage, error) {
	currentUserNavigation, err := navigationProvider.Get(currentUserId)
	if err != nil {
		return nil, err
	}

	searchUserAvailableNavigation, err := navigationProvider.Get(userId)
	if err != nil {
		return nil, err
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
	return shortNavigationInfo, nil
}

func getChildList(userProvider userGetter, userId int) ([]e.UserIdAndLogin, error) {

	parentId, err := userProvider.GetParentId(userId)
	if err != nil {
		return nil, err
	}

	childUsers, err := userProvider.GetChild(parentId)
	if err != nil {
		return nil, err
	}

	var childUsersAndLoginList []e.UserIdAndLogin

	for _, item := range childUsers {
		childUsersAndLoginList = append(childUsersAndLoginList, e.ConvertUserToUserLogin(*item))
	}
	return childUsersAndLoginList, nil
}
