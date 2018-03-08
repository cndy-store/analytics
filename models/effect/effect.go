package effect

import (
	"database/sql"
	"encoding/json"
	"github.com/jmoiron/sqlx"
	"github.com/stellar/go/clients/horizon"
	"log"
	"net/http"
	"time"
)

type Effect struct {
	Id              *uint32 `db:"id"               json:"-"`
	EffectId        *string `db:"effect_id"        json:"id,omitempty"`
	Operation       *string `db:"operation"        json:"operation,omitempty"`
	Succeeds        *string `db:"succeeds"         json:"succeeds,omitempty"` // Currently not used
	Precedes        *string `db:"precedes"         json:"precedes,omitempty"` // Currently not used
	PagingToken     *string `db:"paging_token"     json:"paging_token,omitempty"`
	Account         *string `db:"account"          json:"account,omitempty"`
	Amount          *string `db:"amount"           json:"amount,omitempty"`
	Type            *string `db:"type"             json:"type,omitempty"`
	TypeI           *int32  `db:"type_i"           json:"type_i,omitempty"` // Currently not used
	StartingBalance *string `db:"starting_balance" json:"starting_balance,omitempty"`

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

func New(db *sqlx.DB, effect horizon.Effect) (err error) {
	// Just input the fields we're requiring for now, can be replayed anytime form the chain later.
	_, err = db.Exec(`INSERT INTO effects(effect_id, operation, paging_token, account, amount, type, starting_balance, balance, balance_limit, asset_type, asset_issuer, asset_code, created_at)
	                  VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`,
		effect.ID, effect.Links.Operation.Href, effect.PT, effect.Account, effect.Amount, effect.Type, effect.StartingBalance, effect.Balance.Balance, effect.Balance.Limit,
		effect.Asset.Type, effect.Asset.Issuer, effect.Asset.Code,
		getOperationTime(effect.Links.Operation.Href))
	return
}

func TotalAmount(db *sqlx.DB, t string) (amount float64) {
	err := db.Get(&amount, `SELECT SUM(amount) FROM effects WHERE type=$1`, t)
	if err != nil {
		log.Print(err)
	}
	return
}

func TotalCount(db *sqlx.DB, t string) (count int) {
	err := db.Get(&count, `SELECT COUNT(amount) FROM effects WHERE type=$1`, t)
	if err != nil {
		log.Print(err)
	}
	return
}

func AccountCount(db *sqlx.DB) (count int) {
	err := db.Get(&count, `SELECT COUNT(DISTINCT account) FROM effects`)
	if err != nil {
		log.Print(err)
	}
	return
}

func ItemCount(db *sqlx.DB) (count int) {
	err := db.Get(&count, `SELECT COUNT(*) FROM effects`)
	if err != nil {
		log.Print(err)
	}
	return
}

// GetAllForEvent returns all available contribution_types, marking those active for a specific event
func GetAll(db *sqlx.DB) (effects []Effect, err error) {
	err = db.Select(&effects, `SELECT * FROM effects`)
	if err == sql.ErrNoRows {
		log.Print(err)
	}
	return
}

func getOperationTime(url string) (t time.Time) {
	var h = &http.Client{Timeout: 2 * time.Second}

	r, err := h.Get(url)
	if err != nil {
		log.Printf("GET Error: %s", err)
		return
	}
	defer r.Body.Close()

	type operation struct {
		CreatedAt time.Time `json:"created_at"`
	}
	op := operation{}

	err = json.NewDecoder(r.Body).Decode(&op)
	if err != nil {
		log.Printf("Couldn't decode body: %s", err)
		return
	}

	t = op.CreatedAt
	return
}
