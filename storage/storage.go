package storage

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate"
	migratepg "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/lib/pq"
	"log"
	"path"
	"runtime"
)

var (
	_, file, _, _ = runtime.Caller(0)
	basepath      = path.Dir(file)
)

// Config holds all storage configuration settings
type Config struct {
	Type string
	URI  string
}

// New creates an *sql.DB object and updates the database with the latest migrations
func New(cfg Config) (*sql.DB, error) {
	db, err := sql.Open(cfg.Type, cfg.URI)
	if err != nil {
		return nil, err
	}

	err = updateSchema(db, cfg.Type)
	if err != nil {
		return nil, err
	}
	log.Println("Database is up-to-date")

	return db, nil
}

func updateSchema(db *sql.DB, dbType string) error {
	driver, err := migratepg.WithInstance(db, &migratepg.Config{})
	if err != nil {
		return err
	}

	migrationsURL := fmt.Sprintf("file://%s/migrations", basepath)
	m, err := migrate.NewWithDatabaseInstance(migrationsURL, dbType, driver)
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
