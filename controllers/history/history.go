package history

import (
	"github.com/cndy-store/analytics/models/asset_stat"
	"github.com/cndy-store/analytics/utils/filter"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
)

func Init(db *sqlx.DB, router *gin.Engine) {
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

		// Convert int64 fields to strings
		for i, _ := range assetStats {
			assetStats[i].Convert()
		}

		c.JSON(http.StatusOK, gin.H{
			"history": assetStats,
		})
		return
	})

}
