package effect

import (
	"database/sql"
	"encoding/json"
	"github.com/cndy-store/analytics/models/asset_stat"
	"github.com/jmoiron/sqlx"
	"github.com/stellar/go/clients/horizon"
	"log"
	"net/http"
	"time"
)

type Effect struct {
	Id              *uint32  `db:"id"               json:"-"`
	EffectId        *string  `db:"effect_id"        json:"id,omitempty"`
	Operation       *string  `db:"operation"        json:"operation,omitempty"`
	Succeeds        *string  `db:"succeeds"         json:"succeeds,omitempty"` // Currently not used
	Precedes        *string  `db:"precedes"         json:"precedes,omitempty"` // Currently not used
	PagingToken     *string  `db:"paging_token"     json:"paging_token,omitempty"`
	Account         *string  `db:"account"          json:"account,omitempty"`
	Amount          *float32 `db:"amount"           json:"amount,omitempty"`
	Type            *string  `db:"type"             json:"type,omitempty"`
	TypeI           *int32   `db:"type_i"           json:"type_i,omitempty"` // Currently not used
	StartingBalance *string  `db:"starting_balance" json:"starting_balance,omitempty"`

	Balance      *string `db:"balance"       json:"balance,omitempty"`
	BalanceLimit *string `db:"balance_limit" json:"balance_limit,omitempty"`

	AssetType   *string `db:"asset_type"   json:"asset_type,omitempty"`
	AssetCode   *string `db:"asset_code"   json:"asset_code,omitempty"`
	AssetIssuer *string `db:"asset_issuer" json:"asset_issuer,omitempty"`

	SignerPublicKey *string `db:"signer_public_key" json:"signer_public_key,omitempty"`
	SignerWeight    *string `db:"signer_weight"     json:"signer_weight,omitempty"`
	SignerKey       *string `db:"signer_key"        json:"signer_key,omitempty"`
	SignerType      *string `db:"signer_type"       json:"signer_type,omitempty"`

	CreatedAt *time.Time `db:"created_at"  json:"created_at,omitempty"`
}

type Operation struct {
	From      string    `json:"from"`
	To        string    `json:"to"`
	CreatedAt time.Time `json:"created_at"`
}

func New(db *sqlx.DB, effect horizon.Effect) (err error) {
	// Get operation
	operation := getOperation(effect.Links.Operation.Href)

	// Check whether sender and receiver are the same in account_credited and account_debited
	// operations. If so, ignore this effect.
	// This prevents issues from unchecked transaction to oneself, see:
	// https://stellar.stackexchange.com/questions/1036/why-are-overspending-transactions-added-to-the-ledger
	if operation.To == operation.From && (effect.Asset.Type == "account_credited" || effect.Asset.Type == "account_debited") {
		log.Printf("[WARN] Found effect where sender and receiver are the same account, ignoring.")
		log.Printf("       Debug: %v", effect)
		return
	}

	// Just input the fields we're requiring for now, can be replayed anytime form the chain later.
	_, err = db.Exec(`INSERT INTO effects(effect_id, operation, paging_token, account, amount, type, starting_balance, balance, balance_limit, asset_type, asset_issuer, asset_code, created_at)
		VALUES($1, $2, $3, $4, $5::REAL, $6, $7, $8, $9, $10, $11, $12, $13)`,
		effect.ID, effect.Links.Operation.Href, effect.PT, effect.Account, effect.Amount, effect.Type, effect.StartingBalance, effect.Balance.Balance, effect.Balance.Limit,
		effect.Asset.Type, effect.Asset.Issuer, effect.Asset.Code, operation.CreatedAt)
	if err != nil {
		return
	}

	// Store amount_transfered and amount_issued upon insert in a different table
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

func TotalAmount(db *sqlx.DB, filter Filter) (amount float64) {
	filter.Defaults()
	if filter.Type == "" {
		log.Printf("[ERROR] effect.TotalAmount(): No type given.")
		return
	}

	err := db.Get(&amount, `SELECT SUM(amount) FROM effects WHERE type=$1 AND created_at BETWEEN $2::timestsamp AND $3::timestsamp`,
		filter.Type, filter.From.Unix(), filter.To.Unix())
	if err != nil {
		log.Print(err)
	}
	return
}

// Total assets issued
func TotalIssued(db *sqlx.DB, issuer string, filter Filter) (amount float64) {
	filter.Defaults()
	err := db.Get(&amount, `SELECT SUM(amount) FROM effects WHERE type='account_debited' AND account=$1 AND created_at BETWEEN $2::timestsamp AND $3::timestamp`,
		issuer, filter.From.Unix(), filter.To.Unix())
	if err != nil {
		log.Print(err)
	}
	return
}

func TotalCount(db *sqlx.DB, filter Filter) (count int) {
	filter.Defaults()
	if filter.Type == "" {
		log.Printf("[ERROR] effect.TotalCount(): No type given.")
		return
	}

	err := db.Get(&count, `SELECT COUNT(*) FROM effects WHERE type=$1 AND created_at BETWEEN $2::timestamp AND $3::timestamp`,
		filter.Type, filter.From.Unix(), filter.To.Unix())
	if err != nil {
		log.Printf("[ERROR] effect.TotalCount(): %s", err)
	}
	return
}

func AccountCount(db *sqlx.DB, filter Filter) (count int) {
	filter.Defaults()
	err := db.Get(&count, `SELECT COUNT(DISTINCT account) FROM effects WHERE created_at BETWEEN $1::timestamp AND $2::timestamp`,
		filter.From.Unix(), filter.To.Unix())
	if err != nil {
		log.Printf("[ERROR] effect.AccountCount(): %s", err)
	}
	return
}

func ItemCount(db *sqlx.DB, filter Filter) (count int) {
	filter.Defaults()
	err := db.Get(&count, `SELECT COUNT(*) FROM effects WHERE created_at BETWEEN $1::timestamp AND $2::timestamp`,
		filter.From.Unix(), filter.To.Unix())
	if err != nil {
		log.Printf("[ERROR] effect.ItemCount(): %s", err)
	}
	return
}

func Get(db *sqlx.DB, filter Filter) (effects []Effect, err error) {
	filter.Defaults()
	err = db.Select(&effects, `SELECT * FROM effects WHERE created_at BETWEEN $1::timestsamp AND $2::timestsamp`,
		filter.From.Unix(), filter.To.Unix())
	if err == sql.ErrNoRows {
		log.Printf("[ERROR] effect.Get(): %s", err)
	}
	return
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
