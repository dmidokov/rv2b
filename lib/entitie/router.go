package entitie

import (
	"context"
	"github.com/dmidokov/rv2/config"
	"github.com/dmidokov/rv2/handlers/sse"
	"github.com/dmidokov/rv2/session/cookie"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

type App struct {
	Ctx         context.Context
	DB          *pgxpool.Pool
	Config      *config.Configuration
	CookieStore *cookie.Service
	Logger      *logrus.Logger
	SSE         *sse.EventService
}
