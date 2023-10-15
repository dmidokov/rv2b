package handlers

import (
	"encoding/json"
	e "github.com/dmidokov/rv2/entitie"
	"github.com/dmidokov/rv2/response"
	"github.com/dmidokov/rv2/rights"
	"github.com/dmidokov/rv2/users"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type CreateUserRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (hm *Service) GetUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		return
	}

	log := *hm.Logger
	method := "api.users.get"
	response := resp.Service{Writer: &w, Logger: hm.Logger, Operation: method}

	usersService := users.Service{DB: hm.DB, CookieStore: hm.CookieStore, Log: hm.Logger}

	organizationId, err := hm.GetSessionOrganizationId(r)
	if err != nil {
		log.Errorf("Error: %s", err.Error())
		response.InternalServerError()
		return
	}

	navigationItems, err := usersService.GetByOrganizationId(organizationId)
	if err != nil {
		log.Errorf("Error: %s", err.Error())

		response.InternalServerError()
		return
	}

	response.OKWithData(navigationItems)
}

func (hm *Service) DeleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		return
	}

	log := *hm.Logger
	method := "api.user.delete"
	response := resp.Service{Writer: &w, Logger: hm.Logger, Operation: method}

	vars := mux.Vars(r)
	var userId int
	varsId, ok := vars["id"]

	if !ok {
		response.EmptyData()
		return
	}

	userId, err := strconv.Atoi(varsId)
	if err != nil {
		response.InternalServerError()
		return
	}

	usersService := users.New(hm.DB, hm.CookieStore, hm.Logger)

	rightsService := rights.New()

	currentUserId := usersService.GetUserIdFromSession(r)
	if currentUserId == 0 {
		response.Unauthorized()
		return
	}

	currentUser, err := usersService.GetById(currentUserId)
	if err != nil {
		response.InternalServerError()
		return
	}

	if rightsService.CheckUserRight(currentUser.Rights, rights.DeleteUser) {
		err := usersService.Delete(userId)
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

func (hm *Service) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		return
	}

	log := *hm.Logger
	method := "api.user.add"

	response := resp.Service{Writer: &w, Logger: hm.Logger, Operation: method}

	userData := CreateUserRequest{}
	err := json.NewDecoder(r.Body).Decode(&userData)
	if err != nil {
		response.JsonDecodeError()
		return
	}

	userData.Name = strings.Trim(userData.Name, " ")
	userData.Password = strings.Trim(userData.Password, " ")

	if userData.Name == "" || userData.Password == "" {
		response.EmptyData()
		return
	}

	rightsService := rights.New()
	usersService := users.New(hm.DB, hm.CookieStore, hm.Logger)

	currentUserId := usersService.GetUserIdFromSession(r)
	if currentUserId == 0 {
		response.Unauthorized()
		return
	}

	currentUser, err := usersService.GetById(currentUserId)
	if err != nil {
		response.InternalServerError()
		return
	}

	currentUserOrganizationId := usersService.GetOrganizationIdFromSession(r)
	if currentUserOrganizationId == 0 {
		response.Unauthorized()
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userData.Password), hm.Config.PasswordCost)

	newUser := e.User{
		UserName:       strings.Trim(userData.Name, " "),
		Password:       string(hashedPassword),
		OrganizationId: currentUserOrganizationId,
		Rights:         users.DefaultUserRights,
		CreateTime:     time.Now().Unix(),
		UpdateTime:     time.Now().Unix(),
	}

	if rightsService.CheckUserRight(currentUser.Rights, rights.AddUser) {
		err := usersService.Create(&newUser)
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
