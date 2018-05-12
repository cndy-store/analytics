package cursor

import (
	"github.com/cndy-store/analytics/utils/sql"
	"github.com/stellar/go/clients/horizon"
	"testing"
)

func TestGenesisCursor(t *testing.T) {
	db, err := sql.ResetDB("../../")
	if err != nil {
		t.Error(err)
	}

	genesisCursor := horizon.Cursor("33819440072110101-2") // See db/migrations/0001_initial.up.sql

	currentCursor, err := GetLatest(db)
	if err != nil {
		t.Errorf("cursor.GetLatest(): %s", err)
	}

	if currentCursor != genesisCursor {
		t.Errorf("Latest cursor is %s, expected %s", currentCursor, genesisCursor)
	}
}

func TestNew(t *testing.T) {
	db, err := sql.ResetDB("../../")
	if err != nil {
		t.Error(err)
	}

	newCursor := "33819440072111111-1"

	err = New(db, newCursor)
	if err != nil {
		t.Errorf("cursor.New(): %s", err)
	}

	currentCursor, err := GetLatest(db)
	if err != nil {
		t.Errorf("cursor.GetLastest(): %s", err)
	}

	if currentCursor != horizon.Cursor(newCursor) {
		t.Errorf("Latest cursor is %s, expected %s", currentCursor, newCursor)
	}
}
