package sql

import (
	"database/sql"
	"errors"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var ErrNoRows = sql.ErrNoRows

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
func Exec(db interface{}, query string, args ...interface{}) (result sql.Result, err error) {
	switch db.(type) {
	case *sqlx.DB:
		result, err = db.(*sqlx.DB).Exec(query, args...)
	case *sqlx.Tx:
		result, err = db.(*sqlx.Tx).Exec(query, args...)
	default:
		err = errors.New("Unknown DB interface{} in sql.Exec()")
	}
	return
}

// Getis a type agnostic wrapper for sqlx.Get() (works with sqlx.DB and sqlx.Tx)
func Get(db, obj interface{}, query string, args ...interface{}) (err error) {
	switch db.(type) {
	case *sqlx.DB:
		err = db.(*sqlx.DB).Get(obj, query, args...)
	case *sqlx.Tx:
		err = db.(*sqlx.Tx).Get(obj, query, args...)
	default:
		err = errors.New("Unknown DB interface{} in sql.Get()")
	}
	return
}

// Select is a type agnostic wrapper for sqlx.Select() (works with sqlx.DB and sqlx.Tx)
func Select(db interface{}, obj interface{}, query string, args ...interface{}) (err error) {
	switch db.(type) {
	case *sqlx.DB:
		err = db.(*sqlx.DB).Select(obj, query, args...)
	case *sqlx.Tx:
		err = db.(*sqlx.Tx).Select(obj, query, args...)
	default:
		err = errors.New("Unknown DB interface{} in sql.Select()")
	}
	return
}

// NamedQuery is a type agnostic wrapper for sqlx.NamedQuery() (works with sqlx.DB and sqlx.Tx)
func NamedQuery(db interface{}, obj interface{}, query string, arg interface{}) (err error) {
	var stmt *sqlx.NamedStmt

	switch db.(type) {
	case *sqlx.DB:
		stmt, err = db.(*sqlx.DB).PrepareNamed(query)
	case *sqlx.Tx:
		stmt, err = db.(*sqlx.Tx).PrepareNamed(query)
	default:
		err = errors.New("Unknown DB interface{} in sql.NamedQuery()")
	}

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
func Exists(db interface{}, query string, args ...interface{}) (exists bool, err error) {
	var row *sql.Row

	// Prepare exists query
	query = `SELECT EXISTS(` + query + `) LIMIT 1`

	switch db.(type) {
	case *sqlx.DB:
		row = db.(*sqlx.DB).QueryRow(query, args...)
	case *sqlx.Tx:
		row = db.(*sqlx.Tx).QueryRow(query, args...)
	default:
		err = errors.New("Unknown DB interface{} in sql.Exists()")
	}

	err = row.Scan(&exists)
	return
}
