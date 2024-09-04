package auth

import (
	"encoding/json"
	"fmt"
	"github.com/dmidokov/rv2/lib/entitie"
	resp "github.com/dmidokov/rv2/response"
	"github.com/dmidokov/rv2/session/cookie"
	"github.com/jackc/pgx/v4"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
)

type Request struct {
	UserName string `json:"username"`
	UserPass string `json:"user_pass"`
}

type Response struct {
	resp.Response
}

type OkData struct {
	StartPage string `json:"startPage"`
}

type ErrorResponse struct {
	resp.Response
}

type OrganizationProvider interface {
	GetByHostName(hostName string) (*entitie.Organization, error)
}

type UserProvider interface {
	GetUserByLoginAndOrganization(login string, organizationId int) (*entitie.User, error)
}

func (s *Service) SignIn(userProvider UserProvider, organizationProvider OrganizationProvider) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		//TODO: замаскировать в логе пароли
		fn := "api.signin"

		w.Header().Set("Content-Type", "application/json")

		if r.Method == http.MethodOptions {
			return
		}

		var req Request
		err := json.NewDecoder(r.Body).Decode(&req)

		response := resp.New(&w, s.Logger, fn)

		if err != nil {
			s.Logger.
				WithFields(
					logrus.Fields{
						"req": req,
					}).
				Errorf("При декодировании данных авторизации произошла ошибка: %s", err.Error())

			response.WithError("DecodeError")

			return
		}

		contextLogger := s.Logger.WithFields(logrus.Fields{
			"fn":      fn,
			"request": req,
		})

		if req.UserPass == "" || req.UserName == "" {
			contextLogger.Errorf("Один из переданных параметров пустой")

			response.WithError("OneOfTheSpecifiedParametersIsEmpty")

			return
		}

		foundOrganization, err := organizationProvider.GetByHostName(r.Host)

		if err != nil {
			if err == pgx.ErrNoRows {
				contextLogger.WithFields(
					logrus.Fields{
						"host": r.Host,
					}).Errorf("Организация не найдена: %s", err.Error())

				response.WithError("OrganizationNotFound")

				return
			}

			contextLogger.WithFields(
				logrus.Fields{
					"host": r.Host,
				}).Errorf("Ошибка БД: %s", err.Error())

			response.WithError("DatabaseError")

			return
		}

		login, password := prepareLoginAndPassword(req.UserName, req.UserPass)

		user, err := userProvider.GetUserByLoginAndOrganization(login, foundOrganization.Id)
		if err != nil {
			if err == pgx.ErrNoRows {
				contextLogger.Errorf("Пользователь не найден: %s", err.Error())

				response.WithError("UserNotFound")

				return
			}

			contextLogger.Errorf("Ошибка БД: %s", err.Error())

			response.WithError("DatabaseError")

			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
			contextLogger.Errorf("Неверный пароль: %s", err.Error())

			response.WithError("UserNotFound")

			return
		}

		var savingParams = make(map[string]interface{}, 3)

		fmt.Println(user.Id)

		savingParams[cookie.Authenticated] = true
		savingParams[cookie.UserId] = user.Id
		savingParams[cookie.OrganizationId] = foundOrganization.Id
		s.CookieStore.SetMaxAge(s.Config.SessionMaxAge)

		s.CookieStore.Save(r, w, savingParams)

		if err != nil {
			contextLogger.Errorf("Ошибка сохранения сессии: %s", err.Error())

			response.WithError("SessionSaveError")

			return
		}

		response.OKWithData(OkData{
			StartPage: user.StartPage,
		})
	}
}

func prepareLoginAndPassword(login, password string) (string, string) {
	return strings.Trim(login, " "), strings.Trim(password, " ")
}
