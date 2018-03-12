package cursor

import (
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/stellar/go/clients/horizon"
	"log"
)

// Genesis cursor when the asset was introduced
const GENESIS_CURSOR = "33170762571452437-1"

var ErrNotFound = errors.New("Cursor not found")

type Cursor struct {
	Id          *uint32 `db:"id",          json:"-"`
	PagingToken *string `db:"paging_token" json:"paging_token,omitempty"`
}

func New(db *sqlx.DB, cursor string) (err error) {
	_, err = db.Exec(`UPDATE cursors SET paging_token=$1 WHERE id=1`, cursor)
	return
}

func GetLatest(db *sqlx.DB) (cursor horizon.Cursor) {
	var c string
	err := db.Get(&c, `SELECT paging_token FROM cursors WHERE id=1`)
	if err != nil {
		log.Printf("Couldn't get cursor: %s", err)
		log.Printf("Using genesis cursor %s", GENESIS_CURSOR)
		return horizon.Cursor(GENESIS_CURSOR)
	}

	return horizon.Cursor(c)
}
