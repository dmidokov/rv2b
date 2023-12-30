package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
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
	return nil
}
