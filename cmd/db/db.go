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

func DropSchema(pool *pgxpool.Pool, dbName string) error {

	var DropSchemaSQL = `DROP TABLE IF EXISTS users;
						
						DROP TABLE IF EXISTS organizations;
						
						DROP TABLE IF EXISTS navigation;
						
						DROP TABLE IF EXISTS right_category_ids;
						
						DROP TABLE IF EXISTS rights;
						
						DROP TABLE IF EXISTS rights_names;
						
						DROP TABLE IF EXISTS branches;
						
						DROP TABLE IF EXISTS user_branches;
						
						DROP TABLE IF EXISTS users_create_relations;
						
						DROP TABLE IF EXISTS entity_group_to_entity_name;
						
						DROP TABLE IF EXISTS hot_switch_relations;
`

	_, err := pool.Exec(context.Background(), DropSchemaSQL)
	if err != nil {
		return err
	}
	return nil
}
