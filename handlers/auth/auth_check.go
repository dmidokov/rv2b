package auth

import (
	"github.com/dmidokov/rv2/response"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (s *Service) AuthCheck(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodOptions {
		return
	}

	fn := "api.AuthCheck"

	contextLogger := s.Logger.WithFields(logrus.Fields{
		"fn": fn,
	})

	response := resp.New(&w, s.Logger, fn)

	if auth, ok := s.CookieStore.Get(r, "authenticated"); ok && auth.(bool) {
		response.OK()

		return
	} else {
		contextLogger.Warning("User is not authorized")
		response.Unauthorized()

		return
	}
}
