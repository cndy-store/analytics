package sqlite

import (
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/sqlite3"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// Open database and run migrations
func OpenAndMigrate() (db *sqlx.DB, err error) {
	db, err = sqlx.Open("sqlite3", "database.sqlite3")
	if err != nil {
		return
	}

	driver, err := sqlite3.WithInstance(db.DB, &sqlite3.Config{})
	if err != nil {
		return
	}

	m, err := migrate.NewWithDatabaseInstance("file://db/migrations", "sqlite3", driver)
	if err != nil {
		return
	}

	err = m.Up()
	if err == migrate.ErrNoChange {
		err = nil
	}

	return
}
