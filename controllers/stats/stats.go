package stats

import (
	"github.com/cndy-store/analytics/models/asset_stat"
	"github.com/cndy-store/analytics/models/cursor"
	"github.com/cndy-store/analytics/utils/filter"
	"github.com/cndy-store/analytics/utils/sql"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func Init(db sql.Database, router *gin.Engine) {
	// GET /stats[?from=XXX&to=XXX]
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

		assetStats, err := assetStat.Get(db, assetStat.Filter{From: from, To: to})
		if err != nil {
			log.Printf("[ERROR] Couldn't get asset stats from database: %s", err)
			c.String(http.StatusInternalServerError, "")
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"stats": assetStats,
		})
		return
	})

	router.GET("/stats/latest", func(c *gin.Context) {
		latest, err := assetStat.Latest(db)
		if err != nil {
			log.Printf("[ERROR] Couldn't get asset stats from database: %s", err)
			c.String(http.StatusInternalServerError, "")
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"latest": latest,
		})
		return
	})

	router.GET("/stats/cursor", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"current_cursor": cursor.Current,
		})
		return
	})
}
