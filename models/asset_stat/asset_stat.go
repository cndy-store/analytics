package assetStat

import (
	"github.com/cndy-store/analytics/utils/bigint"
	"github.com/cndy-store/analytics/utils/sql"
	"github.com/stellar/go/clients/horizon"
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

func New(db sql.Database, effect horizon.Effect, timestamp time.Time) (err error) {
	// Store amount_transfered and amount_issued upon insert in a different table
	// (analogue to the asset endpoint of Horizon)

	_, err = db.Exec(`INSERT INTO asset_stats(paging_token, asset_code, asset_issuer, asset_type, created_at, issued, transferred, accounts_with_trustline, accounts_with_payments, payments)
		                   VALUES ($1, $2, $3, $4, $5,
		                       (SELECT COALESCE(SUM(amount), 0) FROM effects WHERE type='account_debited' AND account=$6),
		                       (SELECT COALESCE(SUM(amount), 0) FROM effects WHERE type='account_debited' AND account!=$6),
		                       (SELECT COUNT(DISTINCT account) FROM effects WHERE type='trustline_created' AND account!=$6),
		                       (SELECT COUNT(DISTINCT account) FROM effects WHERE type='account_debited' AND account!=$6),
		                       (SELECT COUNT(*) FROM effects WHERE type='account_debited' AND account!=$6)
		                   )`,
		effect.PT, effect.Asset.Code, effect.Asset.Issuer, effect.Asset.Type, timestamp, effect.Asset.Issuer)

	return
}

type Filter struct {
	From *time.Time
	To   *time.Time
}

func (f *Filter) Defaults() {
	if f.From == nil {
		t := time.Unix(0, 0)
		f.From = &t
	}

	if f.To == nil {
		t := time.Now()
		f.To = &t
	}
}

func Get(db sql.Database, filter Filter) (stats []AssetStat, err error) {
	filter.Defaults()
	err = db.Select(&stats, `SELECT * FROM asset_stats WHERE created_at BETWEEN $1::timestamp AND $2::timestamp ORDER BY id`,
		filter.From, filter.To)
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

func Latest(db sql.Database) (stats AssetStat, err error) {
	err = db.Get(&stats, `SELECT * FROM asset_stats ORDER BY id DESC LIMIT 1`)
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
