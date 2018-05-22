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

	err = LoadLatest(tx)
	if err != nil {
		t.Errorf("cursor.LoadLatest(): %s", err)
	}

	if Current != genesisCursor {
		t.Errorf("Latest cursor is %s, expected %s", Current, genesisCursor)
	}
}

func TestSave(t *testing.T) {
	db, err := sql.OpenAndMigrate("../..")
	if err != nil {
		t.Error(err)
	}

	tx, err := db.Beginx()
	if err != nil {
		t.Error(err)
	}
	defer tx.Rollback()

	newCursor := horizon.Cursor("33819440072111111-1")
	Update(newCursor)

	err = Save(tx)
	if err != nil {
		t.Errorf("cursor.Save(): %s", err)
	}

	err = LoadLatest(tx)
	if err != nil {
		t.Errorf("cursor.GetLastest(): %s", err)
	}

	if Current != horizon.Cursor(newCursor) {
		t.Errorf("Latest cursor is %s, expected %s", Current, newCursor)
	}
}
