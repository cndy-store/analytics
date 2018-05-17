package main

import (
	"github.com/cndy-store/analytics/controllers/effects"
	"github.com/cndy-store/analytics/controllers/history"
	"github.com/cndy-store/analytics/controllers/stats"
	"github.com/cndy-store/analytics/models/cursor"
	"github.com/cndy-store/analytics/models/effect"
	"github.com/cndy-store/analytics/utils/cndy"
	"github.com/cndy-store/analytics/utils/sql"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/stellar/go/clients/horizon"
	"golang.org/x/net/context"
	"log"
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
	history.Init(db, router)

	router.Run(":3144")
}
