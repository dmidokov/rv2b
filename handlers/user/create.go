package user

import (
	"encoding/json"
	"github.com/dmidokov/rv2/lib"
	e "github.com/dmidokov/rv2/lib/entitie"
	resp "github.com/dmidokov/rv2/response"
	"github.com/dmidokov/rv2/storage/postgres/rights"
	"github.com/dmidokov/rv2/storage/postgres/user"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
	"time"
)

type userCreator interface {
	GetById(userId int) (*e.User, error)
	GetOrganizationIdFromSession(r *http.Request) int
	GetUserIdFromSession(r *http.Request) int
	GetByOrganizationId(userId int) ([]*e.UserShort, error)
	Create(user *e.User) (int, error)
	SetUserCreateRelations(creatorId int, createdId int) error
}

type CreateUserRequest struct {
	Name             string `json:"name"`
	Password         string `json:"password"`
	EmployeeName     string `json:"employeeName"`
	EmployeeLastName string `json:"employeeLastName"`
	EmployeeLogin    string `json:"employeeLogin"`
	Type             string `json:"type"`
}

func (s *Service) Create(userProvider userCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			return
		}

		log := s.Logger
		method := "api.user.add"

		response := resp.Service{Writer: &w, Logger: s.Logger, Operation: method}

		userData := CreateUserRequest{}
		err := json.NewDecoder(r.Body).Decode(&userData)
		if err != nil {
			response.JsonDecodeError()
			return
		}

		userData.Name = strings.Trim(userData.Name, " ")
		userData.Password = strings.Trim(userData.Password, " ")

		if (userData.Name == "" || userData.Password == "") && userData.Type == "service" {
			response.EmptyData()
			return
		}

		rightsService := rights.New(s.DB, s.Logger)

		currentUserId := userProvider.GetUserIdFromSession(r)
		if currentUserId == 0 {
			response.Unauthorized()
			return
		}

		currentUser, err := userProvider.GetById(currentUserId)
		if err != nil {
			response.InternalServerError()
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
			StartPage:      "/",
		}

		userData.Type = strings.Trim(userData.Type, " ")

		switch userData.Type {
		case "service":
			newUser.Type = 0
		case "employee":
			newUser.Type = 1
			newUser.UserName = userData.EmployeeLogin
		}

		createdUserId := 0

		if rightsService.CheckUserRight(currentUser, lib.AddUser) {
			createdUserId, err = userProvider.Create(&newUser)
			if err != nil {
				log.Errorf("Error: %s", err.Error())
				response.InternalServerError()
				return
			}
			_ = userProvider.SetUserCreateRelations(currentUserId, createdUserId)
		} else {
			log.Warningf("Method now allowed for user %d", currentUserId)
			response.NotAllowed()
			return
		}

		response.OK()
	}
}
