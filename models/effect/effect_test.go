package effect

import (
	"fmt"
	"github.com/cndy-store/analytics/models/asset_stat"
	"github.com/cndy-store/analytics/utils/sql"
	"github.com/stellar/go/clients/horizon"
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

	// Construct effect to insert
	operation := horizon.Link{
		Href:      "https://horizon-testnet.stellar.org/operations/33820436504518657",
		Templated: false,
	}
	succeeds := horizon.Link{
		Href:      "https://horizon-testnet.stellar.org/effects?order=desc&cursor=33820436504518657-1",
		Templated: false,
	}
	precedes := horizon.Link{
		Href:      "https://horizon-testnet.stellar.org/effects?order=asc&cursor=33820436504518657-1",
		Templated: false,
	}

	asset := horizon.Asset{
		Type:   "credit_alphanum4",
		Code:   "CNDY",
		Issuer: "GCJKC2MI63KSQ6MLE6GBSXPDKTDAK43WR522ZYR3F34NPM7Z5UEPIZNX",
	}

	type links struct {
		Operation horizon.Link
		Succeeds  horizon.Link
		Precedes  horizon.Link
	}

	// We need an account_debited effect from the issuer, so total_amount is not nil for the
	// asset stats tests later
	effect := horizon.Effect{
		ID:      "0033820436504518657-0000000001",
		PT:      "33820436504518657-1",
		Account: "GCJKC2MI63KSQ6MLE6GBSXPDKTDAK43WR522ZYR3F34NPM7Z5UEPIZNX",
		Amount:  "5.0000000",
		Type:    "account_debited",
		TypeI:   2,
		Balance: horizon.Balance{Asset: asset},
		Signer:  horizon.Signer{},
	}
	effect.Links.Operation = operation
	effect.Links.Succeeds = succeeds
	effect.Links.Precedes = precedes

	err = New(tx, effect)
	if err != nil {
		t.Error(err)
	}

	effects, err := Get(tx, Filter{})
	if err != nil {
		t.Error(err)
	}

	if len(effects) != 1 {
		t.Errorf("Expected 1 result, got %d", len(effects))
	}

	e := effects[0]

	// Check links
	if operation.Href != *e.Operation {
		t.Errorf("Expected %v got %v", operation.Href, *e.Operation)
	}
	if succeeds.Href != *e.Succeeds {
		t.Errorf("Expected %v got %v", succeeds.Href, *e.Succeeds)
	}
	if precedes.Href != *e.Precedes {
		t.Errorf("Expected %v got %v", precedes.Href, *e.Precedes)
	}

	// Check asset
	if asset.Type != *e.AssetType {
		t.Errorf("Expected %v got %v", asset.Type, *e.AssetType)
	}
	if asset.Code != *e.AssetCode {
		t.Errorf("Expected %v got %v", asset.Code, *e.AssetCode)
	}
	if asset.Issuer != *e.AssetIssuer {
		t.Errorf("Expected %v got %v", asset.Issuer, *e.AssetIssuer)
	}

	// Check native fields
	if effect.ID != *e.EffectId {
		t.Errorf("Expected %v got %v", effect.ID, *e.EffectId)
	}
	if effect.PT != *e.PagingToken {
		t.Errorf("Expected %v got %v", effect.PT, *e.PagingToken)
	}
	if effect.Account != *e.Account {
		t.Errorf("Expected %v got %v", effect.Account, *e.Account)
	}
	if effect.Type != *e.Type {
		t.Errorf("Expected %v got %v", effect.Type, *e.Type)
	}
	if effect.TypeI != *e.TypeI {
		t.Errorf("Expected %v got %v", effect.TypeI, *e.TypeI)
	}

	// Check converted amount
	parsedAmount, err := parseInt64(effect.Amount)
	if err != nil {
		return
	}
	if *parsedAmount != *e.Amount {
		t.Errorf("Expected %v got %v", *parsedAmount, *e.Amount)
	}

	// Check whether asset_stat data was updated
	assetStats, err := assetStat.Get(tx, assetStat.Filter{})
	if err != nil {
		t.Error(err)
	}
	if len(assetStats) != 1 {
		t.Errorf("Expected 1 result, got %d", len(assetStats))
	}

	a := assetStats[0]

	if effect.PT != *a.PagingToken {
		t.Errorf("Expected %v got %v", effect.PT, *a.PagingToken)
	}
	if asset.Type != *a.AssetType {
		t.Errorf("Expected %v got %v", asset.Type, *a.AssetType)
	}
	if asset.Code != *a.AssetCode {
		t.Errorf("Expected %v got %v", asset.Code, *a.AssetCode)
	}
	if asset.Issuer != *a.AssetIssuer {
		t.Errorf("Expected %v got %v", asset.Issuer, *a.AssetIssuer)
	}
	if *e.Amount != *a.TotalAmount {
		t.Errorf("Expected %v got %v", e.Amount, *a.TotalAmount)
	}
	if *a.NumAccounts != 1 {
		t.Errorf("Expected 1 got %v", *a.NumAccounts)
	}
	if *a.NumEffects != 1 {
		t.Errorf("Expected 1 got %v", *a.NumEffects)
	}
	if *e.CreatedAt != *a.CreatedAt {
		t.Errorf("Expected %v got %v", e.CreatedAt, *a.CreatedAt)
	}
}

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

	err = insertData(tx)
	if err != nil {
		t.Error(err)
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
	db, err := sql.OpenAndMigrate("../..")
	if err != nil {
		t.Error(err)
	}

	tx, err := db.Beginx()
	if err != nil {
		t.Error(err)
	}
	defer tx.Rollback()

	err = insertData(tx)
	if err != nil {
		t.Error(err)
	}

	// Filter{}
	count := AccountCount(tx, Filter{})
	if count != 3 {
		t.Errorf("Expected 3 got %d", count)
	}

	// Filter{From}
	count = AccountCount(tx, Filter{From: &datasets[4].CreatedAt})
	if count != 3 {
		t.Errorf("Expected 3 got %d", count)
	}

	// Filter{To}
	count = AccountCount(tx, Filter{To: &datasets[4].CreatedAt})
	if count != 3 {
		t.Errorf("Expected 3 got %d", count)
	}

	// Filter{From, To}
	count = AccountCount(tx, Filter{From: &datasets[4].CreatedAt, To: &datasets[5].CreatedAt})
	if count != 2 {
		t.Errorf("Expected 2 got %d", count)
	}
}

func TestTotalIssued(t *testing.T) {
	db, err := sql.OpenAndMigrate("../..")
	if err != nil {
		t.Error(err)
	}

	tx, err := db.Beginx()
	if err != nil {
		t.Error(err)
	}
	defer tx.Rollback()

	err = insertData(tx)
	if err != nil {
		t.Error(err)
	}

	// Filter{}
	count := TotalIssued(tx, "GCJKC2MI63KSQ6MLE6GBSXPDKTDAK43WR522ZYR3F34NPM7Z5UEPIZNX", Filter{})
	if count != "1100.0000000" {
		t.Errorf("Expected 1100.0000000 got %s", count)
	}

	// Filter{From}
	count = TotalIssued(tx, "GCJKC2MI63KSQ6MLE6GBSXPDKTDAK43WR522ZYR3F34NPM7Z5UEPIZNX", Filter{From: &datasets[4].CreatedAt})
	if count != "100.0000000" {
		t.Errorf("Expected 100.0000000 got %s", count)
	}

	// Filter{To}
	count = TotalIssued(tx, "GCJKC2MI63KSQ6MLE6GBSXPDKTDAK43WR522ZYR3F34NPM7Z5UEPIZNX", Filter{To: &datasets[4].CreatedAt})
	if count != "1000.0000000" {
		t.Errorf("Expected 1000.0000000 got %s", count)
	}

	// Filter{From, To}
	count = TotalIssued(tx, "GCJKC2MI63KSQ6MLE6GBSXPDKTDAK43WR522ZYR3F34NPM7Z5UEPIZNX", Filter{From: &datasets[1].CreatedAt, To: &datasets[3].CreatedAt})
	if count != "1000.0000000" {
		t.Errorf("Expected 1000.0000000 got %s", count)
	}
}

func TestTotalAmount(t *testing.T) {
	db, err := sql.OpenAndMigrate("../..")
	if err != nil {
		t.Error(err)
	}

	tx, err := db.Beginx()
	if err != nil {
		t.Error(err)
	}
	defer tx.Rollback()

	err = insertData(tx)
	if err != nil {
		t.Error(err)
	}

	// Filter{}
	count := TotalAmount(tx, Filter{Type: "account_credited"})
	if count != "1115.0000000" {
		t.Errorf("Expected 1115.0000000 got %s", count)
	}

	// Filter{From}
	count = TotalAmount(tx, Filter{Type: "account_credited", From: &datasets[4].CreatedAt})
	if count != "115.0000000" {
		t.Errorf("Expected 115.0000000 got %s", count)
	}

	// Filter{To}
	count = TotalAmount(tx, Filter{Type: "account_credited", To: &datasets[4].CreatedAt})
	if count != "1015.0000000" {
		t.Errorf("Expected 1015.0000000 got %s", count)
	}

	// Filter{From, To}
	count = TotalAmount(tx, Filter{Type: "account_credited", From: &datasets[1].CreatedAt, To: &datasets[3].CreatedAt})
	if count != "1000.0000000" {
		t.Errorf("Expected 1000.0000000 got %s", count)
	}
}

func TestTotalCount(t *testing.T) {
	db, err := sql.OpenAndMigrate("../..")
	if err != nil {
		t.Error(err)
	}

	tx, err := db.Beginx()
	if err != nil {
		t.Error(err)
	}
	defer tx.Rollback()

	err = insertData(tx)
	if err != nil {
		t.Error(err)
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
	db, err := sql.OpenAndMigrate("../..")
	if err != nil {
		t.Error(err)
	}

	tx, err := db.Beginx()
	if err != nil {
		t.Error(err)
	}
	defer tx.Rollback()

	err = insertData(tx)
	if err != nil {
		t.Error(err)
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

func TestParseInt64(t *testing.T) {
	var tests = []struct {
		String string
		Int    int64
	}{
		{"", 0},
		{"5.0000000", 50000000},
		{"5.0", 50000000},
		{"5", 50000000},
		{"50000000", 500000000000000},
	}

	for _, test := range tests {
		res, err := parseInt64(test.String)
		if err != nil {
			t.Error(err)
		}

		if test.Int != *res {
			t.Errorf("Expected %d got %d", test.Int, *res)
		}
	}
}

func TestStringFromInt64(t *testing.T) {
	var tests = []struct {
		Int    int64
		String string
	}{
		{0, "0.0000000"},
		{5, "0.0000005"},
		{50000000, "5.0000000"},
	}

	for _, test := range tests {
		res := stringFromInt64(test.Int)
		if test.String != res {
			t.Errorf("Expected %s got %s", test.String, res)
		}
	}
}

// Helper function to insert test data
func insertData(tx interface{}) (err error) {
	for i, data := range datasets {
		var amount *int64
		amount, err = parseInt64(data.Amount)
		if err != nil {
			return
		}

		_, err = sql.Exec(tx, `INSERT INTO effects(effect_id, operation, paging_token, account, amount, type, asset_type, asset_issuer, asset_code, created_at)
			                    VALUES($1, 'https://horizon-testnet.stellar.org/operations/34028708058632193', $2, $3, $4, $5, 'credit_alphanum4', 'GCJKC2MI63KSQ6MLE6GBSXPDKTDAK43WR522ZYR3F34NPM7Z5UEPIZNX', 'CNDY', $6)`,
			fmt.Sprintf("0034028708058632193-000000000%d", i), data.PagingToken, data.Account, amount, data.Type, data.CreatedAt)
		if err != nil {
			return
		}
	}

	return
}
