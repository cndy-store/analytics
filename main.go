package main

import (
	"errors"
	"github.com/chr4/cndy-analytics/models/asset_stat"
	"github.com/chr4/cndy-analytics/models/cursor"
	"github.com/chr4/cndy-analytics/models/effect"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/sqlite3"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stellar/go/clients/horizon"
	"golang.org/x/net/context"
	"log"
	"net/http"
	"time"
)

// Our asset code to watch and the cursor to when the asset was first introduced
const ASSET_CODE = "CNDY"
const ASSET_ISSUER = "GCJKC2MI63KSQ6MLE6GBSXPDKTDAK43WR522ZYR3F34NPM7Z5UEPIZNX"

func main() {
	db, err := initdb("database.sqlite3")
	if err != nil {
		log.Fatal("Fatal error opening database: ", err)
	}

	// Start API in go subroutine
	go api(db)

	client := horizon.DefaultTestNetClient
	ctx := context.Background() // Stream indefinitly

	for {
		currentCursor := cursor.GetLatest(db)
		client.StreamEffects(ctx, &currentCursor, func(e horizon.Effect) {
			if e.Asset.Code == ASSET_CODE && e.Asset.Issuer == ASSET_ISSUER {
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
	}
}

func api(db *sqlx.DB) {
	router := gin.Default()
	router.Use(cors.Default()) // Allow all origins

	// GET /cndy/stats[?from=XXX&to=XXX]
	router.GET("/stats", func(c *gin.Context) {
		from, to, err := getFromAndTo(c)
		if err != nil {
			log.Printf("ERROR: %s", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"asset_code":         ASSET_CODE,
			"effect_count":       effect.ItemCount(db, effect.Filter{From: from, To: to}),
			"accounts_involved":  effect.AccountCount(db, effect.Filter{From: from, To: to}),
			"amount_transferred": effect.TotalAmount(db, effect.Filter{Type: "account_credited", From: from, To: to}),
			"trustlines_created": effect.TotalCount(db, effect.Filter{Type: "trustline_created", From: from, To: to}),
			"amount_issued":      effect.TotalIssued(db, ASSET_ISSUER, effect.Filter{From: from, To: to}),
			"current_cursor":     cursor.GetLatest(db),
		})
		return
	})

	// GET /effects[?from=XXX&to=XXX]
	router.GET("/effects", func(c *gin.Context) {
		from, to, err := getFromAndTo(c)
		if err != nil {
			log.Printf("ERROR: %s", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		effects, err := effect.Get(db, effect.Filter{From: from, To: to})
		if err != nil {
			log.Printf("ERROR: %s", err)
			c.String(http.StatusInternalServerError, "")
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"effects": effects,
		})
		return
	})

	// GET /history[?from=XXX&to=XXX]
	router.GET("/history", func(c *gin.Context) {
		from, to, err := getFromAndTo(c)
		if err != nil {
			log.Printf("ERROR: %s", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		assetStats, err := assetStat.Get(db, assetStat.Filter{From: from, To: to})
		if err != nil {
			log.Printf("ERROR: %s", err)
			c.String(http.StatusInternalServerError, "")
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"history": assetStats,
		})
		return
	})

	router.Run(":3144")
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

func getFromAndTo(c *gin.Context) (from *time.Time, to *time.Time, err error) {
	if query := c.Query("from"); query != "" {
		t, e := time.Parse(time.RFC3339, query)
		if e != nil {
			err = errors.New("Invalid date in 'from' parameter.")
			return
		}
		from = &t
	}

	if query := c.Query("to"); query != "" {
		t, e := time.Parse(time.RFC3339, query)
		if e != nil {
			err = errors.New("Invalid date in 'to' parameter.")
			return
		}
		to = &t
	}

	return
}
