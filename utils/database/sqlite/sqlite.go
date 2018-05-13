package sqlite

import (
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/sqlite3"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"os"
)

// Open production database and run migrations
func OpenAndMigrate() (db *sqlx.DB, err error) {
	db, err = sqlx.Open("sqlite3", "database.sqlite3")
	if err != nil {
		return
	}

	err = runMigrations(db, ".")
	if err != nil {
		return
	}

	return
}

// Delete and re-create test database, run migrations and insert test data
func InitTestDB(relPath string) (db *sqlx.DB, err error) {
	testDBPath := relPath + "/test.sqlite3"

	// Remove test database if it exists
	_, err = os.Stat(testDBPath)
	if err == nil {
		err2 := os.Remove(testDBPath)
		if err2 != nil {
			err = err2
			return
		}
	}

	db, err = sqlx.Open("sqlite3", testDBPath)
	if err != nil {
		return
	}

	err = runMigrations(db, relPath)
	if err != nil {
		return
	}

	// err = addTestData(db, relPath)
	// if err != nil {
	// 	return
	// }

	return
}

func runMigrations(db *sqlx.DB, relPath string) (err error) {
	driver, err := sqlite3.WithInstance(db.DB, &sqlite3.Config{})
	if err != nil {
		return
	}

	m, err := migrate.NewWithDatabaseInstance("file://"+relPath+"/db/migrations", "sqlite3", driver)
	if err != nil {
		return
	}

	err = m.Up()
	if err == migrate.ErrNoChange {
		err = nil
	}

	return
}
