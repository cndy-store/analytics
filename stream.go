package main

import (
	"github.com/stellar/go/clients/horizon"
	"golang.org/x/net/context"
	"log"
	"strconv"
)

// TODO:
// - Store effect entries in database or binary file
// - Store PT and use it as cursor, so we can catch up from last time

const ASSET_CODE = "CNDY"
const GENESIS_CURSOR = "33170762571452437-1"

type Collection struct {
	Effects []horizon.Effect
}

// Aggregate amount of all entries with type t
func (c *Collection) Total(t string) (total float64) {
	for _, e := range c.Effects {
		if e.Type == t {
			amount, err := strconv.ParseFloat(e.Amount, 64)
			if err == nil {
				total += amount
			}
		}
	}
	return
}

func (c *Collection) Append(effect horizon.Effect) {
	c.Effects = append(c.Effects, effect)
}

func main() {
	client := horizon.DefaultTestNetClient
	cursor := horizon.Cursor(GENESIS_CURSOR)
	ctx := context.Background() // Stream indefinitly

	var collection Collection

	for {
		err := client.StreamEffects(ctx, &cursor, func(e horizon.Effect) {
			if e.Asset.Code == ASSET_CODE {
				log.Printf("--+--[ %s ]", e.Asset.Code)
				log.Printf("  |")
				log.Printf("  +->  Account: %s", e.Account)
				log.Printf("  +->  Issuer:  %s", e.Asset.Issuer)
				log.Printf("  +->  Type:    %s", e.Type)
				log.Printf("  +->  Amount:  %s\n\n", e.Amount)

				collection.Append(e)
			}

			log.Printf("DEBUG: Total %s transferred: %f", ASSET_CODE, collection.Total("account_credited"))
			log.Printf("DEBUG: Current cursor: %s", e.PT)
		})

		if err != nil {
			log.Print(err)
		}
	}
}
