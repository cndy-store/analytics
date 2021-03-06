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
		args, err := filter.Parse(c)
		if err != nil {
			log.Printf("[ERROR] Couldn't parse URL parameters: %s", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		assetStats, err := assetStat.Get(db, args)
		if err != nil {
			log.Printf("[ERROR] Couldn't get asset stats from database: %s", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Internal server error",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"stats":  assetStats,
		})
		return
	})

	router.GET("/stats/latest", func(c *gin.Context) {
		args, err := filter.Parse(c)
		if err != nil {
			log.Printf("[ERROR] Couldn't parse URL parameters: %s", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		latest, err := assetStat.Latest(db, args)
		if err != nil {
			log.Printf("[ERROR] Couldn't get asset stats from database: %s", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Internal server error",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"latest": latest,
		})
		return
	})

	router.GET("/stats/cursor", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":         "ok",
			"current_cursor": cursor.Current,
		})
		return
	})
}
