package auth

import (
	"encoding/json"
	e "github.com/dmidokov/rv2/entitie"
	resp "github.com/dmidokov/rv2/response"
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

type ErrorResponse struct {
	resp.Response
}

type OrganizationProvider interface {
	GetByHostName(hostName string) (*e.Organization, error)
}

type UserProvider interface {
	GetUserByLoginAndOrganization(login string, organizationId int) (*e.User, error)
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

		if err != nil {
			s.Logger.WithFields(
				logrus.Fields{
					"req": req,
				}).
				Errorf("При декодировании данных авторизации произошла ошибка: %s", err.Error())

			json.NewEncoder(w).Encode(
				resp.Error("DecodeError"))

			return
		}

		contextLogger := s.Logger.WithFields(logrus.Fields{
			"fn":      fn,
			"request": req,
		})

		if req.UserPass == "" || req.UserName == "" {
			contextLogger.Errorf("Один из переданных параметров пустой")

			json.NewEncoder(w).Encode(
				resp.Error("OneOfTheSpecifiedParametersIsEmpty"))

			return
		}

		foundOrganization, err := organizationProvider.GetByHostName(r.Host)

		if err != nil {
			if err == pgx.ErrNoRows {
				contextLogger.WithFields(
					logrus.Fields{
						"host": r.Host,
					}).Errorf("Организация не найдена: %s", err.Error())

				json.NewEncoder(w).Encode(
					resp.Error("OrganizationNotFound"))

				return
			}

			contextLogger.WithFields(
				logrus.Fields{
					"host": r.Host,
				}).Errorf("Ошибка БД: %s", err.Error())

			json.NewEncoder(w).Encode(
				resp.Error("DatabaseError"))

			return
		}

		login, password := prepareLoginAndPassword(req.UserName, req.UserPass)

		user, err := userProvider.GetUserByLoginAndOrganization(login, foundOrganization.Id)
		if err != nil {
			if err == pgx.ErrNoRows {
				contextLogger.Errorf("Пользователь не найден: %s", err.Error())

				json.NewEncoder(w).Encode(
					resp.Error("UserNotFound"))

				return
			}

			contextLogger.Errorf("Ошибка БД: %s", err.Error())

			json.NewEncoder(w).Encode(
				resp.Error("DatabaseError"))

			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
			contextLogger.Errorf("Неверный пароль: %s", err.Error())

			json.NewEncoder(w).Encode(
				resp.Error("UserNotFound"))

			return
		}

		var savingParams = make(map[string]interface{}, 3)

		savingParams["authenticated"] = true
		savingParams["userid"] = user.Id
		savingParams["organizationid"] = foundOrganization.Id
		s.CookieStore.SetMaxAge(s.Config.SessionMaxAge)

		s.CookieStore.Save(r, w, savingParams)

		if err != nil {
			contextLogger.Errorf("Ошибка сохранения сессии: %s", err.Error())

			json.NewEncoder(w).Encode(
				resp.Error("SessionSaveError"))

			return
		}

		json.NewEncoder(w).Encode(
			resp.OK())
	}
}

func prepareLoginAndPassword(login, password string) (string, string) {
	return strings.Trim(login, " "), strings.Trim(password, " ")
}
