package migrations

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const migrationsPath = "file://migrations/"

func Init() {

	db, _ := sql.Open("postgres", "postgres://remonttiuser:deltad2dt@127.0.0.1:5432/remonttiv2db?sslmode=disable")

	err := db.Ping()

	if err != nil {
		panic("Troubles")
	} else {
		fmt.Println("Ping is OK!")
	}

	err = migrateSQL(db, "postgres")

	if err != nil {
		panic(err)
	} else {
		fmt.Println("Migration successfully finished.")
	}
	db.Close()

}

// This function executes the migration scripts.
func migrateSQL(db *sql.DB, driverName string) error {
	driver, _ := postgres.WithInstance(db, &postgres.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		migrationsPath,
		driverName,
		driver,
	)
	if err != nil {
		return err
	}

	println("Migration start")

	if err := m.Down(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}
