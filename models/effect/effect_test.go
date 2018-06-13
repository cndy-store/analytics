package effect

import (
	"github.com/cndy-store/analytics/models/asset_stat"
	"github.com/cndy-store/analytics/utils/bigint"
	"github.com/cndy-store/analytics/utils/cndy"
	"github.com/cndy-store/analytics/utils/sql"
	"github.com/cndy-store/analytics/utils/test"
	"github.com/stellar/go/clients/horizon"
	"testing"
)

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
		Code:   cndy.AssetCode,
		Issuer: cndy.AssetIssuer,
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
		Account: cndy.AssetIssuer,
		Amount:  "5.0000000",
		Type:    "account_debited",
		TypeI:   2,
		Balance: horizon.Balance{Asset: asset},
		Signer:  horizon.Signer{},
	}
	effect.Links.Operation = operation
	effect.Links.Succeeds = succeeds
	effect.Links.Precedes = precedes

	// TODO: Mock effect.GetOperationTime()
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
	parsedAmount, err := bigint.Parse(effect.Amount)
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
	if *a.Payments != 1 {
		t.Errorf("Expected 1 got %v", *a.Payments)
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

	err = test.InsertTestData(tx)
	if err != nil {
		t.Error(err)
	}

	// Filter{}
	effects, err := Get(tx, Filter{})
	if err != nil {
		t.Errorf("effect.Get(): %s", err)
	}
	if len(test.Effects) != len(effects) {
		t.Errorf("Expected %d effects got %d", len(test.Effects), len(effects))
	}

	// Filter{From}
	effects, err = Get(tx, Filter{From: &test.Effects[5].CreatedAt})
	if err != nil {
		t.Errorf("effect.Get(): %s", err)
	}
	if len(effects) != 3 {
		t.Errorf("Expected 3 effects got %d", len(effects))
	}
	if test.Effects[5].PagingToken != *effects[0].PagingToken {
		t.Errorf("Expected paging_token to be %s got: %s", test.Effects[3].PagingToken, *effects[0].PagingToken)
	}
	if test.Effects[6].PagingToken != *effects[1].PagingToken {
		t.Errorf("Expected paging_token to be %s got: %s", test.Effects[4].PagingToken, *effects[1].PagingToken)
	}
	if test.Effects[7].PagingToken != *effects[2].PagingToken {
		t.Errorf("Expected paging_token to be %s got: %s", test.Effects[5].PagingToken, *effects[2].PagingToken)
	}

	// Filter{To}
	effects, err = Get(tx, Filter{To: &test.Effects[2].CreatedAt})
	if err != nil {
		t.Errorf("effect.Get(): %s", err)
	}
	if len(effects) != 3 {
		t.Errorf("Expected 3 effects got %d", len(effects))
	}
	if test.Effects[0].PagingToken != *effects[0].PagingToken {
		t.Errorf("Expected paging_token to be %s got: %s", test.Effects[0].PagingToken, *effects[0].PagingToken)
	}
	if test.Effects[1].PagingToken != *effects[1].PagingToken {
		t.Errorf("Expected paging_token to be %s got: %s", test.Effects[1].PagingToken, *effects[1].PagingToken)
	}
	if test.Effects[2].PagingToken != *effects[2].PagingToken {
		t.Errorf("Expected paging_token to be %s got: %s", test.Effects[2].PagingToken, *effects[2].PagingToken)
	}

	// Filter{From, To}
	effects, err = Get(tx, Filter{From: &test.Effects[3].CreatedAt, To: &test.Effects[4].CreatedAt})
	if err != nil {
		t.Errorf("effect.Get(): %s", err)
	}
	if len(effects) != 2 {
		t.Errorf("Expected 2 effects got %d", len(effects))
	}
	if test.Effects[3].PagingToken != *effects[0].PagingToken {
		t.Errorf("Expected paging_token to be %s got: %s", test.Effects[3].PagingToken, *effects[0].PagingToken)
	}
	if test.Effects[4].PagingToken != *effects[1].PagingToken {
		t.Errorf("Expected paging_token to be %s got: %s", test.Effects[4].PagingToken, *effects[1].PagingToken)
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

	err = test.InsertTestData(tx)
	if err != nil {
		t.Error(err)
	}

	// Filter{}
	count := AccountCount(tx, Filter{})
	if count != 3 {
		t.Errorf("Expected 3 got %d", count)
	}

	// Filter{From}
	count = AccountCount(tx, Filter{From: &test.Effects[4].CreatedAt})
	if count != 3 {
		t.Errorf("Expected 3 got %d", count)
	}

	// Filter{To}
	count = AccountCount(tx, Filter{To: &test.Effects[4].CreatedAt})
	if count != 3 {
		t.Errorf("Expected 3 got %d", count)
	}

	// Filter{From, To}
	count = AccountCount(tx, Filter{From: &test.Effects[4].CreatedAt, To: &test.Effects[5].CreatedAt})
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

	err = test.InsertTestData(tx)
	if err != nil {
		t.Error(err)
	}

	// Filter{}
	count := TotalIssued(tx, cndy.AssetIssuer, Filter{})
	if count != 11000000000 {
		t.Errorf("Expected 1100.0000000 got %d", count)
	}

	// Filter{From}
	count = TotalIssued(tx, cndy.AssetIssuer, Filter{From: &test.Effects[4].CreatedAt})
	if count != 1000000000 {
		t.Errorf("Expected 100.0000000 got %d", count)
	}

	// Filter{To}
	count = TotalIssued(tx, cndy.AssetIssuer, Filter{To: &test.Effects[4].CreatedAt})
	if count != 10000000000 {
		t.Errorf("Expected 1000.0000000 got %d", count)
	}

	// Filter{From, To}
	count = TotalIssued(tx, cndy.AssetIssuer, Filter{From: &test.Effects[1].CreatedAt, To: &test.Effects[3].CreatedAt})
	if count != 10000000000 {
		t.Errorf("Expected 1000.0000000 got %d", count)
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

	err = test.InsertTestData(tx)
	if err != nil {
		t.Error(err)
	}

	// Filter{}
	count := TotalAmount(tx, Filter{Type: "account_credited"})
	if count != 11150000000 {
		t.Errorf("Expected 1115.0000000 got %d", count)
	}

	// Filter{From}
	count = TotalAmount(tx, Filter{Type: "account_credited", From: &test.Effects[4].CreatedAt})
	if count != 1150000000 {
		t.Errorf("Expected 115.0000000 got %d", count)
	}

	// Filter{To}
	count = TotalAmount(tx, Filter{Type: "account_credited", To: &test.Effects[4].CreatedAt})
	if count != 10150000000 {
		t.Errorf("Expected 1015.0000000 got %d", count)
	}

	// Filter{From, To}
	count = TotalAmount(tx, Filter{Type: "account_credited", From: &test.Effects[1].CreatedAt, To: &test.Effects[3].CreatedAt})
	if count != 10000000000 {
		t.Errorf("Expected 1000.0000000 got %d", count)
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

	err = test.InsertTestData(tx)
	if err != nil {
		t.Error(err)
	}

	// Filter{}
	count := TotalCount(tx, Filter{Type: "trustline_created"})
	if count != 2 {
		t.Errorf("Expected 2 got %d", count)
	}

	// Filter{From}
	count = TotalCount(tx, Filter{Type: "trustline_created", From: &test.Effects[3].CreatedAt})
	if count != 1 {
		t.Errorf("Expected 1 got %d", count)
	}

	// Filter{To}
	count = TotalCount(tx, Filter{Type: "trustline_created", To: &test.Effects[2].CreatedAt})
	if count != 1 {
		t.Errorf("Expected 1 got %d", count)
	}

	// Filter{From, To}
	count = TotalCount(tx, Filter{Type: "trustline_created", From: &test.Effects[1].CreatedAt, To: &test.Effects[2].CreatedAt})
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

	err = test.InsertTestData(tx)
	if err != nil {
		t.Error(err)
	}

	// Filter{}
	count := ItemCount(tx, Filter{})
	if count != len(test.Effects) {
		t.Errorf("Expected %d got %d", len(test.Effects), count)
	}

	// Filter{From}
	count = ItemCount(tx, Filter{From: &test.Effects[3].CreatedAt})
	if count != 5 {
		t.Errorf("Expected 5 got %d", count)
	}

	// Filter{To}
	count = ItemCount(tx, Filter{To: &test.Effects[2].CreatedAt})
	if count != 3 {
		t.Errorf("Expected 3 got %d", count)
	}

	// Filter{From, To}
	count = ItemCount(tx, Filter{From: &test.Effects[1].CreatedAt, To: &test.Effects[4].CreatedAt})
	if count != 4 {
		t.Errorf("Expected 4 got %d", count)
	}
}

func TestConvert(t *testing.T) {
	amount := int64(500000000)
	effect := Effect{Amount: &amount, Balance: &amount, BalanceLimit: &amount}
	effect.Convert()

	if *effect.JsonAmount != "50.0000000" {
		t.Errorf("Expected 50.0000000, got %s", *effect.JsonAmount)
	}
	if *effect.JsonBalance != "50.0000000" {
		t.Errorf("Expected 50.0000000, got %s", *effect.JsonBalance)
	}
	if *effect.JsonBalanceLimit != "50.0000000" {
		t.Errorf("Expected 50.0000000, got %s", *effect.JsonBalanceLimit)
	}
}
