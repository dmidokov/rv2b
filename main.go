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

	log := setupLog(cfg)

	sse := setupSSE()
	sse.Run()

	conn := connectDB(cfg)

	if cfg.DeleteTablesBeforeStart == 1 {
		db.DropSchema(conn)
	}

	migrations.Init()

	store := cookie.New(cfg.SessionsSecret)

	handler := handlers.New(conn, cfg, store, log, sse)

	router, err := handler.Router()

	if err != nil {
		log.Fatalf("Регистрация завершилась с ошибкой: %s", err)
	}

	go func() {
		log.Fatal(http.ListenAndServeTLS(":443", "secrets/server.crt", "secrets/server.key", router))
	}()

	go func() {
		log.Fatal(http.ListenAndServe(":80", http.HandlerFunc(handlers.Redirect)))
	}()

	finish := make(chan bool)
	<-finish
}

func setupLog(cfg *config.Configuration) *logrus.Logger {

	log := &logrus.Logger{}

	switch cfg.MODE {
	case config.DEV:
		log = &logrus.Logger{
			Out: os.Stdout,
			Formatter: &logrus.TextFormatter{
				DisableColors: false,
			},
			Hooks:        make(logrus.LevelHooks),
			Level:        logrus.DebugLevel,
			ReportCaller: true,
		}
	}

	return log

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
