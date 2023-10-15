package auth

import (
	"encoding/json"
	resp "github.com/dmidokov/rv2/response"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (s *Service) AuthCheck(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodOptions {
		return
	}

	fn := "api.authCheck"

	contextLogger := s.Logger.WithFields(logrus.Fields{
		"fn": fn,
	})

	if auth, ok := s.CookieStore.Get(r, "authenticated"); ok && auth.(bool) {
		json.NewEncoder(w).Encode(
			resp.OK())

		return
	} else {
		contextLogger.Warning("User is not authorized")
		json.NewEncoder(w).Encode(
			resp.Error("UserUnauthorized"))
		return
	}
}
