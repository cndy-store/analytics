package main

import (
	"encoding/gob"
	"encoding/json"
	"github.com/gin-gonic/gin"
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

// Our asset code to watch and the cursor to when the asset was first introduced
const ASSET_CODE = "CNDY"
const ASSET_ISSUER = "GD7YB3R3TKUU3OHTE3DO5BIVBLQVFKYRHPW5Y6NHVSQVNNEOQ5I2RKLU"
const GENESIS_CURSOR = "33170762571452437-1"

type Collection struct {
	Items  []Item         `json:"transactions"`
	Cursor horizon.Cursor `json:"cursor"`
}

type Item struct {
	Effect          horizon.Effect `json:"effect"`
	LedgerCloseTime time.Time      `json:"created_at"`
}

// Save current collection to file
func (c *Collection) Save() {
	f, err := os.Create("collection.save")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	enc := gob.NewEncoder(f)
	enc.Encode(c)
	f.Close()
}

// Load saved collection from file
func (c *Collection) Load() bool {
	f, err := os.Open("collection.save")
	if err != nil {
		return false
	}
	defer f.Close()

	dec := gob.NewDecoder(f)
	dec.Decode(&c)
	f.Close()
	return true
}

// Aggregate amount of all entries with type t
func (c *Collection) TotalAmount(t string) (total float64) {
	for _, i := range c.Items {
		if i.Effect.Type == t {
			amount, err := strconv.ParseFloat(i.Effect.Amount, 64)
			if err == nil {
				total += amount
			}
		}
	}
	return
}

func (c *Collection) TotalCount(t string) (count int) {
	for _, i := range c.Items {
		if i.Effect.Type == t {
			count += 1
		}
	}
	return
}

func (c *Collection) Append(effect Item) {
	c.Items = append(c.Items, effect)
}

func (c *Collection) AccountCount() int {
	accounts := make(map[string]struct{})

	for _, i := range c.Items {
		accounts[i.Effect.Account] = struct{}{}
	}

	return len(accounts)
}

func (c *Collection) ItemCount() int {
	return len(c.Items)
}

var col Collection

func init() {

	// Load collection upon startup
	if col.Load() {
		log.Printf("%d transactions of %d accounts loaded. Resuming operation from cursor %s",
			col.ItemCount(), col.AccountCount(), col.Cursor)
	} else {
		col.Cursor = horizon.Cursor(GENESIS_CURSOR)
		log.Printf("Retrieving data from blockchain beginning with cursor %s", col.Cursor)
	}

	// Start routine for API
	go api()

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
		log.Printf("Saving collection at cursor %s", col.Cursor)
		col.Save()
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
			log.Printf("DEBUG: Total %s transferred: %f", ASSET_CODE, col.TotalAmount("account_credited"))
			log.Printf("       Total trustlines created: %d", col.TotalCount("trustline_created"))
			log.Printf("       Cursor: %s", col.Cursor)
		}
	}()

	for {
		err := client.StreamEffects(ctx, &col.Cursor, func(e horizon.Effect) {
			if e.Asset.Code == ASSET_CODE && e.Asset.Issuer == ASSET_ISSUER {
				log.Printf("--+--[ %s ]", e.Asset.Code)
				log.Printf("  |")
				log.Printf("  +->  Type:    %s", e.Type)
				log.Printf("  +->  Account: %s", e.Account)
				log.Printf("  +->  Amount:  %s\n\n", e.Amount)

				col.Append(Item{
					Effect:          e,
					LedgerCloseTime: getOperationTime(e.Links.Operation.Href),
				})
			}

			// Save cursor position to resume operation in case connection drops
			col.Cursor = horizon.Cursor(e.PT)
		})

		if err != nil {
			log.Printf("INFO: %s", err)
		}
	}
}

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

func getOperationTime(url string) (t time.Time) {
	var h = &http.Client{Timeout: 2 * time.Second}

	log.Printf("DEBUG: GET %s", url)

	r, err := h.Get(url)
	if err != nil {
		log.Printf("GET Error: %s", err)
		return
	}
	defer r.Body.Close()

	operation := Operation{}
	err = json.NewDecoder(r.Body).Decode(&operation)
	if err != nil {
		log.Printf("Couldn't decode body: %s", err)
		return
	}

	t = operation.LedgerCloseTime
	log.Printf("       %+v", t)
	return
}

func api() {
	router := gin.Default()

	router.GET("/api/cndy", func(c *gin.Context) {
		// if err != nil {
		// 	log.Printf("ERROR: %s", err)
		// 	c.String(http.StatusInternalServerError, "")
		// 	return
		// }

		c.JSON(http.StatusOK, gin.H{
			"tx_count":                 col.ItemCount(),
			"accounts":                 col.AccountCount(),
			"total_amount_transferred": col.TotalAmount("account_credited"),
			"trustlines_created":       col.TotalCount("trustline_created"),
			"current_cursor":           col.Cursor,
			// "transactions":       col.Items,
		})
		return
	})

	router.Run(":8080")
}
