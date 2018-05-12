package sql

import (
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Open production database and run migrations
func OpenAndMigrate() (db *sqlx.DB, err error) {
	/* Open connection to postgresql to the server specified in the following environment variables:
	 *
	 * PGHOST=127.0.0.1
	 * PGUSER=cndy
	 * PGDATABASE=cndy
	 * PGSSLMODE=disable
	 */
	db, err = sqlx.Open("postgres", "postgres://")
	if err != nil {
		return
	}

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return
	}

	m, err := migrate.NewWithDatabaseInstance("file://db/migrations", "postgres", driver)
	if err != nil {
		return
	}

	err = m.Up()
	if err == migrate.ErrNoChange {
		err = nil
	}

	return
}

// Delete and re-create test database, run migrations and insert test data
func ResetDB(relPath string) (db *sqlx.DB, err error) {
	db, err = sqlx.Open("postgres", "postgres://")
	if err != nil {
		return
	}

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return
	}

	m, err := migrate.NewWithDatabaseInstance("file://"+relPath+"/db/migrations", "postgres", driver)
	if err != nil {
		return
	}

	err = m.Drop()
	if err != nil {
		return
	}

	err = m.Up()
	if err == migrate.ErrNoChange {
		err = nil
	}

	return
}
