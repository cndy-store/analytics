package stats

import (
	"github.com/cndy-store/analytics/models/cursor"
	"github.com/cndy-store/analytics/models/effect"
	"github.com/cndy-store/analytics/utils/bigint"
	"github.com/cndy-store/analytics/utils/cndy"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Init(db interface{}, router *gin.Engine) {
	// GET /cndy/stats
	router.GET("/stats", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"asset_code":         cndy.AssetCode,
			"payments":           effect.TotalCount(db, effect.Filter{Type: "account_debited"}),
			"accounts_involved":  effect.AccountCount(db, effect.Filter{}),
			"amount_transferred": bigint.ToString(effect.TotalAmount(db, effect.Filter{Type: "account_credited"})),
			"trustlines_created": effect.TotalCount(db, effect.Filter{Type: "trustline_created"}),
			"amount_issued":      bigint.ToString(effect.TotalIssued(db, cndy.AssetIssuer, effect.Filter{})),
			"current_cursor":     cursor.Current,
		})
		return
	})
}
