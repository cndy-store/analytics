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
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	db, err := sql.OpenAndMigrate(".")
	if err != nil {
		log.Fatal("[ERROR] Couldn't open database: ", err)
	}

	// Start API in go subroutine
	go api(db)

	// Intercept signals
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	// Save cursor on exit
	go func() {
		signal := <-signalChannel
		log.Printf("Received signal: %v\n", signal)
		log.Printf("Saving cursor to database: %s\n", cursor.Current)

		err = cursor.Save(db)
		if err != nil {
			log.Printf("[ERROR] Couldn't save cursor to database: %s", err)
		}
		os.Exit(0)
	}()

	// Also, save cursor every 5 minutes
	go func() {
		ticker := time.NewTicker(time.Minute * 5)
		for _ = range ticker.C {
			log.Printf("Saving cursor to database: %s\n", cursor.Current)

			err = cursor.Save(db)
			if err != nil {
				log.Printf("[ERROR] Couldn't save cursor to database: %s", err)
			}
		}
	}()

	client := horizon.DefaultTestNetClient
	ctx := context.Background() // Stream indefinitly

	// Load latest cursor from database
	err = cursor.LoadLatest(db)
	if err != nil {
		log.Printf("[ERROR] Couldn't get latest cursor from database: %s", err)
		os.Exit(1)
	}

	for {
		client.StreamEffects(ctx, &cursor.Current, func(e horizon.Effect) {
			if e.Asset.Code == cndy.AssetCode && e.Asset.Issuer == cndy.AssetIssuer {
				err = effect.New(db, e)
				if err != nil {
					log.Printf("[ERROR] Couldn't save effect to database: %s", err)
				}

				// Make sure to also safe the current cursor, so database is consistent
				err = cursor.Save(db)
				if err != nil {
					log.Printf("[ERROR] Couldn't save cursor to database: %s", err)
				}
			}

			cursor.Update(horizon.Cursor(e.PT))
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
