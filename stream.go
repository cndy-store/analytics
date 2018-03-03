package main

import (
	"github.com/stellar/go/clients/horizon"
	"golang.org/x/net/context"
	"log"
)

const ASSET_CODE = "CNDY"

func main() {
	client := horizon.DefaultTestNetClient
	cursor := horizon.Cursor("now")
	ctx := context.Background() // Stream indefinitly

	for {
		err := client.StreamEffects(ctx, &cursor, func(e horizon.Effect) {
			if e.Asset.Code == ASSET_CODE {
				log.Printf("--+--[ %s ]\n", e.Asset.Code)
				log.Printf("  |\n")
				log.Printf("  +->  Account: %s\n", e.Account)
				log.Printf("  +->  Issuer:  %s\n", e.Asset.Issuer)
				log.Printf("  +->  Type:    %s\n", e.Type)
				log.Printf("  +->  Amount:  %s\n\n", e.Amount)
			}
		})

		if err != nil {
			log.Fatal(err)
		}
	}
}
