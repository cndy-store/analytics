package sql

import (
	"context"
	"database/sql"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var ErrNoRows = sql.ErrNoRows

// This type serves as an abstraction for sqlx.DB and sqlx.Tx and supports all functions both
// of the types have in common. This allows as to use their functions type agnostically.
type Database interface {
	// Common functions of sql.DB and sql.Tx
	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row

	// Common functions of sqlx.DB and sqlx.Tx
	BindNamed(query string, arg interface{}) (string, []interface{}, error)
	DriverName() string
	Get(dest interface{}, query string, args ...interface{}) error
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	MustExec(query string, args ...interface{}) sql.Result
	MustExecContext(ctx context.Context, query string, args ...interface{}) sql.Result
	NamedExec(query string, arg interface{}) (sql.Result, error)
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
	NamedQuery(query string, arg interface{}) (*sqlx.Rows, error)
	PrepareNamed(query string) (*sqlx.NamedStmt, error)
	PrepareNamedContext(ctx context.Context, query string) (*sqlx.NamedStmt, error)
	Preparex(query string) (*sqlx.Stmt, error)
	PreparexContext(ctx context.Context, query string) (*sqlx.Stmt, error)
	QueryRowx(query string, args ...interface{}) *sqlx.Row
	QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row
	Queryx(query string, args ...interface{}) (*sqlx.Rows, error)
	QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error)
	Rebind(query string) string
	Select(dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
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
