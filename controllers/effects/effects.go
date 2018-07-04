package effects

import (
	"github.com/cndy-store/analytics/models/effect"
	"github.com/cndy-store/analytics/utils/filter"
	"github.com/cndy-store/analytics/utils/sql"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func Init(db sql.Database, router *gin.Engine) {
	// GET /effects[?from=XXX&to=XXX]
	router.GET("/effects", func(c *gin.Context) {
		args, err := filter.Parse(c)
		if err != nil {
			log.Printf("[ERROR] Couldn't parse URL parameters: %s", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		effects, err := effect.Get(db, args)
		if err != nil {
			log.Printf("[ERROR] Couldn't get effect from database: %s", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Internal server error",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"effects": effects,
		})
		return
	})
}
