package cursor

import (
	"github.com/cndy-store/analytics/utils/cndy"
	"github.com/cndy-store/analytics/utils/sql"
	"github.com/stellar/go/clients/horizon"
	"testing"
)

func TestGenesisCursor(t *testing.T) {
	db, err := sql.OpenAndMigrate("../..")
	if err != nil {
		t.Error(err)
	}

	tx, err := db.Beginx()
	if err != nil {
		t.Error(err)
	}
	defer tx.Rollback()

	genesisCursor := horizon.Cursor(cndy.GenesisCursor)

	currentCursor, err := GetLatest(tx)
	if err != nil {
		t.Errorf("cursor.GetLatest(): %s", err)
	}

	if currentCursor != genesisCursor {
		t.Errorf("Latest cursor is %s, expected %s", currentCursor, genesisCursor)
	}
}

func TestNew(t *testing.T) {
	db, err := sql.OpenAndMigrate("../..")
	if err != nil {
		t.Error(err)
	}

	tx, err := db.Beginx()
	if err != nil {
		t.Error(err)
	}
	defer tx.Rollback()

	newCursor := "33819440072111111-1"

	err = New(tx, newCursor)
	if err != nil {
		t.Errorf("cursor.New(): %s", err)
	}

	currentCursor, err := GetLatest(tx)
	if err != nil {
		t.Errorf("cursor.GetLastest(): %s", err)
	}

	if currentCursor != horizon.Cursor(newCursor) {
		t.Errorf("Latest cursor is %s, expected %s", currentCursor, newCursor)
	}
}
