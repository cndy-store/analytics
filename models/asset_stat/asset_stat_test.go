package assetStat

import (
	"github.com/cndy-store/analytics/utils/sql"
	"github.com/cndy-store/analytics/utils/test"
	"testing"
)

func TestGet(t *testing.T) {
	db, err := sql.OpenAndMigrate("../..")
	if err != nil {
		t.Error(err)
	}

	tx, err := db.Beginx()
	if err != nil {
		t.Error(err)
	}
	defer tx.Rollback()

	err = test.InsertTestData(tx)
	if err != nil {
		t.Error(err)
	}

	// Filter{}
	assetStats, err := Get(tx, Filter{})
	if err != nil {
		t.Errorf("assetStat.Get(): %s", err)
	}
	if len(test.Effects) != len(assetStats) {
		t.Errorf("Expected %d assetStats got %d", len(test.Effects), len(assetStats))
	}

	// Filter{From}
	assetStats, err = Get(tx, Filter{From: &test.Effects[2].CreatedAt})
	if err != nil {
		t.Errorf("assetStat.Get(): %s", err)
	}
	if len(test.Effects[2:]) != len(assetStats) {
		t.Errorf("Expected %d assetStats got %d", len(test.Effects[2:]), len(assetStats))
	}

	for i, e := range test.Effects[2:] {
		if e.TotalAmount != *assetStats[i].TotalAmount {
			t.Errorf("Expected %d got: %d", e.TotalAmount, *assetStats[i].TotalAmount)
		}

		if e.NumAccounts != *assetStats[i].NumAccounts {
			t.Errorf("Expected %d got: %d", e.NumAccounts, *assetStats[i].NumAccounts)
		}

		if e.Payments != *assetStats[i].Payments {
			t.Errorf("Expected %d got: %d", e.Payments, *assetStats[i].Payments)
		}
	}

	// Filter{To}
	assetStats, err = Get(tx, Filter{To: &test.Effects[1].CreatedAt})
	if err != nil {
		t.Errorf("assetStat.Get(): %s", err)
	}
	if len(test.Effects[:2]) != len(assetStats) {
		t.Errorf("Expected %d assetStats got %d", len(test.Effects[:2]), len(assetStats))
	}

	for i, e := range test.Effects[:2] {
		if e.TotalAmount != *assetStats[i].TotalAmount {
			t.Errorf("Expected %d got: %d", e.TotalAmount, *assetStats[i].TotalAmount)
		}

		if e.NumAccounts != *assetStats[i].NumAccounts {
			t.Errorf("Expected %d got: %d", e.NumAccounts, *assetStats[i].NumAccounts)
		}

		if e.Payments != *assetStats[i].Payments {
			t.Errorf("Expected %d got: %d", e.Payments, *assetStats[i].Payments)
		}
	}

	// Filter{From, To}
	assetStats, err = Get(tx, Filter{From: &test.Effects[1].CreatedAt, To: &test.Effects[2].CreatedAt})
	if err != nil {
		t.Errorf("assetStat.Get(): %s", err)
	}
	if len(test.Effects[1:3]) != len(assetStats) {
		t.Errorf("Expected %d assetStats got %d", len(test.Effects[1:3]), len(assetStats))
	}

	for i, e := range test.Effects[1:3] {
		if e.TotalAmount != *assetStats[i].TotalAmount {
			t.Errorf("Expected %d got: %d", e.TotalAmount, *assetStats[i].TotalAmount)
		}

		if e.NumAccounts != *assetStats[i].NumAccounts {
			t.Errorf("Expected %d got: %d", e.NumAccounts, *assetStats[i].NumAccounts)
		}

		if e.Payments != *assetStats[i].Payments {
			t.Errorf("Expected %d got: %d", e.Payments, *assetStats[i].Payments)
		}
	}
}

func TestLatest(t *testing.T) {
	db, err := sql.OpenAndMigrate("../..")
	if err != nil {
		t.Error(err)
	}

	tx, err := db.Beginx()
	if err != nil {
		t.Error(err)
	}
	defer tx.Rollback()

	err = test.InsertTestData(tx)
	if err != nil {
		t.Error(err)
	}

	// Filter{}
	assetStats, err := Latest(tx)
	if err != nil {
		t.Errorf("assetStat.Latest(): %s", err)
	}

	// Compare to latest test.Effects[]
	lastEffect := test.Effects[len(test.Effects)-1]
	if lastEffect.TotalAmount != *assetStats.TotalAmount {
		t.Errorf("Expected %d got %d", lastEffect.TotalAmount, *assetStats.TotalAmount)
	}
	if lastEffect.NumAccounts != *assetStats.NumAccounts {
		t.Errorf("Expected %d got %d", lastEffect.NumAccounts, *assetStats.NumAccounts)
	}
	if lastEffect.Payments != *assetStats.Payments {
		t.Errorf("Expected %d got %d", lastEffect.Payments, *assetStats.Payments)
	}
}

func TestConvert(t *testing.T) {
	amount := int64(500000000)
	assetStat := AssetStat{TotalAmount: &amount}
	assetStat.Convert()

	if *assetStat.JsonTotalAmount != "50.0000000" {
		t.Errorf("Expected 50.0000000, got %s", *assetStat.JsonTotalAmount)
	}
}
