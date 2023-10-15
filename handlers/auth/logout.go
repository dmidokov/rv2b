package auth

import (
	"encoding/json"
	resp "github.com/dmidokov/rv2/response"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (s *Service) Logout(w http.ResponseWriter, r *http.Request) {
	fn := "api.logout"

	//session, _ := s.CookieStore.Get(r, s.Config.SessionsSecret)
	//session.Values["authenticated"] = false
	//session.Values["userid"] = nil
	//session.Values["organizationid"] = nil
	//session.Options.MaxAge = -1
	//err := session.Save(r, w)

	var savingParams = make(map[string]interface{}, 3)

	savingParams["authenticated"] = false
	savingParams["userid"] = nil
	savingParams["organizationid"] = nil
	err := s.CookieStore.Save(r, w, savingParams)

	contextLogger := s.Logger.WithFields(logrus.Fields{
		"fn": fn,
	})

	if !err {
		contextLogger.Errorf("Ошибка сохранения сессии при логауте")

		json.NewEncoder(w).Encode(
			resp.Error("SessionSaveError"))

		return
	}

	http.Redirect(w, r, "/#/login", http.StatusMovedPermanently)
}
