package organization

import (
	"github.com/dmidokov/rv2/config"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

type Service struct {
	Logger *logrus.Logger
	DB     *pgxpool.Pool
	Config *config.Configuration
}

type OrgProvider interface {
}

func New(Logger *logrus.Logger, DB *pgxpool.Pool, Config *config.Configuration) Service {
	return Service{
		Logger: Logger,
		DB:     DB,
		Config: Config,
	}
}
