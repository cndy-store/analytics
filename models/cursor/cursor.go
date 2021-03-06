package cursor

import (
	"errors"
	"github.com/cndy-store/analytics/utils/sql"
	"github.com/stellar/go/clients/horizon"
)

var ErrNotFound = errors.New("Cursor not found")
var Current horizon.Cursor

type Cursor struct {
	Id          *uint32 `db:"id",          json:"-"`
	PagingToken *string `db:"paging_token" json:"paging_token,omitempty"`
}

func Update(cursor horizon.Cursor) {
	Current = cursor
}

func Save(db sql.Database) (err error) {
	_, err = db.Exec(`UPDATE cursors SET paging_token=$1 WHERE id=1`, Current)
	return
}

func LoadLatest(db sql.Database) (err error) {
	var c string
	err = db.Get(&c, `SELECT paging_token FROM cursors WHERE id=1`)
	if err != nil {
		return
	}

	Current = horizon.Cursor(c)
	return
}
