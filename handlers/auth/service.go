package auth

import (
	"github.com/dmidokov/rv2/config"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Service struct {
	Logger      *logrus.Logger
	DB          *pgxpool.Pool
	CookieStore SessionStorage
	Config      *config.Configuration
}

type SessionStorage interface {
	Save(r *http.Request, w http.ResponseWriter, data map[string]interface{}) bool
	Get(r *http.Request, key string) (interface{}, bool)
	SetMaxAge(maxAge int)
}

func New(Logger *logrus.Logger, DB *pgxpool.Pool, CookieStore SessionStorage, Config *config.Configuration) Service {
	return Service{
		Logger:      Logger,
		DB:          DB,
		CookieStore: CookieStore,
		Config:      Config,
	}
}

func (s *Service) Get() {

}
