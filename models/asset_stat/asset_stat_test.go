package assetStat

import (
	"github.com/cndy-store/analytics/utils/cndy"
	"github.com/cndy-store/analytics/utils/sql"
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	var datasets = []struct {
		PagingToken string
		TotalAmount int64
		NumAccounts int32
		NumEffects  int32
		CreatedAt   time.Time
	}{
		{"33819440072110101-0", 10000000000, 10, 50, time.Date(2018, time.March, 12, 0, 0, 0, 0, time.UTC)},
		{"33819440072110101-1", 10000000000, 12, 60, time.Date(2018, time.March, 14, 0, 0, 0, 0, time.UTC)},
		{"33819440072110101-2", 20000000000, 15, 70, time.Date(2018, time.March, 16, 0, 0, 0, 0, time.UTC)},
		{"33819440072110101-3", 20000000000, 22, 80, time.Date(2018, time.March, 18, 0, 0, 0, 0, time.UTC)},
	}

	db, err := sql.OpenAndMigrate("../..")
	if err != nil {
		t.Error(err)
	}

	tx, err := db.Beginx()
	if err != nil {
		t.Error(err)
	}
	defer tx.Rollback()

	for _, data := range datasets {
		_, err = tx.Exec(`INSERT INTO asset_stats(paging_token, asset_code, asset_issuer, asset_type, total_amount, num_accounts, num_effects, created_at)
			              VALUES($1, $2, $3, 'credit_alphanum4', $4, $5, $6, $7)`,
			data.PagingToken, cndy.AssetCode, cndy.AssetIssuer, data.TotalAmount, data.NumAccounts, data.NumEffects, data.CreatedAt)
		if err != nil {
			t.Error(err)
		}
	}

	// Filter{}
	assetStats, err := Get(tx, Filter{})
	if err != nil {
		t.Errorf("assetStat.Get(): %s", err)
	}
	if len(datasets) != len(assetStats) {
		t.Errorf("Expected %d assetStats got %d", len(datasets), len(assetStats))
	}

	// Filter{From}
	assetStats, err = Get(tx, Filter{From: &datasets[2].CreatedAt})
	if err != nil {
		t.Errorf("assetStat.Get(): %s", err)
	}
	if len(assetStats) != 2 {
		t.Errorf("Expected 2 assetStats got %d", len(assetStats))
	}
	if datasets[2].NumAccounts != *assetStats[0].NumAccounts {
		t.Errorf("Expected num_accounts to be %d got: %d", datasets[2].NumAccounts, *assetStats[0].NumAccounts)
	}
	if datasets[3].NumAccounts != *assetStats[1].NumAccounts {
		t.Errorf("Expected num_accounts to be %d got: %d", datasets[3].NumAccounts, *assetStats[1].NumAccounts)
	}

	// Filter{To}
	assetStats, err = Get(tx, Filter{To: &datasets[1].CreatedAt})
	if err != nil {
		t.Errorf("assetStat.Get(): %s", err)
	}
	if len(assetStats) != 2 {
		t.Errorf("Expected 2 assetStats got %d", len(assetStats))
	}
	if datasets[0].NumAccounts != *assetStats[0].NumAccounts {
		t.Errorf("Expected num_accounts to be %d got: %d", datasets[0].NumAccounts, *assetStats[0].NumAccounts)
	}
	if datasets[1].NumAccounts != *assetStats[1].NumAccounts {
		t.Errorf("Expected num_accounts to be %d got: %d", datasets[1].NumAccounts, *assetStats[1].NumAccounts)
	}

	// Filter{From, To}
	assetStats, err = Get(tx, Filter{From: &datasets[1].CreatedAt, To: &datasets[2].CreatedAt})
	if err != nil {
		t.Errorf("assetStat.Get(): %s", err)
	}
	if len(assetStats) != 2 {
		t.Errorf("Expected 2 assetStats got %d", len(assetStats))
	}
	if datasets[1].NumAccounts != *assetStats[0].NumAccounts {
		t.Errorf("Expected num_accounts to be %d got: %d", datasets[1].NumAccounts, *assetStats[0].NumAccounts)
	}
	if datasets[2].NumAccounts != *assetStats[1].NumAccounts {
		t.Errorf("Expected num_accounts to be %d got: %d", datasets[2].NumAccounts, *assetStats[1].NumAccounts)
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
