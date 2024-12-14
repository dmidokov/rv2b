package rights

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

type Service struct {
	DB  *pgxpool.Pool
	Log *logrus.Logger
}

func New(DB *pgxpool.Pool, Log *logrus.Logger) *Service {
	return &Service{
		DB:  DB,
		Log: Log,
	}
}
