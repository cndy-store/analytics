package assetStat

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/stellar/go/clients/horizon"
	"log"
	"time"
)

type AssetStat struct {
	Id          *uint32    `db:"id"           json:"-"`
	AssetType   *string    `db:"asset_type"   json:"asset_type,omitempty"`
	AssetCode   *string    `db:"asset_code"   json:"asset_code,omitempty"`
	AssetIssuer *string    `db:"asset_issuer" json:"asset_issuer,omitempty"`
	TotalAmount *string    `db:"total_amount" json:"total_amount,omitempty"`
	NumAccounts *int32     `db:"num_accounts" json:"num_accounts,omitempty"`
	NumEffects  *int32     `db:"num_effects"  json:"num_effects,omitempty"`
	CreatedAt   *time.Time `db:"created_at"   json:"created_at,omitempty"`
}

func New(db *sqlx.DB, effect horizon.Effect, timestamp time.Time) (err error) {
	// Store amount_transfered and amount_issued upon insert in a different table
	// (analogue to the asset endpoint of Horizon)
	_, err = db.Exec(`INSERT INTO asset_stats(asset_code, asset_issuer, asset_type, created_at, total_amount, num_accounts, num_effects)
	                  VALUES ($1, $2, $3, $4,
					         (SELECT SUM(amount) FROM effects WHERE type='account_debited' AND account=$1),
					         (SELECT COUNT(DISTINCT account) FROM effects),
							 (SELECT COUNT(*) FROM effects))`,
		effect.Asset.Code, effect.Asset.Issuer, effect.Asset.Type, timestamp)

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

func Get(db *sqlx.DB, filter Filter) (stats []AssetStat, err error) {
	filter.Defaults()
	err = db.Select(&stats, `SELECT * FROM asset_stats WHERE cast(strftime('%s', created_at) AS INT) BETWEEN $2 AND $3`,
		filter.From.Unix(), filter.To.Unix())
	if err == sql.ErrNoRows {
		log.Printf("[ERROR] asset_stat.Get(): %s", err)
	}
	return
}
