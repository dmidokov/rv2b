package organization

import (
	"encoding/json"
	e "github.com/dmidokov/rv2/entitie"
	resp "github.com/dmidokov/rv2/response"
	"github.com/dmidokov/rv2/rights"
	"github.com/dmidokov/rv2/storage/postgres/user"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
	"time"
)

type CreateOrganizationRequest struct {
	Name     string `json:"name"`
	Host     string `json:"host"`
	UserName string `json:"user-name"`
	UserPass string `json:"user-pass"`
}

type OrgCreator interface {
	Create(org *e.Organization) (*e.Organization, error)
}

type UserProvider interface {
	Create(user *e.User) error
	GetOrganizationIdFromSession(r *http.Request) int
	GetUserIdFromSession(r *http.Request) int
	GetById(userId int) (*e.User, error)
}

type CreateUserRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (s *Service) Create(orgCreator OrgCreator, userProvider UserProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			return
		}

		log := s.Logger
		method := "api.organizations.add"
		response := resp.Service{Writer: &w, Logger: s.Logger, Operation: method}

		var orgData CreateOrganizationRequest
		err := json.NewDecoder(r.Body).Decode(&orgData)

		if err != nil {
			errorText := "При декодировании данных авторизации произошла ошибка"
			log.Errorf(errorText+": %s", err.Error())
			http.Error(w, errorText, http.StatusInternalServerError)
			return
		}

		newOrganization := &e.Organization{
			Name: strings.Trim(orgData.Name, " "),
			Host: strings.Trim(orgData.Host, " "),
		}

		userData := CreateUserRequest{}

		userData.Name = strings.Trim(orgData.UserName, " ")
		userData.Password = strings.Trim(orgData.UserPass, " ")

		if userData.Name == "" || userData.Password == "" {
			response.EmptyData()
			return
		}

		rightsService := rights.New()

		currentUserId := userProvider.GetUserIdFromSession(r)
		if currentUserId == 0 {
			response.Unauthorized()
			return
		}

		currentUserOrganizationId := userProvider.GetOrganizationIdFromSession(r)
		if currentUserOrganizationId == 0 {
			response.Unauthorized()
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userData.Password), s.Config.PasswordCost)

		newUser := e.User{
			UserName:       strings.Trim(userData.Name, " "),
			Password:       string(hashedPassword),
			OrganizationId: currentUserOrganizationId,
			Rights:         user.DefaultUserRights,
			CreateTime:     time.Now().Unix(),
			UpdateTime:     time.Now().Unix(),
		}

		currentUser, err := userProvider.GetById(currentUserId)
		if err != nil {
			response.InternalServerError()
			return
		}

		if rightsService.CheckUserRight(currentUser, rights.AddUser&rights.AddOrganization) {

			createdOrganization, err := orgCreator.Create(newOrganization)
			if err != nil {
				log.Errorf("Error: %s", err.Error())

				response.InternalServerError()
				return
			}

			newUser.OrganizationId = createdOrganization.Id

			err = userProvider.Create(&newUser)
			if err != nil {
				log.Errorf("Error: %s", err.Error())
				response.InternalServerError()
				return
			}
		} else {
			log.Warningf("Method now allowed for user %d", currentUserId)
			response.NotAllowed()
			return
		}

		response.OK()
	}
}
