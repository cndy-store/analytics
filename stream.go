package main

import (
	"encoding/gob"
	"github.com/stellar/go/clients/horizon"
	"golang.org/x/net/context"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

// Our asset code to watch and the cursor to when the asset was first introduced
const ASSET_CODE = "CNDY"
const ASSET_ISSUER = "GD7YB3R3TKUU3OHTE3DO5BIVBLQVFKYRHPW5Y6NHVSQVNNEOQ5I2RKLU"
const GENESIS_CURSOR = "33170762571452437-1"

var collection Collection

type Collection struct {
	Effects []horizon.Effect
	Cursor  horizon.Cursor
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

func (c *Collection) AccountCount() int {
	accounts := make(map[string]struct{})

	for _, e := range c.Effects {
		accounts[e.Account] = struct{}{}
	}

	return len(accounts)
}

func (c *Collection) TxCount() int {
	return len(c.Effects)
}

func init() {
	// Load collection upon startup
	if loadCollection() {
		log.Printf("%d transactions of %d accounts loaded. Resuming operation from cursor %s",
			collection.TxCount(), collection.AccountCount(), collection.Cursor)
	} else {
		collection.Cursor = horizon.Cursor(GENESIS_CURSOR)
		log.Printf("Retrieving data from blockchain beginning with cursor %s", collection.Cursor)
	}

	// Intercept signals
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	// Save collection on exit
	go func() {
		signal := <-signalChannel
		log.Printf("Received signal: %v\n", signal)
		log.Printf("Saving collection at cursor %s", collection.Cursor)
		saveCollection()
		os.Exit(0)
	}()
}

func main() {
	client := horizon.DefaultTestNetClient
	ctx := context.Background() // Stream indefinitly

	// Go subroutine to periodically print status
	go func() {
		ticker := time.NewTicker(time.Millisecond * 5000)
		for _ = range ticker.C {
			log.Printf("DEBUG: Total %s transferred: %f (Cursor: %s)",
				ASSET_CODE,
				collection.Total("account_credited"),
				collection.Cursor)
		}
	}()

	for {
		err := client.StreamEffects(ctx, &collection.Cursor, func(e horizon.Effect) {
			if e.Asset.Code == ASSET_CODE && e.Asset.Issuer == ASSET_ISSUER {
				log.Printf("--+--[ %s ]", e.Asset.Code)
				log.Printf("  |")
				log.Printf("  +->  Type:    %s", e.Type)
				log.Printf("  +->  Account: %s", e.Account)
				log.Printf("  +->  Amount:  %s\n\n", e.Amount)

				collection.Append(e)
			}

			// Update cursor pointer to resume operation in case connection gets lost
			collection.Cursor = horizon.Cursor(e.PT)
		})

		if err != nil {
			log.Print(err)
		}
	}
}

// Save current collection to file
func saveCollection() {
	f, err := os.Create("collection.save")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	enc := gob.NewEncoder(f)
	enc.Encode(collection)
	f.Close()
}

// Load saved collection from file
func loadCollection() bool {
	f, err := os.Open("collection.save")
	if err != nil {
		return false
	}
	defer f.Close()

	dec := gob.NewDecoder(f)
	dec.Decode(&collection)
	f.Close()
	return true
}
