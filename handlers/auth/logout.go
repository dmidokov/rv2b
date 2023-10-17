package auth

import (
	"encoding/json"
	resp "github.com/dmidokov/rv2/response"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (s *Service) Logout(w http.ResponseWriter, r *http.Request) {
	fn := "api.logout"

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

		_ = json.NewEncoder(w).Encode(
			resp.Error("SessionSaveError"))

		return
	}

	http.Redirect(w, r, "/#/login", http.StatusMovedPermanently)
}
