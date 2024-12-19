package group

import (
	"github.com/dmidokov/rv2/config"
	"github.com/dmidokov/rv2/lib"
	resp "github.com/dmidokov/rv2/response"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Service struct {
	Logger *logrus.Logger
	DB     *pgxpool.Pool
	Config *config.Configuration
}

func New(logger *logrus.Logger, db *pgxpool.Pool, cfg *config.Configuration) Service {
	return Service{
		Logger: logger,
		DB:     db,
		Config: cfg,
	}
}

func (s *Service) GetGroups(userProvider userProvider, rightsProvider rightsProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			return
		}

		method := "api.getGroups"
		response := resp.Service{Writer: &w, Logger: s.Logger, Operation: method}

		currentUserId := userProvider.GetUserIdFromSession(r)
		if currentUserId == 0 {
			response.Unauthorized()
			return
		}

		currentUser, err := userProvider.GetById(currentUserId)
		if err != nil {
			s.Logger.Warning("User not authorized")
			response.InternalServerError()
			return
		}

		if !rightsProvider.CheckUserRight(currentUser, lib.ViewUserGroups) {
			s.Logger.Warning("User have no rights to view group list")
			response.NotAllowed()
			return
		}

		groups, err := rightsProvider.GetGroupsByOrganizationId(currentUser.OrganizationId)
		if err != nil {
			s.Logger.Warningf("can't get groups list: %s", err.Error())
		}

		response.OKWithData(groups)

	}
}

func (s *Service) GetAvailableRights(userProvider userProvider, rightsProvider rightsProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			return
		}

		method := "api.getAvailableRights"
		response := resp.Service{Writer: &w, Logger: s.Logger, Operation: method}

		currentUserId := userProvider.GetUserIdFromSession(r)
		if currentUserId == 0 {
			response.Unauthorized()
			return
		}

		currentUser, err := userProvider.GetById(currentUserId)
		if err != nil {
			s.Logger.Warning("Can't get user by id: ", err.Error())
			response.InternalServerError()
			return
		}

		currentUserGroupsRights, err := rightsProvider.GetUserGroupsRights(currentUserId)
		if err != nil {
			s.Logger.Warning("Can't get user group rights: ", err.Error())
			response.InternalServerError()
			return
		}

		currentUserRightsWithDescriptions, err := rightsProvider.GetByUserRights(currentUser.Rights | currentUserGroupsRights)
		if err != nil {
			s.Logger.Warning("can't get user rights descrioptions: ", err.Error())
			response.InternalServerError()
			return
		}

		response.OKWithData(currentUserRightsWithDescriptions)
	}
}
