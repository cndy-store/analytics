package main

import (
	"encoding/gob"
	"encoding/json"
	"github.com/stellar/go/clients/horizon"
	"golang.org/x/net/context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

// This struct is taken from the internal package github.com/stellar/go/services/horizon/internal/resource/operations"
type Operation struct {
	Links struct {
		Self        horizon.Link `json:"self"`
		Transaction horizon.Link `json:"transaction"`
		Effects     horizon.Link `json:"effects"`
		Succeeds    horizon.Link `json:"succeeds"`
		Precedes    horizon.Link `json:"precedes"`
	} `json:"_links"`

	ID              string    `json:"id"`
	PT              string    `json:"paging_token"`
	SourceAccount   string    `json:"source_account"`
	Type            string    `json:"type"`
	TypeI           int32     `json:"type_i"`
	LedgerCloseTime time.Time `json:"created_at"`
	TransactionHash string    `json:"transaction_hash"`
}

// Our asset code to watch and the cursor to when the asset was first introduced
const ASSET_CODE = "CNDY"
const ASSET_ISSUER = "GD7YB3R3TKUU3OHTE3DO5BIVBLQVFKYRHPW5Y6NHVSQVNNEOQ5I2RKLU"
const GENESIS_CURSOR = "33170762571452437-1"

var collection Collection

type EffectWithTime struct {
	Effect          horizon.Effect
	LedgerCloseTime time.Time
}

type Collection struct {
	Effects []EffectWithTime
	Cursor  horizon.Cursor
}

// Aggregate amount of all entries with type t
func (c *Collection) TotalAmount(t string) (total float64) {
	for _, e := range c.Effects {
		if e.Effect.Type == t {
			amount, err := strconv.ParseFloat(e.Effect.Amount, 64)
			if err == nil {
				total += amount
			}
		}
	}
	return
}

func (c *Collection) TotalCount(t string) (count int) {
	for _, e := range c.Effects {
		if e.Effect.Type == t {
			count += 1
		}
	}
	return
}

func (c *Collection) Append(effect EffectWithTime) {
	c.Effects = append(c.Effects, effect)
}

func (c *Collection) AccountCount() int {
	accounts := make(map[string]struct{})

	for _, e := range c.Effects {
		accounts[e.Effect.Account] = struct{}{}
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
			log.Printf("DEBUG: Total %s transferred: %f", ASSET_CODE, collection.TotalAmount("account_credited"))
			log.Printf("       Total trustlines created: %d", collection.TotalCount("trustline_created"))
			log.Printf("       Cursor: %s", collection.Cursor)
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

				// Try getting operation via GET request
				var myClient = &http.Client{Timeout: 2 * time.Second}

				r, err := myClient.Get(e.Links.Operation.Href)
				if err != nil {
					log.Printf("GET Error: %s", err)
				}
				defer r.Body.Close()

				operation := Operation{}
				err = json.NewDecoder(r.Body).Decode(&operation)
				if err != nil {
					log.Printf("Couldn't decode body: %s", err)
				}

				log.Printf("DEBUG: %+v", operation.LedgerCloseTime)

				newEffect := EffectWithTime{
					Effect:          e,
					LedgerCloseTime: operation.LedgerCloseTime,
				}

				collection.Append(newEffect)
			}

			// Save cursor position to resume operation in case connection gets lost
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
