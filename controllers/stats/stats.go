package stats

import (
	"github.com/cndy-store/analytics/models/asset_stat"
	"github.com/cndy-store/analytics/models/cursor"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func Init(db interface{}, router *gin.Engine) {
	// GET /stats
	router.GET("/stats", func(c *gin.Context) {
		latest, err := assetStat.Latest(db)
		if err != nil {
			log.Printf("[ERROR] Couldn't get asset stats from database: %s", err)
			c.String(http.StatusInternalServerError, "")
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"stats":          latest,
			"current_cursor": cursor.Current,
		})
		return
	})
}
