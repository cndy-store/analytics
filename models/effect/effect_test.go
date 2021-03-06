package effect

import (
	"github.com/cndy-store/analytics/models/asset_stat"
	"github.com/cndy-store/analytics/utils/bigint"
	"github.com/cndy-store/analytics/utils/cndy"
	"github.com/cndy-store/analytics/utils/filter"
	"github.com/cndy-store/analytics/utils/sql"
	"github.com/cndy-store/analytics/utils/test"
	hProtocol "github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/protocols/horizon/base"
	"github.com/stellar/go/support/render/hal"
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
	operation := hal.Link{
		Href:      "https://horizon-testnet.stellar.org/operations/33820436504518657",
		Templated: false,
	}
	succeeds := hal.Link{
		Href:      "https://horizon-testnet.stellar.org/effects?order=desc&cursor=33820436504518657-1",
		Templated: false,
	}
	precedes := hal.Link{
		Href:      "https://horizon-testnet.stellar.org/effects?order=asc&cursor=33820436504518657-1",
		Templated: false,
	}

	asset := base.Asset{
		Type:   "credit_alphanum4",
		Code:   cndy.AssetCode,
		Issuer: cndy.AssetIssuer,
	}

	type links struct {
		Operation hal.Link
		Succeeds  hal.Link
		Precedes  hal.Link
	}

	// We need an account_debited effect from the issuer, so total_amount is not nil for the
	// asset stats tests later
	effect := hProtocol.Effect{
		ID:      "0033820436504518657-0000000001",
		PT:      "33820436504518657-1",
		Account: cndy.AssetIssuer,
		Amount:  "5.0000000",
		Type:    "account_debited",
		TypeI:   2,
		Balance: hProtocol.Balance{Asset: asset},
		Signer:  hProtocol.Signer{},
	}
	effect.Links.Operation = operation
	effect.Links.Succeeds = succeeds
	effect.Links.Precedes = precedes

	// TODO: Mock effect.GetOperationTime()
	err = New(tx, effect)
	if err != nil {
		t.Error(err)
	}

	effects, err := Get(tx, filter.NewCNDYFilter(nil, nil))
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
	a, err := assetStat.Latest(tx, filter.NewCNDYFilter(nil, nil))
	if err != nil {
		t.Error(err)
	}

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
	if *e.CreatedAt != *a.CreatedAt {
		t.Errorf("Expected %v got %v", *e.CreatedAt, *a.CreatedAt)
	}

	// We've used an account_debited type from the asset issuer
	// This should increase the total available supply
	if *e.Amount != *a.Issued {
		t.Errorf("Expected %v got %v", *e.Amount, *a.Issued)
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

	// NewCNDYFilter(nil, nil)
	effects, err := Get(tx, filter.NewCNDYFilter(nil, nil))
	if err != nil {
		t.Errorf("effect.Get(): %s", err)
	}
	if len(test.CNDYEffects) != len(effects) {
		t.Errorf("Expected %d effects got %d", len(test.CNDYEffects), len(effects))
	}

	// Filter{From}
	effects, err = Get(tx, filter.NewCNDYFilter(&test.CNDYEffects[5].CreatedAt, nil))
	if err != nil {
		t.Errorf("effect.Get(): %s", err)
	}
	if len(test.CNDYEffects[5:]) != len(effects) {
		t.Errorf("Expected %d effects got %d", len(test.CNDYEffects), len(effects))
	}

	for i, e := range test.CNDYEffects[5:] {

		if e.PagingToken != *effects[i].PagingToken {
			t.Errorf("Expected paging_token to be %s got: %s", e.PagingToken, *effects[i].PagingToken)
		}
	}

	// Filter{To}
	effects, err = Get(tx, filter.NewCNDYFilter(nil, &test.CNDYEffects[2].CreatedAt))
	if err != nil {
		t.Errorf("effect.Get(): %s", err)
	}
	if len(test.CNDYEffects[:3]) != len(effects) {
		t.Errorf("Expected %d effects got %d", len(test.CNDYEffects[:3]), len(effects))
	}

	for i, e := range test.CNDYEffects[:3] {

		if e.PagingToken != *effects[i].PagingToken {
			t.Errorf("Expected paging_token to be %s got: %s", e.PagingToken, *effects[i].PagingToken)
		}
	}

	// Filter{From, To}
	effects, err = Get(tx, filter.NewCNDYFilter(&test.CNDYEffects[3].CreatedAt, &test.CNDYEffects[4].CreatedAt))
	if err != nil {
		t.Errorf("effect.Get(): %s", err)
	}
	if len(test.CNDYEffects[3:5]) != len(effects) {
		t.Errorf("Expected %d effects got %d", len(test.CNDYEffects[3:5]), len(effects))
	}

	for i, e := range test.CNDYEffects[3:5] {
		if e.PagingToken != *effects[i].PagingToken {
			t.Errorf("Expected paging_token to be %s got: %s", e.PagingToken, *effects[i].PagingToken)
		}
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
