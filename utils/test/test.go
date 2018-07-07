package test

import (
	"fmt"
	"github.com/cndy-store/analytics/utils/bigint"
	"github.com/cndy-store/analytics/utils/cndy"
	"github.com/jmoiron/sqlx"
	"time"
)

type Effect struct {
	PagingToken string
	Account     string
	Amount      string
	Type        string

	// AssetStat tests
	Issued                int64
	Transferred           int64
	AccountsWithTrustline int32
	AccountsWithPayments  int32
	Payments              int32

	CreatedAt time.Time
}

// NOTE: CreatedAt timestamps need to be in the past
var Effects = []Effect{
	{"34028708058632193-0", "GDNH64DRUT4CY3UJLWQIB655PQ6OG34UGYB4NC5DC4TYWLNJIBCEYTTD", "", "trustline_created", 0, 0, 1, 0, 0, time.Date(2018, time.March, 4, 0, 0, 0, 0, time.UTC)},
	{"34028708058632193-1", "GDNH64DRUT4CY3UJLWQIB655PQ6OG34UGYB4NC5DC4TYWLNJIBCEYTTD", "1000.0000000", "account_credited", 0, 0, 1, 0, 0, time.Date(2018, time.March, 5, 0, 0, 0, 0, time.UTC)},
	{"34028708058632193-2", cndy.AssetIssuer, "1000.0000000", "account_debited", 10000000000, 0, 1, 0, 0, time.Date(2018, time.March, 6, 0, 0, 0, 0, time.UTC)},
	{"34028708058632193-3", "GBEYRLI7OCZU7JVT33GBVVI5XWVCQNSUI3TXDE7Z5MWC6CLQMSTPDT6A", "", "trustline_created", 10000000000, 0, 2, 0, 0, time.Date(2018, time.March, 7, 0, 0, 0, 0, time.UTC)},
	{"34028708058632193-4", "GBEYRLI7OCZU7JVT33GBVVI5XWVCQNSUI3TXDE7Z5MWC6CLQMSTPDT6A", "15.0000000", "account_credited", 10000000000, 0, 2, 0, 0, time.Date(2018, time.March, 8, 0, 0, 0, 0, time.UTC)},
	{"34028708058632193-5", "GDNH64DRUT4CY3UJLWQIB655PQ6OG34UGYB4NC5DC4TYWLNJIBCEYTTD", "15.0000000", "account_debited", 10000000000, 150000000, 2, 1, 1, time.Date(2018, time.March, 9, 0, 0, 0, 0, time.UTC)},
	{"34028708058632193-6", "GCFXHS4GXL6BVUCXBWXGTITROWLVYXQKQLF4YH5O5JT3YZXCYPAFBJZB", "", "trustline_created", 10000000000, 150000000, 3, 1, 1, time.Date(2018, time.March, 10, 0, 0, 0, 0, time.UTC)},
	{"34028708058632193-7", "GDNH64DRUT4CY3UJLWQIB655PQ6OG34UGYB4NC5DC4TYWLNJIBCEYTTD", "100.0000000", "account_credited", 10000000000, 150000000, 3, 1, 1, time.Date(2018, time.March, 11, 0, 0, 0, 0, time.UTC)},
	{"34028708058632193-8", "GCFXHS4GXL6BVUCXBWXGTITROWLVYXQKQLF4YH5O5JT3YZXCYPAFBJZB", "100.0000000", "account_debited", 10000000000, 1150000000, 3, 2, 2, time.Date(2018, time.March, 12, 0, 0, 0, 0, time.UTC)},
	{"34028708058632193-9", "GBEYRLI7OCZU7JVT33GBVVI5XWVCQNSUI3TXDE7Z5MWC6CLQMSTPDT6A", "10.0000000", "account_credited", 10000000000, 1150000000, 3, 2, 2, time.Date(2018, time.March, 13, 0, 0, 0, 0, time.UTC)},
	{"34028708058632194-0", "GCFXHS4GXL6BVUCXBWXGTITROWLVYXQKQLF4YH5O5JT3YZXCYPAFBJZB", "10.0000000", "account_debited", 10000000000, 1250000000, 3, 2, 3, time.Date(2018, time.March, 14, 0, 0, 0, 0, time.UTC)},
}

// Helper function to insert test data
// TODO: Write a test that asset_stats where correctly inserted?
func InsertTestData(tx *sqlx.Tx) (err error) {
	for i, data := range Effects {
		var amount *int64
		amount, err = bigint.Parse(data.Amount)
		if err != nil {
			return
		}

		_, err = tx.Exec(`INSERT INTO effects(effect_id, operation, paging_token, account, amount, type, asset_type, asset_issuer, asset_code, created_at)
			                    VALUES($1, 'https://horizon-testnet.stellar.org/operations/34028708058632193', $2, $3, $4, $5, 'credit_alphanum4', $6, $7, $8)`,
			fmt.Sprintf("0034028708058632193-%09d", i), data.PagingToken, data.Account, amount, data.Type, cndy.AssetIssuer, cndy.AssetCode, data.CreatedAt)
		if err != nil {
			return
		}
	}

	return
}
