package migrations

import (
	"database/sql"
	"fmt"
	"github.com/dmidokov/rv2/config"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Init(cfg *config.Configuration) {

	s := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.DbUser, cfg.DbPassword, cfg.DbHost, cfg.DbPort, cfg.DbName)

	db, _ := sql.Open("postgres", s)

	err := db.Ping()

	if err != nil {
		panic("Troubles")
	} else {
		fmt.Println("Ping is OK!")
	}

	err = migrateSQL(db, "postgres", cfg.MigrationPath)

	if err != nil {
		panic(err.Error())
	} else {
		fmt.Println("Migration successfully finished.")
	}
	_ = db.Close()

}

// This function executes the migration scripts.
func migrateSQL(db *sql.DB, driverName string, migrationsPath string) error {
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
