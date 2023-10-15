package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

func ConnectToDB(dbHost, dbPort, dbUser, dbPassword, dbName string) (*pgxpool.Pool, error) {

	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)
	conn, err := pgxpool.Connect(context.Background(), psqlInfo)
	if err != nil {
		return nil, err
	}

	err = conn.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	return conn, nil

}

func DropSchema(pool *pgxpool.Pool) error {
	_, err := pool.Exec(context.Background(), DropSchemaSQL)
	if err != nil {
		return err
	}

	err = CreateTables(pool)
	if err != nil {
		return err
	}

	return nil
}

func CreateTables(pool *pgxpool.Pool) error {
	_, err := pool.Exec(context.Background(), TablesStructSQL)
	if err != nil {
		return err
	}
	return nil
}

func FillTables(pool *pgxpool.Pool, adminPassword string) error {

	str, err := bcrypt.GenerateFromPassword([]byte(adminPassword), 14)
	logrus.Info(string(str))

	if err != nil {
		return err
	}

	_, err = pool.Exec(context.Background(), TablesDataSQL)
	if err != nil {
		return err
	}
	return nil
}
