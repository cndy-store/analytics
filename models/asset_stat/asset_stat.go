package assetStat

import (
	"github.com/cndy-store/analytics/utils/bigint"
	"github.com/cndy-store/analytics/utils/filter"
	"github.com/cndy-store/analytics/utils/sql"
	"time"
)

type AssetStat struct {
	EffectId    *string `db:"effect_id",   json:"effect_id"` // TODO: omitempty, also add to tests?
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

func Get(db sql.Database, filter filter.Filter) (stats []AssetStat, err error) {
	filter.Defaults()
	err = db.Select(&stats, `SELECT * FROM asset_stats($1, $2) WHERE created_at BETWEEN $3::timestamp AND $4::timestamp ORDER BY effect_id`,
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
	err = db.Get(&stats, `SELECT * FROM asset_stats($1, $2) ORDER BY effect_id DESC LIMIT 1`,
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
