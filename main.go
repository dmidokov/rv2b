package main

import (
	"context"
	"github.com/dmidokov/rv2/config"
	"github.com/dmidokov/rv2/db"
	"github.com/dmidokov/rv2/handlers"
	"github.com/dmidokov/rv2/migrations"
	"github.com/dmidokov/rv2/session/cookie"
	"github.com/dmidokov/rv2/sse"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"os"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {

	cfg := config.LoadConfig()

	logger := setupLog(cfg)

	eventService := setupSSE()
	eventService.Run()

	conn := connectDB(cfg)

	if cfg.DeleteTablesBeforeStart == 1 {
		logger.Info("Удаление схемы")
		err := db.DropSchema(conn)
		if err != nil {
			logger.Fatal("Can't drop schema")
		}
	}

	migrations.Init(cfg)

	store := cookie.New(cfg.SessionsSecret)

	handler := handlers.New(conn, cfg, store, logger, eventService)

	router, err := handler.Router()

	if err != nil {
		logger.Fatalf("Регистрация завершилась с ошибкой: %s", err)
	}

	go func() {
		logger.Fatal(http.ListenAndServeTLS(":"+cfg.SSLPort, cfg.SecretsPath+"server.crt", cfg.SecretsPath+"server.key", router))
	}()

	go func() {
		logger.Fatal(http.ListenAndServe(":"+cfg.HttpPort, http.HandlerFunc(handlers.Redirect)))
	}()

	finish := make(chan bool)
	<-finish
}

func setupLog(cfg *config.Configuration) *logrus.Logger {

	logger := &logrus.Logger{}

	switch cfg.MODE {
	case config.DEV:
		logger = &logrus.Logger{
			Out: os.Stdout,
			Formatter: &logrus.TextFormatter{
				DisableColors: false,
			},
			Hooks:        make(logrus.LevelHooks),
			Level:        logrus.DebugLevel,
			ReportCaller: true,
		}
	}

	return logger

}

func connectDB(cfg *config.Configuration) *pgxpool.Pool {
	conn, err := db.ConnectToDB(
		cfg.DbHost,
		cfg.DbPort,
		cfg.DbUser,
		cfg.DbPassword,
		cfg.DbName)
	if err != nil {
		log.Fatalf("Подключение завершилось с ошибкой : %s", err)
	}

	err = conn.Ping(context.Background())
	if err != nil {
		log.Fatalf("Ping завершился с ошибкой : %s", err)
	}

	return conn
}

func setupSSE() *sse.EventService {
	return &sse.EventService{
		Chanel:    make(chan sse.Event, 10),
		Receivers: make(map[sse.EventName]map[int]sse.Receiver, 1000),
	}
}
