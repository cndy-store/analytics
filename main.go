package main

import (
	"github.com/cndy-store/analytics/controllers/effects"
	"github.com/cndy-store/analytics/controllers/stats"
	"github.com/cndy-store/analytics/models/asset_stat"
	"github.com/cndy-store/analytics/models/cursor"
	"github.com/cndy-store/analytics/models/effect"
	"github.com/cndy-store/analytics/utils/cndy"
	"github.com/cndy-store/analytics/utils/filter"
	"github.com/cndy-store/analytics/utils/sql"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/stellar/go/clients/horizon"
	"golang.org/x/net/context"
	"log"
	"net/http"
	"time"
)

func main() {
	db, err := sql.OpenAndMigrate(".")
	if err != nil {
		log.Fatal("[ERROR] Couldn't open database: ", err)
	}

	// Start API in go subroutine
	go api(db)

	client := horizon.DefaultTestNetClient
	ctx := context.Background() // Stream indefinitly

	for {
		currentCursor, err := cursor.GetLatest(db)
		if err != nil {
			log.Printf("[ERROR] Couldn't get latest cursor from database: %s", err)
			log.Printf("        Re-trying in 5 seconds...")
			time.Sleep(5 * time.Second)
			continue
		}

		client.StreamEffects(ctx, &currentCursor, func(e horizon.Effect) {
			if e.Asset.Code == cndy.AssetCode && e.Asset.Issuer == cndy.AssetIssuer {
				err = effect.New(db, e)
				if err != nil {
					log.Printf("[ERROR] Couldn't save effect to database: %s", err)
				}
			}

			err = cursor.New(db, e.PT)
			if err != nil {
				log.Printf("[ERROR] Couldn't save cursor to database: %s", err)
			}
		})
	}
}

func api(db *sqlx.DB) {
	router := gin.Default()
	router.Use(cors.Default()) // Allow all origins
	stats.Init(db, router)
	effects.Init(db, router)

	// GET /history[?from=XXX&to=XXX]
	router.GET("/history", func(c *gin.Context) {
		from, to, err := filter.Parse(c)
		if err != nil {
			log.Printf("[ERROR] Couldn't parse URL parameters: %s", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		assetStats, err := assetStat.Get(db, assetStat.Filter{From: from, To: to})
		if err != nil {
			log.Printf("[ERROR] Couldn't get asset stats from database: %s", err)
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
