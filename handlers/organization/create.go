package organization

import (
	"encoding/json"
	"github.com/dmidokov/rv2/lib"
	"github.com/dmidokov/rv2/lib/entitie"
	resp "github.com/dmidokov/rv2/response"
	"github.com/dmidokov/rv2/storage/postgres/rights"
	"github.com/dmidokov/rv2/storage/postgres/user"
	"github.com/jackc/pgx/v4"
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
	Create(org *entitie.Organization) (*entitie.Organization, error)
	GetByHostName(hostName string) (*entitie.Organization, error)
}

type UserProvider interface {
	Create(user *entitie.User) (int, error)
	GetOrganizationIdFromSession(r *http.Request) int
	GetUserIdFromSession(r *http.Request) int
	GetById(userId int) (*entitie.User, error)
}

type CreateUserRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type CreateOKResponse struct {
	resp.Response
	Data entitie.Organization `json:"data"`
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
			response.JsonDecodeError()
			return
		}

		newOrganization := &entitie.Organization{
			Name: strings.Trim(orgData.Name, " "),
			Host: strings.Trim(orgData.Host, " "),
		}

		_, err = orgCreator.GetByHostName(newOrganization.Host)
		if err != nil {
			if err.Error() == pgx.ErrNoRows.Error() {
				log.Info("Creating organization is not created before")
			} else {
				log.Errorf("Duplicate")
				response.InternalServerError()
				return
			}
		}

		userData := CreateUserRequest{}

		userData.Name = strings.Trim(orgData.UserName, " ")
		userData.Password = strings.Trim(orgData.UserPass, " ")

		if userData.Name == "" || userData.Password == "" {
			response.EmptyData()
			return
		}

		rightsService := rights.New(s.DB, s.Logger)

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

		newUser := entitie.User{
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

		newOrganization.Creator = currentUserId

		if rightsService.CheckUserRight(currentUser, lib.AddUser&lib.AddOrganization) {

			createdOrganization, err := orgCreator.Create(newOrganization)
			if err != nil {
				log.Errorf("Error: %s", err.Error())

				response.InternalServerError()
				return
			}

			newUser.OrganizationId = createdOrganization.Id

			_, err = userProvider.Create(&newUser)
			if err != nil {
				log.Errorf("Error: %s", err.Error())
				response.InternalServerError()
				return
			}

			_ = json.NewEncoder(w).Encode(CreateOKResponse{
				Response: resp.Response{
					Status: resp.StatusOK,
				},
				Data: *createdOrganization,
			})
		} else {
			log.Warningf("Method now allowed for user %d", currentUserId)
			response.NotAllowed()
			return
		}

	}
}
