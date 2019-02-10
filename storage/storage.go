package storage

import (
	"database/sql"
	"github.com/golang-migrate/migrate"
	migratepg "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/lib/pq"
	"log"
)

const (
	DbType           = "postgres"
	ConnectionString = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	MigrationsPath = "file://storage/migrations"
)

// UpdateSchema makes sure the DB schema is up-to-date and all migrations have been applied
func UpdateSchema(db *sql.DB) error {
	driver, err := migratepg.WithInstance(db, &migratepg.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(MigrationsPath, DbType, driver)
	if err != nil {
		return err
	}
	err = m.Up()
	if err == migrate.ErrNoChange {
		log.Println("No changes to the database schema have been made")
		err = nil
	}
	return err
}
