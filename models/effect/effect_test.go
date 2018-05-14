package effect

import (
	"github.com/cndy-store/analytics/utils/sql"
	"testing"
	"time"
)

var datasets = []struct {
	PagingToken string
	Account     string
	Amount      string
	Type        string
	CreatedAt   time.Time
}{
	{"34028708058632193-0", "GDNH64DRUT4CY3UJLWQIB655PQ6OG34UGYB4NC5DC4TYWLNJIBCEYTTD", "", "trustline_created", time.Date(2018, time.March, 8, 0, 0, 0, 0, time.UTC)},
	{"34028708058632193-1", "GDNH64DRUT4CY3UJLWQIB655PQ6OG34UGYB4NC5DC4TYWLNJIBCEYTTD", "1000.0000000", "account_credited", time.Date(2018, time.March, 10, 0, 0, 0, 0, time.UTC)},
	{"34028708058632193-2", "GCJKC2MI63KSQ6MLE6GBSXPDKTDAK43WR522ZYR3F34NPM7Z5UEPIZNX", "1000.0000000", "account_debited", time.Date(2018, time.March, 12, 0, 0, 0, 0, time.UTC)},
	{"34028708058632193-3", "GBEYRLI7OCZU7JVT33GBVVI5XWVCQNSUI3TXDE7Z5MWC6CLQMSTPDT6A", "", "trustline_created", time.Date(2018, time.March, 14, 0, 0, 0, 0, time.UTC)},
	{"34028708058632193-4", "GBEYRLI7OCZU7JVT33GBVVI5XWVCQNSUI3TXDE7Z5MWC6CLQMSTPDT6A", "15.0000000", "account_credited", time.Date(2018, time.March, 16, 0, 0, 0, 0, time.UTC)},
	{"34028708058632193-5", "GDNH64DRUT4CY3UJLWQIB655PQ6OG34UGYB4NC5DC4TYWLNJIBCEYTTD", "15.0000000", "account_debited", time.Date(2018, time.March, 18, 0, 0, 0, 0, time.UTC)},
	{"34028708058632193-6", "GDNH64DRUT4CY3UJLWQIB655PQ6OG34UGYB4NC5DC4TYWLNJIBCEYTTD", "100.0000000", "account_credited", time.Date(2018, time.March, 20, 0, 0, 0, 0, time.UTC)},
	{"34028708058632193-7", "GCJKC2MI63KSQ6MLE6GBSXPDKTDAK43WR522ZYR3F34NPM7Z5UEPIZNX", "100.0000000", "account_debited", time.Date(2018, time.March, 22, 0, 0, 0, 0, time.UTC)},
}

// TODO: Check effect.New() and verify that assetStat.New() was called correctly

func TestGet(t *testing.T) {
	db, err := sql.OpenAndMigrate("../../")
	if err != nil {
		t.Error(err)
	}

	tx, err := db.Beginx()
	if err != nil {
		t.Error(err)
	}
	defer tx.Rollback()

	for _, data := range datasets {
		_, err := sql.Exec(tx, `INSERT INTO effects(effect_id, operation, paging_token, account, amount, type, asset_type, asset_issuer, asset_code, created_at)
			                    VALUES('0034028708058632193-0000000001', 'https://horizon-testnet.stellar.org/operations/34028708058632193', $1, $2, $3, $4, 'credit_alphanum4', 'GCJKC2MI63KSQ6MLE6GBSXPDKTDAK43WR522ZYR3F34NPM7Z5UEPIZNX', 'CNDY', $5)`,
			data.PagingToken, data.Account, data.Amount, data.Type, data.CreatedAt)
		if err != nil {
			return
		}
	}

	// Filter{}
	effects, err := Get(tx, Filter{})
	if err != nil {
		t.Errorf("effect.Get(): %s", err)
	}
	if len(datasets) != len(effects) {
		t.Errorf("Expected %d effects got %d", len(datasets), len(effects))
	}

	// Filter{From}
	effects, err = Get(tx, Filter{From: &datasets[5].CreatedAt})
	if err != nil {
		t.Errorf("effect.Get(): %s", err)
	}
	if len(effects) != 3 {
		t.Errorf("Expected 3 effects got %d", len(effects))
	}
	if datasets[5].PagingToken != *effects[0].PagingToken {
		t.Errorf("Expected paging_token to be %s got: %s", datasets[3].PagingToken, *effects[0].PagingToken)
	}
	if datasets[6].PagingToken != *effects[1].PagingToken {
		t.Errorf("Expected paging_token to be %s got: %s", datasets[4].PagingToken, *effects[1].PagingToken)
	}
	if datasets[7].PagingToken != *effects[2].PagingToken {
		t.Errorf("Expected paging_token to be %s got: %s", datasets[5].PagingToken, *effects[2].PagingToken)
	}

	// Filter{To}
	effects, err = Get(tx, Filter{To: &datasets[2].CreatedAt})
	if err != nil {
		t.Errorf("effect.Get(): %s", err)
	}
	if len(effects) != 3 {
		t.Errorf("Expected 3 effects got %d", len(effects))
	}
	if datasets[0].PagingToken != *effects[0].PagingToken {
		t.Errorf("Expected paging_token to be %s got: %s", datasets[0].PagingToken, *effects[0].PagingToken)
	}
	if datasets[1].PagingToken != *effects[1].PagingToken {
		t.Errorf("Expected paging_token to be %s got: %s", datasets[1].PagingToken, *effects[1].PagingToken)
	}
	if datasets[2].PagingToken != *effects[2].PagingToken {
		t.Errorf("Expected paging_token to be %s got: %s", datasets[2].PagingToken, *effects[2].PagingToken)
	}

	// Filter{From, To}
	effects, err = Get(tx, Filter{From: &datasets[3].CreatedAt, To: &datasets[4].CreatedAt})
	if err != nil {
		t.Errorf("effect.Get(): %s", err)
	}
	if len(effects) != 2 {
		t.Errorf("Expected 2 effects got %d", len(effects))
	}
	if datasets[3].PagingToken != *effects[0].PagingToken {
		t.Errorf("Expected paging_token to be %s got: %s", datasets[3].PagingToken, *effects[0].PagingToken)
	}
	if datasets[4].PagingToken != *effects[1].PagingToken {
		t.Errorf("Expected paging_token to be %s got: %s", datasets[4].PagingToken, *effects[1].PagingToken)
	}
}

func TestAccountCount(t *testing.T) {
	db, err := sql.OpenAndMigrate("../../")
	if err != nil {
		t.Error(err)
	}

	tx, err := db.Beginx()
	if err != nil {
		t.Error(err)
	}
	defer tx.Rollback()

	for _, data := range datasets {
		_, err := sql.Exec(tx, `INSERT INTO effects(effect_id, operation, paging_token, account, amount, type, asset_type, asset_issuer, asset_code, created_at)
			                    VALUES('0034028708058632193-0000000001', 'https://horizon-testnet.stellar.org/operations/34028708058632193', '34028708058632193-1', $1, $2, $3, 'credit_alphanum4', 'GCJKC2MI63KSQ6MLE6GBSXPDKTDAK43WR522ZYR3F34NPM7Z5UEPIZNX', 'CNDY', $4)`,
			data.Account, data.Amount, data.Type, data.CreatedAt)
		if err != nil {
			return
		}
	}

	// Filter{}
	count := AccountCount(tx, Filter{Type: "account_credited"})
	if count != 3 {
		t.Errorf("Expected 3 got %d", count)
	}

	// Filter{From}
	count = AccountCount(tx, Filter{Type: "account_credited", From: &datasets[4].CreatedAt})
	if count != 2 {
		t.Errorf("Expected 2 got %d", count)
	}

	// Filter{To}
	count = AccountCount(tx, Filter{Type: "account_credited", To: &datasets[4].CreatedAt})
	if count != 2 {
		t.Errorf("Expected 2 got %d", count)
	}

	// Filter{From, To}
	count = AccountCount(tx, Filter{Type: "account_credited", From: &datasets[3].CreatedAt, To: &datasets[4].CreatedAt})
	if count != 1 {
		t.Errorf("Expected 1 got %d", count)
	}
}

func TestTotalIssued(t *testing.T) {
	db, err := sql.OpenAndMigrate("../../")
	if err != nil {
		t.Error(err)
	}

	tx, err := db.Beginx()
	if err != nil {
		t.Error(err)
	}
	defer tx.Rollback()

	for _, data := range datasets {
		_, err := sql.Exec(tx, `INSERT INTO effects(effect_id, operation, paging_token, account, amount, type, asset_type, asset_issuer, asset_code, created_at)
			                    VALUES('0034028708058632193-0000000001', 'https://horizon-testnet.stellar.org/operations/34028708058632193', '34028708058632193-1', $1, $2, $3, 'credit_alphanum4', 'GCJKC2MI63KSQ6MLE6GBSXPDKTDAK43WR522ZYR3F34NPM7Z5UEPIZNX', 'CNDY', $4)`,
			data.Account, data.Amount, data.Type, data.CreatedAt)
		if err != nil {
			return
		}
	}

	// Filter{}
	count := TotalIssued(tx, "GCJKC2MI63KSQ6MLE6GBSXPDKTDAK43WR522ZYR3F34NPM7Z5UEPIZNX", Filter{})
	if count != 1100 {
		t.Errorf("Expected 1100 got %f", count)
	}

	// Filter{From}
	count = TotalIssued(tx, "GCJKC2MI63KSQ6MLE6GBSXPDKTDAK43WR522ZYR3F34NPM7Z5UEPIZNX", Filter{From: &datasets[4].CreatedAt})
	if count != 100 {
		t.Errorf("Expected 100 got %f", count)
	}

	// Filter{To}
	count = TotalIssued(tx, "GCJKC2MI63KSQ6MLE6GBSXPDKTDAK43WR522ZYR3F34NPM7Z5UEPIZNX", Filter{To: &datasets[4].CreatedAt})
	if count != 1000 {
		t.Errorf("Expected 1000 got %f", count)
	}

	// Filter{From, To}
	count = TotalIssued(tx, "GCJKC2MI63KSQ6MLE6GBSXPDKTDAK43WR522ZYR3F34NPM7Z5UEPIZNX", Filter{From: &datasets[1].CreatedAt, To: &datasets[3].CreatedAt})
	if count != 1000 {
		t.Errorf("Expected 1000 got %f", count)
	}
}

func TestTotalCount(t *testing.T) {
	db, err := sql.OpenAndMigrate("../../")
	if err != nil {
		t.Error(err)
	}

	tx, err := db.Beginx()
	if err != nil {
		t.Error(err)
	}
	defer tx.Rollback()

	for _, data := range datasets {
		_, err := sql.Exec(tx, `INSERT INTO effects(effect_id, operation, paging_token, account, amount, type, asset_type, asset_issuer, asset_code, created_at)
			                    VALUES('0034028708058632193-0000000001', 'https://horizon-testnet.stellar.org/operations/34028708058632193', '34028708058632193-1', $1, $2, $3, 'credit_alphanum4', 'GCJKC2MI63KSQ6MLE6GBSXPDKTDAK43WR522ZYR3F34NPM7Z5UEPIZNX', 'CNDY', $4)`,
			data.Account, data.Amount, data.Type, data.CreatedAt)
		if err != nil {
			return
		}
	}

	// Filter{}
	count := TotalCount(tx, Filter{Type: "trustline_created"})
	if count != 2 {
		t.Errorf("Expected 2 got %d", count)
	}

	// Filter{From}
	count = TotalCount(tx, Filter{Type: "trustline_created", From: &datasets[3].CreatedAt})
	if count != 1 {
		t.Errorf("Expected 1 got %d", count)
	}

	// Filter{To}
	count = TotalCount(tx, Filter{Type: "trustline_created", To: &datasets[2].CreatedAt})
	if count != 1 {
		t.Errorf("Expected 1 got %d", count)
	}

	// Filter{From, To}
	count = TotalCount(tx, Filter{Type: "trustline_created", From: &datasets[1].CreatedAt, To: &datasets[2].CreatedAt})
	if count != 0 {
		t.Errorf("Expected 0 got %d", count)
	}
}

func TestItemCount(t *testing.T) {
	db, err := sql.OpenAndMigrate("../../")
	if err != nil {
		t.Error(err)
	}

	tx, err := db.Beginx()
	if err != nil {
		t.Error(err)
	}
	defer tx.Rollback()

	for _, data := range datasets {
		_, err := sql.Exec(tx, `INSERT INTO effects(effect_id, operation, paging_token, account, amount, type, asset_type, asset_issuer, asset_code, created_at)
			                    VALUES('0034028708058632193-0000000001', 'https://horizon-testnet.stellar.org/operations/34028708058632193', '34028708058632193-1', $1, $2, $3, 'credit_alphanum4', 'GCJKC2MI63KSQ6MLE6GBSXPDKTDAK43WR522ZYR3F34NPM7Z5UEPIZNX', 'CNDY', $4)`,
			data.Account, data.Amount, data.Type, data.CreatedAt)
		if err != nil {
			return
		}
	}

	// Filter{}
	count := ItemCount(tx, Filter{})
	if count != len(datasets) {
		t.Errorf("Expected %d got %d", len(datasets), count)
	}

	// Filter{From}
	count = ItemCount(tx, Filter{From: &datasets[3].CreatedAt})
	if count != 5 {
		t.Errorf("Expected 5 got %d", count)
	}

	// Filter{To}
	count = ItemCount(tx, Filter{To: &datasets[2].CreatedAt})
	if count != 3 {
		t.Errorf("Expected 3 got %d", count)
	}

	// Filter{From, To}
	count = ItemCount(tx, Filter{From: &datasets[1].CreatedAt, To: &datasets[4].CreatedAt})
	if count != 4 {
		t.Errorf("Expected 4 got %d", count)
	}
}
