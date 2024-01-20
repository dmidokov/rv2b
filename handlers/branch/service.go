package branch

import (
	"github.com/dmidokov/rv2/config"
	e "github.com/dmidokov/rv2/lib/entitie"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Service struct {
	Logger *logrus.Logger
	DB     *pgxpool.Pool
	Config *config.Configuration
}

type OrgProvider interface {
}

type userProvider interface {
	Create(user *e.User) (int, error)
	GetOrganizationIdFromSession(r *http.Request) int
	GetUserIdFromSession(r *http.Request) int
	GetById(userId int) (*e.User, error)
	IsAuthorized(r *http.Request) bool
}

func New(Logger *logrus.Logger, DB *pgxpool.Pool, Config *config.Configuration) Service {
	return Service{
		Logger: Logger,
		DB:     DB,
		Config: Config,
	}
}
