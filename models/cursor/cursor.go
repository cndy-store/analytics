package cursor

import (
	"errors"
	"github.com/cndy-store/analytics/utils/sql"
	"github.com/stellar/go/clients/horizon"
)

var ErrNotFound = errors.New("Cursor not found")

type Cursor struct {
	Id          *uint32 `db:"id",          json:"-"`
	PagingToken *string `db:"paging_token" json:"paging_token,omitempty"`
}

func New(db interface{}, cursor string) (err error) {
	_, err = sql.Exec(db, `UPDATE cursors SET paging_token=$1 WHERE id=1`, cursor)
	return
}

func GetLatest(db interface{}) (cursor horizon.Cursor, err error) {
	var c string
	err = sql.Get(db, &c, `SELECT paging_token FROM cursors WHERE id=1`)
	if err != nil {
		return
	}

	cursor = horizon.Cursor(c)
	return
}
