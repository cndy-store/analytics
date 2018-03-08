package main

import (
	"github.com/chr4/cndy-analytics/models/cursor"
	"github.com/chr4/cndy-analytics/models/effect"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/mattes/migrate"
	"github.com/mattes/migrate/database/sqlite3"
	_ "github.com/mattes/migrate/source/file"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stellar/go/clients/horizon"
	"golang.org/x/net/context"
	"log"
	"net/http"
	"time"
)

// Our asset code to watch and the cursor to when the asset was first introduced
const ASSET_CODE = "CNDY"
const ASSET_ISSUER = "GD7YB3R3TKUU3OHTE3DO5BIVBLQVFKYRHPW5Y6NHVSQVNNEOQ5I2RKLU"

func main() {
	db, err := initdb("database.sqlite3")
	if err != nil {
		log.Fatal("Fatal error opening database: ", err)
	}

	go api(db)

	client := horizon.DefaultTestNetClient
	ctx := context.Background() // Stream indefinitly

	// Go subroutine to periodically print status
	go func() {
		ticker := time.NewTicker(time.Millisecond * 5000)
		for _ = range ticker.C {
			log.Printf("DEBUG: Total %s transferred: %f", ASSET_CODE, effect.TotalAmount(db, "account_credited"))
			log.Printf("       Total trustlines created: %d", effect.TotalCount(db, "trustline_created"))
			log.Printf("       Cursor: %s", cursor.GetLatest(db))
		}
	}()

	for {
		currentCursor := cursor.GetLatest(db)
		err := client.StreamEffects(ctx, &currentCursor, func(e horizon.Effect) {
			if e.Asset.Code == ASSET_CODE && e.Asset.Issuer == ASSET_ISSUER {
				log.Printf("--+--[ %s ]", e.Asset.Code)
				log.Printf("  |")
				log.Printf("  +->  Type:    %s", e.Type)
				log.Printf("  +->  Account: %s", e.Account)
				log.Printf("  +->  Amount:  %s\n\n", e.Amount)

				err = effect.New(db, e)
				if err != nil {
					log.Printf("Error saving to database: %s", err)
				}
			}

			err = cursor.New(db, e.PT)
			if err != nil {
				log.Printf("Error saving to database: %s", err)
			}
		})

		if err != nil {
			log.Printf("INFO: %s", err)
		}
	}
}

func api(db *sqlx.DB) {
	router := gin.Default()

	// TODO: From-until
	router.GET("/api/cndy/stats", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"asset_code":         ASSET_CODE,
			"tx_count":           effect.ItemCount(db),
			"accounts_involved":  effect.AccountCount(db),
			"amount_transferred": effect.TotalAmount(db, "account_credited"),
			"trustlines_created": effect.TotalCount(db, "trustline_created"),
			"current_cursor":     cursor.GetLatest(db),
		})
		return
	})

	router.GET("/api/cndy/transactions", func(c *gin.Context) {
		effects, err := effect.GetAll(db)
		if err != nil {
			log.Printf("ERROR: %s", err)
			c.String(http.StatusInternalServerError, "")
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"transactions": effects,
		})
		return
	})

	router.Run(":8080")
}

// Open database and run migrations
func initdb(uri string) (db *sqlx.DB, err error) {
	db, err = sqlx.Open("sqlite3", uri)
	if err != nil {
		return
	}

	driver, err := sqlite3.WithInstance(db.DB, &sqlite3.Config{})
	if err != nil {
		return
	}

	m, err := migrate.NewWithDatabaseInstance("file://db/migrations", "sqlite3", driver)
	if err != nil {
		return
	}

	err = m.Up()
	if err == migrate.ErrNoChange {
		err = nil
	}

	return
}
