package effect

import (
	"encoding/json"
	"github.com/cndy-store/analytics/models/asset_stat"
	"github.com/cndy-store/analytics/utils/bigint"
	"github.com/cndy-store/analytics/utils/sql"
	"github.com/stellar/go/clients/horizon"
	"log"
	"net/http"
	"time"
)

type Effect struct {
	EffectId        *string `db:"effect_id"        json:"id,omitempty"`
	Operation       *string `db:"operation"        json:"operation,omitempty"`
	Succeeds        *string `db:"succeeds"         json:"succeeds,omitempty"`
	Precedes        *string `db:"precedes"         json:"precedes,omitempty"`
	PagingToken     *string `db:"paging_token"     json:"paging_token,omitempty"`
	Account         *string `db:"account"          json:"account,omitempty"`
	Amount          *int64  `db:"amount"           json:"-"`
	Type            *string `db:"type"             json:"type,omitempty"`
	TypeI           *int32  `db:"type_i"           json:"type_i,omitempty"`
	StartingBalance *string `db:"starting_balance" json:"starting_balance,omitempty"`

	Balance      *int64 `db:"balance"       json:"-"`
	BalanceLimit *int64 `db:"balance_limit" json:"-"`

	AssetType   *string `db:"asset_type"   json:"asset_type,omitempty"`
	AssetCode   *string `db:"asset_code"   json:"asset_code,omitempty"`
	AssetIssuer *string `db:"asset_issuer" json:"asset_issuer,omitempty"`

	SignerPublicKey *string `db:"signer_public_key" json:"signer_public_key,omitempty"`
	SignerWeight    *int32  `db:"signer_weight"     json:"signer_weight,omitempty"`
	SignerKey       *string `db:"signer_key"        json:"signer_key,omitempty"`
	SignerType      *string `db:"signer_type"       json:"signer_type,omitempty"`

	CreatedAt *time.Time `db:"created_at"  json:"created_at,omitempty"`

	// These fields are used by .Convert()
	JsonAmount       *string `db:"-" json:"amount,omitempty"`
	JsonBalance      *string `db:"-" json:"balance,omitempty"`
	JsonBalanceLimit *string `db:"-" json:"balance_limit,omitempty"`
}

type Operation struct {
	From      string    `json:"from"`
	To        string    `json:"to"`
	CreatedAt time.Time `json:"created_at"`
}

func New(db interface{}, effect horizon.Effect) (err error) {
	// Get operation
	operation := getOperation(effect.Links.Operation.Href)

	// Check whether sender and receiver are the same in account_credited and account_debited
	// operations. If so, ignore this effect.
	// This prevents issues from unchecked transaction to oneself, see:
	// https://stellar.stackexchange.com/questions/1036/why-are-overspending-transactions-added-to-the-ledger
	if operation.To == operation.From && (effect.Asset.Type == "account_credited" || effect.Asset.Type == "account_debited") {
		log.Printf("[WARN] Found effect where sender and receiver are the same account, ignoring.")
		log.Printf("[DEBUG] %+v", effect)
		return
	}

	// Parse strings into integers
	parsedAmount, err := bigint.Parse(effect.Amount)
	if err != nil {
		return
	}
	parsedBalance, err := bigint.Parse(effect.Balance.Balance)
	if err != nil {
		return
	}
	parsedBalanceLimit, err := bigint.Parse(effect.Balance.Limit)
	if err != nil {
		return
	}

	// Just input the fields we're requiring for now, can be replayed anytime form the chain later.
	_, err = sql.Exec(db, `INSERT INTO effects(
			effect_id,
			operation, succeeds, precedes,
			paging_token, account, amount, type, type_i, starting_balance,
			balance, balance_limit,
			asset_type, asset_issuer, asset_code,
			signer_public_key, signer_weight, signer_key, signer_type,
			created_at
		)
		VALUES(
			$1,
			$2, $3, $4,
			$5, $6, $7, $8, $9, $10,
			$11, $12,
			$13, $14, $15,
			$16, $17, $18, $19,
			$20
		)`,
		effect.ID,
		effect.Links.Operation.Href, effect.Links.Succeeds.Href, effect.Links.Precedes.Href,
		effect.PT, effect.Account, parsedAmount, effect.Type, effect.TypeI, effect.StartingBalance,
		parsedBalance, parsedBalanceLimit,
		effect.Asset.Type, effect.Asset.Issuer, effect.Asset.Code,
		effect.Signer.PublicKey, effect.Signer.Weight, effect.Signer.Key, effect.Signer.Type,
		operation.CreatedAt)
	if err != nil {
		return
	}

	// Store asset stats upon insert in a different table
	err = assetStat.New(db, effect, operation.CreatedAt)
	if err != nil {
		return
	}

	log.Printf("--+--[ %s ]", effect.Asset.Code)
	log.Printf("  |")
	log.Printf("  +->  Type:      %s", effect.Type)
	log.Printf("  +->  Account:   %s", effect.Account)
	log.Printf("  +->  Amount:    %s", effect.Amount)
	log.Printf("  +->  Timestamp: %s\n\n", operation.CreatedAt)
	return
}

type Filter struct {
	Type string
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

func TotalAmount(db interface{}, filter Filter) (amount int64) {
	filter.Defaults()
	if filter.Type == "" {
		log.Printf("[ERROR] effect.TotalAmount(): No type given.")
		return
	}

	err := sql.Get(db, &amount, `SELECT SUM(amount) FROM effects WHERE type=$1 AND created_at BETWEEN $2::timestamp AND $3::timestamp`,
		filter.Type, filter.From, filter.To)
	if err != nil {
		log.Print(err)
		return
	}

	return
}

// Total assets issued
func TotalIssued(db interface{}, issuer string, filter Filter) (issued int64) {
	filter.Defaults()

	err := sql.Get(db, &issued, `SELECT SUM(amount) FROM effects WHERE type='account_debited' AND account=$1 AND created_at BETWEEN $2::timestamp AND $3::timestamp`,
		issuer, filter.From, filter.To)
	if err != nil {
		log.Print(err)
		return
	}

	return
}

func TotalCount(db interface{}, filter Filter) (count int) {
	filter.Defaults()
	if filter.Type == "" {
		log.Printf("[ERROR] effect.TotalCount(): No type given.")
		return
	}

	err := sql.Get(db, &count, `SELECT COUNT(*) FROM effects WHERE type=$1 AND created_at BETWEEN $2::timestamp AND $3::timestamp`,
		filter.Type, filter.From, filter.To)
	if err != nil {
		log.Printf("[ERROR] effect.TotalCount(): %s", err)
	}
	return
}

func AccountCount(db interface{}, filter Filter) (count int) {
	filter.Defaults()
	err := sql.Get(db, &count, `SELECT COUNT(DISTINCT account) FROM effects WHERE created_at BETWEEN $1::timestamp AND $2::timestamp`,
		filter.From, filter.To)
	if err != nil {
		log.Printf("[ERROR] effect.AccountCount(): %s", err)
	}
	return
}

func ItemCount(db interface{}, filter Filter) (count int) {
	filter.Defaults()
	err := sql.Get(db, &count, `SELECT COUNT(*) FROM effects WHERE created_at BETWEEN $1::timestamp AND $2::timestamp`,
		filter.From, filter.To)
	if err != nil {
		log.Printf("[ERROR] effect.ItemCount(): %s", err)
	}
	return
}

func Get(db interface{}, filter Filter) (effects []Effect, err error) {
	filter.Defaults()
	err = sql.Select(db, &effects, `SELECT * FROM effects WHERE created_at BETWEEN $1::timestamp AND $2::timestamp ORDER BY created_at`,
		filter.From, filter.To)
	if err == sql.ErrNoRows {
		log.Printf("[ERROR] effect.Get(): %s", err)
	}
	return
}

// Convert int64 fields of to strings
func (e *Effect) Convert() {
	if e.Amount != nil {
		amount := bigint.ToString(*e.Amount)
		e.JsonAmount = &amount
	}

	if e.Balance != nil {
		balance := bigint.ToString(*e.Balance)
		e.JsonBalance = &balance
	}

	if e.BalanceLimit != nil {
		balanceLimit := bigint.ToString(*e.BalanceLimit)
		e.JsonBalanceLimit = &balanceLimit
	}
}

func getOperation(url string) (op Operation) {
	var h = &http.Client{Timeout: 2 * time.Second}

	r, err := h.Get(url)
	if err != nil {
		log.Printf("[ERROR] effect.getOperationTime() HTTP request 'GET %s'", url)
		log.Printf("        %s", err)
		return
	}
	defer r.Body.Close()

	err = json.NewDecoder(r.Body).Decode(&op)
	if err != nil {
		log.Printf("[ERROR] effect.getOperationTime(): Couldn't decode JSON body: %s", err)
		return
	}
	return
}
