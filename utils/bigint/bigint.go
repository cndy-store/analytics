package bigint

import (
	"github.com/stellar/go/amount"
	"math/big"
)

func Parse(s string) (*int64, error) {
	// Prase empty strings as nil
	if s == "" {
		zero := int64(0)
		return &zero, nil
	}

	p, err := amount.ParseInt64(s)
	return &p, err
}

// ToString returns an "amount string" from the provided raw int64 value `v`.
// Taken from: github.com/stellar/go/amount/main.go
func ToString(v int64) string {
	One := int64(10000000)
	bigOne := big.NewRat(One, 1)
	r := big.NewRat(v, 1)
	r.Quo(r, bigOne)
	return r.FloatString(7)
}
