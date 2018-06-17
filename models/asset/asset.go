package asset

import (
	"github.com/cndy-store/analytics/utils/sql"
	"time"
)

type Asset struct {
	Type      *string    `db:"type"       json:"type,omitempty"`
	Code      *string    `db:"code"       json:"code,omitempty"       binding="required"`
	Issuer    *string    `db:"issuer"     json:"issuer,omitempty      binding="required"`
	CreatedAt *time.Time `db:"created_at" json:"created_at,omitempty"`
}

// Variable to hold all registered assets in memory
var Registered []Asset

func New(db sql.Database, asset Asset) (ret Asset, err error) {
	err = db.Get(&ret, `INSERT INTO assets(type, code, issuer) VALUES($1, $2, $3) RETURNING *`,
		asset.Type, asset.Code, asset.Issuer)
	if err != nil {
		return
	}

	err = UpdateRegistered(db)
	return
}

func Exists(db sql.Database, asset Asset) (yes bool, err error) {
	return sql.Exists(db, `SELECT 1 FROM assets WHERE code=$1 AND issuer=$2`, asset.Code, asset.Issuer)
}

func Get(db sql.Database) (assets []Asset, err error) {
	err = db.Select(&assets, `SELECT * FROM assets`)
	if err == sql.ErrNoRows {
		err = nil
		return
	}
	if err != nil {
		return
	}

	return
}

func UpdateRegistered(db sql.Database) (err error) {
	Registered, err = Get(db)
	return
}
