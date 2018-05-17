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

	err = test.InsertAssetStats(tx)
	if err != nil {
		t.Error(err)
	}

	// Filter{}
	assetStats, err := Get(tx, Filter{})
	if err != nil {
		t.Errorf("assetStat.Get(): %s", err)
	}
	if len(test.AssetStats) != len(assetStats) {
		t.Errorf("Expected %d assetStats got %d", len(test.AssetStats), len(assetStats))
	}

	// Filter{From}
	assetStats, err = Get(tx, Filter{From: &test.AssetStats[2].CreatedAt})
	if err != nil {
		t.Errorf("assetStat.Get(): %s", err)
	}
	if len(assetStats) != 2 {
		t.Errorf("Expected 2 assetStats got %d", len(assetStats))
	}
	if test.AssetStats[2].NumAccounts != *assetStats[0].NumAccounts {
		t.Errorf("Expected num_accounts to be %d got: %d", test.AssetStats[2].NumAccounts, *assetStats[0].NumAccounts)
	}
	if test.AssetStats[3].NumAccounts != *assetStats[1].NumAccounts {
		t.Errorf("Expected num_accounts to be %d got: %d", test.AssetStats[3].NumAccounts, *assetStats[1].NumAccounts)
	}

	// Filter{To}
	assetStats, err = Get(tx, Filter{To: &test.AssetStats[1].CreatedAt})
	if err != nil {
		t.Errorf("assetStat.Get(): %s", err)
	}
	if len(assetStats) != 2 {
		t.Errorf("Expected 2 assetStats got %d", len(assetStats))
	}
	if test.AssetStats[0].NumAccounts != *assetStats[0].NumAccounts {
		t.Errorf("Expected num_accounts to be %d got: %d", test.AssetStats[0].NumAccounts, *assetStats[0].NumAccounts)
	}
	if test.AssetStats[1].NumAccounts != *assetStats[1].NumAccounts {
		t.Errorf("Expected num_accounts to be %d got: %d", test.AssetStats[1].NumAccounts, *assetStats[1].NumAccounts)
	}

	// Filter{From, To}
	assetStats, err = Get(tx, Filter{From: &test.AssetStats[1].CreatedAt, To: &test.AssetStats[2].CreatedAt})
	if err != nil {
		t.Errorf("assetStat.Get(): %s", err)
	}
	if len(assetStats) != 2 {
		t.Errorf("Expected 2 assetStats got %d", len(assetStats))
	}
	if test.AssetStats[1].NumAccounts != *assetStats[0].NumAccounts {
		t.Errorf("Expected num_accounts to be %d got: %d", test.AssetStats[1].NumAccounts, *assetStats[0].NumAccounts)
	}
	if test.AssetStats[2].NumAccounts != *assetStats[1].NumAccounts {
		t.Errorf("Expected num_accounts to be %d got: %d", test.AssetStats[2].NumAccounts, *assetStats[1].NumAccounts)
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
