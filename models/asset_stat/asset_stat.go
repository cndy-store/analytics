package assetStat

import (
	"github.com/cndy-store/analytics/utils/bigint"
	"github.com/cndy-store/analytics/utils/filter"
	"github.com/cndy-store/analytics/utils/sql"
	hProtocol "github.com/stellar/go/protocols/horizon"
	"time"
)

type AssetStat struct {
	Id          *uint32 `db:"id"           json:"-"`
	PagingToken *string `db:"paging_token" json:"paging_token,omitempty"`
	AssetType   *string `db:"asset_type"   json:"asset_type,omitempty"`
	AssetCode   *string `db:"asset_code"   json:"asset_code,omitempty"`
	AssetIssuer *string `db:"asset_issuer" json:"asset_issuer,omitempty"`

	Issued      *int64 `db:"issued"      json:"-"`
	Transferred *int64 `db:"transferred" json:"-"`
	Payments    *int32 `db:"payments"    json:"payments"`

	AccountsWithTrustline *int32 `db:"accounts_with_trustline" json:"accounts_with_trustline"`
	AccountsWithPayments  *int32 `db:"accounts_with_payments"  json:"accounts_with_payments"`

	CreatedAt *time.Time `db:"created_at" json:"created_at,omitempty"`

	// These fields are used by .Convert()
	JsonIssued      *string `db:"-" json:"issued"`
	JsonTransferred *string `db:"-" json:"transferred"`
}

func New(db sql.Database, effect hProtocol.Effect, timestamp time.Time) (err error) {
	// Store amount_transfered and amount_issued upon insert in a different table
	// (analogue to the asset endpoint of Horizon)

	// This INSERT statement is closely related to the repopulate_asset_stats() procedure
	// See database migrations for reference
	_, err = db.Exec(`INSERT INTO asset_stats(paging_token, asset_code, asset_issuer, asset_type, created_at, issued, transferred, accounts_with_trustline, accounts_with_payments, payments)
		                  VALUES ($1, $2, $3, $4, $5,
		                      (SELECT COALESCE(SUM(amount), 0) FROM effects
						       WHERE type='account_debited'
							       AND account=$6
							       AND asset_code=$7
							       AND asset_issuer=$6),
		                      (SELECT COALESCE(SUM(amount), 0) FROM effects
						       WHERE type='account_debited'
							       AND account!=$6
							       AND asset_code=$7
							       AND asset_issuer=$6),
		                      (SELECT COUNT(DISTINCT account) FROM effects
						       WHERE type='trustline_created'
							       AND account!=$6
							       AND asset_code=$7
							       AND asset_issuer=$6),
		                      (SELECT COUNT(DISTINCT account) FROM effects
						       WHERE type='account_debited'
							       AND account!=$6
							       AND asset_code=$7
							       AND asset_issuer=$6),
		                      (SELECT COUNT(*) FROM effects
						       WHERE type='account_debited'
							       AND account!=$6
							       AND asset_code=$7
							       AND asset_issuer=$6)
		                   )`,
		effect.PT, effect.Asset.Code, effect.Asset.Issuer, effect.Asset.Type, timestamp, effect.Asset.Issuer, effect.Asset.Code)
	return
}

func Get(db sql.Database, filter filter.Filter) (stats []AssetStat, err error) {
	filter.Defaults()
	err = db.Select(&stats, `SELECT * FROM asset_stats WHERE asset_code=$1 AND asset_issuer=$2 AND created_at BETWEEN $3::timestamp AND $4::timestamp ORDER BY id`,
		filter.AssetCode, filter.AssetIssuer, filter.From, filter.To)
	if err == sql.ErrNoRows {
		err = nil
		return
	}
	if err != nil {
		return
	}

	// Convert int64 fields to strings
	for i, _ := range stats {
		stats[i].Convert()
	}

	return
}

func Latest(db sql.Database, filter filter.Filter) (stats AssetStat, err error) {
	filter.Defaults()
	err = db.Get(&stats, `SELECT * FROM asset_stats WHERE asset_code=$1 AND asset_issuer=$2 ORDER BY id DESC LIMIT 1`,
		filter.AssetCode, filter.AssetIssuer)
	if err == sql.ErrNoRows {
		err = nil
		return
	}
	if err != nil {
		return
	}

	// Convert int64 fields to strings
	stats.Convert()
	return
}

// Convert int64 fields of to strings
func (a *AssetStat) Convert() {
	if a.Issued != nil {
		issued := bigint.ToString(*a.Issued)
		a.JsonIssued = &issued
	} else {
		issued := "0.0000000"
		a.JsonIssued = &issued
	}

	if a.Transferred != nil {
		transferred := bigint.ToString(*a.Transferred)
		a.JsonTransferred = &transferred
	} else {
		transferred := "0.0000000"
		a.JsonTransferred = &transferred
	}

	if a.AccountsWithTrustline == nil {
		accounts := int32(0)
		a.AccountsWithTrustline = &accounts
	}

	if a.AccountsWithPayments == nil {
		accounts := int32(0)
		a.AccountsWithPayments = &accounts
	}
	if a.Payments == nil {
		payments := int32(0)
		a.Payments = &payments
	}
}
