package user

import (
	"github.com/dmidokov/rv2/handlers/navigation"
	"github.com/dmidokov/rv2/lib/entitie"
	resp "github.com/dmidokov/rv2/response"
	"github.com/dmidokov/rv2/storage/postgres/rights"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type navigationProvider interface {
	Get(userId int, navigationType string) (*[]entitie.Navigation, error)
}

type userGetter interface {
	GetById(userId int) (*entitie.User, error)
	GetOrganizationIdFromSession(r *http.Request) int
	GetUserIdFromSession(r *http.Request) int
	GetByOrganizationId(userId int) ([]*entitie.UserShort, error)
	GetInfo(userId int, infoLevel int) (*entitie.UserInfoFull, error)
	GetParentId(userId int) (int, error)
	GetChild(userId int) ([]*entitie.UserShort, error)
	GetHotSwitch(userId int) ([]*entitie.UserShort, error)
}

type rightsProvider interface {
	GetGroupsByOrganizationId(organizationId int) ([]entitie.Group, error)
	GetUserGroupsWithName(userId int) ([]entitie.GroupNameAndIds, error)
}

func (s *Service) GetUser(userProvider userGetter, navigationProvider navigationProvider, rightsProvider rightsProvider) http.HandlerFunc {
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

		organizationId := userProvider.GetOrganizationIdFromSession(r)
		fullUserInfo.Groups, err = getGroupList(rightsProvider, organizationId)
		if err != nil {
			log.Errorf("Error: %s", err.Error())
			response.InternalServerError()
			return
		}

		fullUserInfo.AssignedGroups, err = rightsProvider.GetUserGroupsWithName(userId)
		if err != nil {
			log.Errorf("Error: %s", err.Error())
			response.InternalServerError()
			return
		}

		response.OKWithData(fullUserInfo)
	}
}

func getGroupList(provider rightsProvider, organizationId int) ([]entitie.Group, error) {
	return provider.GetGroupsByOrganizationId(organizationId)
}

func getHotSwitchFromUser(userProvider userGetter, userId int) ([]entitie.UserIdAndLogin, error) {
	hotSwitch, err := userProvider.GetHotSwitch(userId)
	if err != nil {
		return nil, err
	}

	var hotSwitchList []entitie.UserIdAndLogin

	for _, item := range hotSwitch {
		hotSwitchList = append(hotSwitchList, entitie.ConvertUserToUserLogin(*item))
	}
	return hotSwitchList, nil

}

func getCurrentUserRightsWithDescription(userProvider userGetter, currentUserId int, rightsService *rights.Service) ([]entitie.Right, error) {

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

func getUserNavigation(navigationProvider navigationProvider, currentUserId int, userId int) ([]entitie.NavigationInfoPage, error) {
	currentUserNavigation, err := navigationProvider.Get(currentUserId, navigation.TypeNavigationLeft)
	if err != nil {
		return nil, err
	}

	searchUserAvailableNavigation, err := navigationProvider.Get(userId, navigation.TypeNavigationLeft)
	if err != nil {
		return nil, err
	}

	var shortNavigationInfo []entitie.NavigationInfoPage

	for _, v := range *currentUserNavigation {
		elem := entitie.NavigationInfoPage{
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

func getChildList(userProvider userGetter, userId int) ([]entitie.UserIdAndLogin, error) {

	parentId, err := userProvider.GetParentId(userId)
	if err != nil {
		return nil, err
	}

	childUsers, err := userProvider.GetChild(parentId)
	if err != nil {
		return nil, err
	}

	var childUsersAndLoginList []entitie.UserIdAndLogin

	for _, item := range childUsers {
		childUsersAndLoginList = append(childUsersAndLoginList, entitie.ConvertUserToUserLogin(*item))
	}
	return childUsersAndLoginList, nil
}
