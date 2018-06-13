package test

import (
	"fmt"
	"github.com/cndy-store/analytics/utils/bigint"
	"github.com/cndy-store/analytics/utils/cndy"
	"github.com/cndy-store/analytics/utils/sql"
	"github.com/jmoiron/sqlx"
	"time"
)

type Effect struct {
	PagingToken string
	Account     string
	Amount      string
	Type        string
	CreatedAt   time.Time
}

var Effects = []Effect{
	{"34028708058632193-0", "GDNH64DRUT4CY3UJLWQIB655PQ6OG34UGYB4NC5DC4TYWLNJIBCEYTTD", "", "trustline_created", time.Date(2018, time.March, 8, 0, 0, 0, 0, time.UTC)},
	{"34028708058632193-1", "GDNH64DRUT4CY3UJLWQIB655PQ6OG34UGYB4NC5DC4TYWLNJIBCEYTTD", "1000.0000000", "account_credited", time.Date(2018, time.March, 10, 0, 0, 0, 0, time.UTC)},
	{"34028708058632193-2", cndy.AssetIssuer, "1000.0000000", "account_debited", time.Date(2018, time.March, 12, 0, 0, 0, 0, time.UTC)},
	{"34028708058632193-3", "GBEYRLI7OCZU7JVT33GBVVI5XWVCQNSUI3TXDE7Z5MWC6CLQMSTPDT6A", "", "trustline_created", time.Date(2018, time.March, 14, 0, 0, 0, 0, time.UTC)},
	{"34028708058632193-4", "GBEYRLI7OCZU7JVT33GBVVI5XWVCQNSUI3TXDE7Z5MWC6CLQMSTPDT6A", "15.0000000", "account_credited", time.Date(2018, time.March, 16, 0, 0, 0, 0, time.UTC)},
	{"34028708058632193-5", "GDNH64DRUT4CY3UJLWQIB655PQ6OG34UGYB4NC5DC4TYWLNJIBCEYTTD", "15.0000000", "account_debited", time.Date(2018, time.March, 18, 0, 0, 0, 0, time.UTC)},
	{"34028708058632193-6", "GDNH64DRUT4CY3UJLWQIB655PQ6OG34UGYB4NC5DC4TYWLNJIBCEYTTD", "100.0000000", "account_credited", time.Date(2018, time.March, 20, 0, 0, 0, 0, time.UTC)},
	{"34028708058632193-7", cndy.AssetIssuer, "100.0000000", "account_debited", time.Date(2018, time.March, 22, 0, 0, 0, 0, time.UTC)},
}

type AssetStat struct {
	PagingToken string
	TotalAmount int64
	NumAccounts int32
	Payments    int32
	CreatedAt   time.Time
}

var AssetStats = []AssetStat{
	{"39819440072110101-0", 10000000000, 10, 20, time.Date(2018, time.March, 12, 0, 0, 0, 0, time.UTC)},
	{"39819440072110101-1", 10000000000, 12, 25, time.Date(2018, time.March, 14, 0, 0, 0, 0, time.UTC)},
	{"39819440072110101-2", 20000000000, 15, 30, time.Date(2018, time.March, 16, 0, 0, 0, 0, time.UTC)},
	{"39819440072110101-3", 20000000000, 22, 35, time.Date(2018, time.March, 18, 0, 0, 0, 0, time.UTC)},
}

// Helper function to insert test data
func InsertEffects(tx *sqlx.Tx) (err error) {
	for i, data := range Effects {
		var amount *int64
		amount, err = bigint.Parse(data.Amount)
		if err != nil {
			return
		}

		_, err = sql.Exec(tx, `INSERT INTO effects(effect_id, operation, paging_token, account, amount, type, asset_type, asset_issuer, asset_code, created_at)
			                    VALUES($1, 'https://horizon-testnet.stellar.org/operations/34028708058632193', $2, $3, $4, $5, 'credit_alphanum4', $6, $7, $8)`,
			fmt.Sprintf("0034028708058632193-%09d", i), data.PagingToken, data.Account, amount, data.Type, cndy.AssetIssuer, cndy.AssetCode, data.CreatedAt)
		if err != nil {
			return
		}
	}

	return
}

func InsertAssetStats(tx *sqlx.Tx) (err error) {
	for _, data := range AssetStats {
		_, err = tx.Exec(`INSERT INTO asset_stats(paging_token, asset_code, asset_issuer, asset_type, total_amount, num_accounts, payments, created_at)
			              VALUES($1, $2, $3, 'credit_alphanum4', $4, $5, $6, $7)`,
			data.PagingToken, cndy.AssetCode, cndy.AssetIssuer, data.TotalAmount, data.NumAccounts, data.Payments, data.CreatedAt)
		if err != nil {
			return
		}
	}

	return
}
