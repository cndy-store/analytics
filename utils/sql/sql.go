package sql

import (
	"database/sql"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var ErrNoRows = sql.ErrNoRows

type Database interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Get(dest interface{}, query string, args ...interface{}) error
	Select(dest interface{}, query string, args ...interface{}) error
	PrepareNamed(query string) (*sqlx.NamedStmt, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

// Open production database and run migrations
func OpenAndMigrate(relPath string) (db *sqlx.DB, err error) {
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

	m, err := migrate.NewWithDatabaseInstance("file://"+relPath+"/db/migrations", "postgres", driver)
	if err != nil {
		return
	}

	err = m.Up()
	if err == migrate.ErrNoChange {
		err = nil
	}

	return
}

// Exec is a type agnostic wrapper for sqlx.Exec() (works with sqlx.DB and sqlx.Tx)
func Exec(db Database, query string, args ...interface{}) (result sql.Result, err error) {
	return db.Exec(query, args...)
}

// Getis a type agnostic wrapper for sqlx.Get() (works with sqlx.DB and sqlx.Tx)
func Get(db Database, obj interface{}, query string, args ...interface{}) (err error) {
	return db.Get(obj, query, args...)
}

// Select is a type agnostic wrapper for sqlx.Select() (works with sqlx.DB and sqlx.Tx)
func Select(db Database, obj interface{}, query string, args ...interface{}) (err error) {
	return db.Select(obj, query, args...)
}

// NamedQuery is a type agnostic wrapper for sqlx.NamedQuery() (works with sqlx.DB and sqlx.Tx)
func NamedQuery(db Database, obj interface{}, query string, arg interface{}) (err error) {
	var stmt *sqlx.NamedStmt
	stmt, err = db.PrepareNamed(query)

	if err != nil {
		return
	}

	// Close the statement (probably this calls DEALLOCATE).
	// This is necessary, as our postgresql connection is long-living, and all prepared
	// statements are kept in memory for the current database session.
	defer stmt.Stmt.Stmt.Close()

	err = stmt.Select(obj, arg)
	return
}

// Exists is a type agnostic function that checks whether a statement returns a row
func Exists(db Database, query string, args ...interface{}) (exists bool, err error) {
	var row *sql.Row

	// Prepare exists query
	query = `SELECT EXISTS(` + query + `) LIMIT 1`
	row = db.QueryRow(query, args...)

	err = row.Scan(&exists)
	return
}
