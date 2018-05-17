package stats

import (
	"github.com/cndy-store/analytics/models/cursor"
	"github.com/cndy-store/analytics/models/effect"
	"github.com/cndy-store/analytics/utils/bigint"
	"github.com/cndy-store/analytics/utils/cndy"
	"github.com/cndy-store/analytics/utils/filter"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
)

func Init(db *sqlx.DB, router *gin.Engine) {
	// GET /cndy/stats[?from=XXX&to=XXX]
	router.GET("/stats", func(c *gin.Context) {
		from, to, err := filter.Parse(c)
		if err != nil {
			log.Printf("[ERROR] Couldn't parse URL parameters: %s", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		currentCursor, err := cursor.GetLatest(db)
		if err != nil {
			log.Printf("[ERROR] Couldn't get latest cursor from database: %s", err)
		}

		c.JSON(http.StatusOK, gin.H{
			"asset_code":         cndy.AssetCode,
			"effect_count":       effect.ItemCount(db, effect.Filter{From: from, To: to}),
			"accounts_involved":  effect.AccountCount(db, effect.Filter{From: from, To: to}),
			"amount_transferred": bigint.ToString(effect.TotalAmount(db, effect.Filter{Type: "account_credited", From: from, To: to})),
			"trustlines_created": effect.TotalCount(db, effect.Filter{Type: "trustline_created", From: from, To: to}),
			"amount_issued":      bigint.ToString(effect.TotalIssued(db, cndy.AssetIssuer, effect.Filter{From: from, To: to})),
			"current_cursor":     currentCursor,
		})
		return
	})
}
