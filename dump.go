package main

import (
	"encoding/gob"
	"fmt"
	"github.com/stellar/go/clients/horizon"
	"os"
	"strconv"
	"time"
)

type Collection struct {
	Items  []Item
	Cursor horizon.Cursor
}

type Item struct {
	Effect          horizon.Effect
	LedgerCloseTime time.Time
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

func main() {
	// Load collection upon startup
	if col.Load() {
		fmt.Printf("%d transactions of %d accounts loaded. Latest cursor: %s\n\n",
			col.ItemCount(), col.AccountCount(), col.Cursor)
	} else {
		panic("Couldn't load data")
	}

	for _, item := range col.Items {
		fmt.Printf("--+--[ %s ]\n", item.Effect.Asset.Code)
		fmt.Printf("  |\n")
		fmt.Printf("  +->  Type:    %s\n", item.Effect.Type)
		fmt.Printf("  +->  Account: %s\n", item.Effect.Account)
		fmt.Printf("  +->  Amount:  %s\n", item.Effect.Amount)
		fmt.Printf("  +->  Time:    %s\n\n", item.LedgerCloseTime)
	}
}
