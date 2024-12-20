package navigation

import (
	e "github.com/dmidokov/rv2/lib/entitie"
	resp "github.com/dmidokov/rv2/response"
	"net/http"
)

type userProvider interface {
	GetUserIdFromSession(r *http.Request) int
}

type navigationProvider interface {
	Get(userId int, navigationType string) (*[]e.Navigation, error)
}

func (s *Service) GetNavigation(userProvider userProvider, navigationProvider navigationProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			return
		}

		log := s.Logger
		method := "api.navigation.get"
		response := resp.New(&w, s.Logger, method)

		//navigationService := navigation.Service{DB: s.DB, CookieStore: s.CookieStore, Log: s.Logger}
		//userService := user.Service{DB: s.DB, CookieStore: s.CookieStore, Log: s.Logger}

		userId := userProvider.GetUserIdFromSession(r)
		if userId == 0 {
			response.Unauthorized()
			return
		}

		navigationItems, err := navigationProvider.Get(userId, TypeNavigationLeft)
		if err != nil {
			log.Errorf("Error: %s", err.Error())

			response.InternalServerError()
			return
		}

		response.OKWithData(navigationItems)

	}
}
