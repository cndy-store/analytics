package assetStat

import (
	"github.com/cndy-store/analytics/utils/filter"
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
	assetStats, err := Get(tx, filter.NewCNDYFilter(nil, nil))
	if err != nil {
		t.Errorf("assetStat.Get(): %s", err)
	}
	if len(test.Effects) != len(assetStats) {
		t.Errorf("Expected %d assetStats got %d", len(test.Effects), len(assetStats))
	}

	// Filter{From}
	assetStats, err = Get(tx, filter.NewCNDYFilter(&test.Effects[2].CreatedAt, nil))
	if err != nil {
		t.Errorf("assetStat.Get(): %s", err)
	}
	if len(test.Effects[2:]) != len(assetStats) {
		t.Errorf("Expected %d assetStats got %d", len(test.Effects[2:]), len(assetStats))
	}

	for i, e := range test.Effects[2:] {
		if e.Issued != *assetStats[i].Issued {
			t.Errorf("Expected %d got: %d", e.Issued, *assetStats[i].Issued)
		}

		if e.Transferred != *assetStats[i].Transferred {
			t.Errorf("Expected %d got: %d", e.Transferred, *assetStats[i].Transferred)
		}

		if e.AccountsWithTrustline != *assetStats[i].AccountsWithTrustline {
			t.Errorf("Expected %d got: %d", e.AccountsWithTrustline, *assetStats[i].AccountsWithTrustline)
		}

		if e.AccountsWithPayments != *assetStats[i].AccountsWithPayments {
			t.Errorf("Expected %d got: %d", e.AccountsWithPayments, *assetStats[i].AccountsWithPayments)
		}

		if e.Payments != *assetStats[i].Payments {
			t.Errorf("Expected %d got: %d", e.Payments, *assetStats[i].Payments)
		}
	}

	// Filter{To}
	assetStats, err = Get(tx, filter.NewCNDYFilter(nil, &test.Effects[1].CreatedAt))
	if err != nil {
		t.Errorf("assetStat.Get(): %s", err)
	}
	if len(test.Effects[:2]) != len(assetStats) {
		t.Errorf("Expected %d assetStats got %d", len(test.Effects[:2]), len(assetStats))
	}

	for i, e := range test.Effects[:2] {
		if e.Issued != *assetStats[i].Issued {
			t.Errorf("Expected %d got: %d", e.Issued, *assetStats[i].Issued)
		}

		if e.Transferred != *assetStats[i].Transferred {
			t.Errorf("Expected %d got: %d", e.Transferred, *assetStats[i].Transferred)
		}

		if e.AccountsWithTrustline != *assetStats[i].AccountsWithTrustline {
			t.Errorf("Expected %d got: %d", e.AccountsWithTrustline, *assetStats[i].AccountsWithTrustline)
		}

		if e.AccountsWithPayments != *assetStats[i].AccountsWithPayments {
			t.Errorf("Expected %d got: %d", e.AccountsWithPayments, *assetStats[i].AccountsWithPayments)
		}

		if e.Payments != *assetStats[i].Payments {
			t.Errorf("Expected %d got: %d", e.Payments, *assetStats[i].Payments)
		}
	}

	// Filter{From, To}
	assetStats, err = Get(tx, filter.NewCNDYFilter(&test.Effects[1].CreatedAt, &test.Effects[2].CreatedAt))
	if err != nil {
		t.Errorf("assetStat.Get(): %s", err)
	}
	if len(test.Effects[1:3]) != len(assetStats) {
		t.Errorf("Expected %d assetStats got %d", len(test.Effects[1:3]), len(assetStats))
	}

	for i, e := range test.Effects[1:3] {
		if e.Issued != *assetStats[i].Issued {
			t.Errorf("Expected %d got: %d", e.Issued, *assetStats[i].Issued)
		}

		if e.Transferred != *assetStats[i].Transferred {
			t.Errorf("Expected %d got: %d", e.Transferred, *assetStats[i].Transferred)
		}

		if e.AccountsWithTrustline != *assetStats[i].AccountsWithTrustline {
			t.Errorf("Expected %d got: %d", e.AccountsWithTrustline, *assetStats[i].AccountsWithTrustline)
		}

		if e.AccountsWithPayments != *assetStats[i].AccountsWithPayments {
			t.Errorf("Expected %d got: %d", e.AccountsWithPayments, *assetStats[i].AccountsWithPayments)
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
	assetStats, err := Latest(tx, filter.NewCNDYFilter(nil, nil))
	if err != nil {
		t.Errorf("assetStat.Latest(): %s", err)
	}

	// Compare to latest test.Effects[]
	lastEffect := test.Effects[len(test.Effects)-1]

	if lastEffect.Issued != *assetStats.Issued {
		t.Errorf("Expected %d got: %d", lastEffect.Issued, *assetStats.Issued)
	}

	if lastEffect.Transferred != *assetStats.Transferred {
		t.Errorf("Expected %d got: %d", lastEffect.Transferred, *assetStats.Transferred)
	}

	if lastEffect.AccountsWithTrustline != *assetStats.AccountsWithTrustline {
		t.Errorf("Expected %d got: %d", lastEffect.AccountsWithTrustline, *assetStats.AccountsWithTrustline)
	}

	if lastEffect.AccountsWithPayments != *assetStats.AccountsWithPayments {
		t.Errorf("Expected %d got: %d", lastEffect.AccountsWithPayments, *assetStats.AccountsWithPayments)
	}

	if lastEffect.Payments != *assetStats.Payments {
		t.Errorf("Expected %d got: %d", lastEffect.Payments, *assetStats.Payments)
	}
}

func TestConvert(t *testing.T) {
	amount := int64(500000000)
	assetStat := AssetStat{Issued: &amount}
	assetStat.Convert()

	if *assetStat.JsonIssued != "50.0000000" {
		t.Errorf("Expected 50.0000000, got %s", *assetStat.JsonIssued)
	}
}
