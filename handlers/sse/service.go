package sse

import (
	"github.com/dmidokov/rv2/config"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
	"net/http"
)

type SessionStorage interface {
	Save(r *http.Request, w http.ResponseWriter, data map[string]interface{}) bool
	GetByKey(r *http.Request, key string) (interface{}, bool)
	SetMaxAge(maxAge int)
}

type Service struct {
	Logger      *logrus.Logger
	DB          *pgxpool.Pool
	CookieStore SessionStorage
	Config      *config.Configuration
	SSE         *EventService
}

func New(Logger *logrus.Logger, DB *pgxpool.Pool, CookieStore SessionStorage, Config *config.Configuration, sse *EventService) Service {
	return Service{
		Logger:      Logger,
		DB:          DB,
		CookieStore: CookieStore,
		Config:      Config,
		SSE:         sse,
	}
}
