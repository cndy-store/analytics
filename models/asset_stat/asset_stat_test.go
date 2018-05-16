package assetStat

import (
	"github.com/cndy-store/analytics/utils/sql"
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	var datasets = []struct {
		PagingToken string
		TotalAmount string
		NumAccounts int32
		NumEffects  int32
		CreatedAt   time.Time
	}{
		{"33819440072110101-0", "1000", 10, 50, time.Date(2018, time.March, 12, 0, 0, 0, 0, time.UTC)},
		{"33819440072110101-1", "1000", 12, 60, time.Date(2018, time.March, 14, 0, 0, 0, 0, time.UTC)},
		{"33819440072110101-2", "2000", 15, 70, time.Date(2018, time.March, 16, 0, 0, 0, 0, time.UTC)},
		{"33819440072110101-3", "2000", 22, 80, time.Date(2018, time.March, 18, 0, 0, 0, 0, time.UTC)},
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
			              VALUES($1, 'CNDY', 'GCJKC2MI63KSQ6MLE6GBSXPDKTDAK43WR522ZYR3F34NPM7Z5UEPIZNX', 'credit_alphanum4', $2, $3, $4, $5)`,
			data.PagingToken, data.TotalAmount, data.NumAccounts, data.NumEffects, data.CreatedAt)
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
