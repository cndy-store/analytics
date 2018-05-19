package effects

import (
	"github.com/cndy-store/analytics/models/effect"
	"github.com/cndy-store/analytics/utils/filter"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func Init(db interface{}, router *gin.Engine) {
	// GET /effects[?from=XXX&to=XXX]
	router.GET("/effects", func(c *gin.Context) {
		from, to, err := filter.Parse(c)
		if err != nil {
			log.Printf("[ERROR] Couldn't parse URL parameters: %s", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		effects, err := effect.Get(db, effect.Filter{From: from, To: to})
		if err != nil {
			log.Printf("[ERROR] Couldn't get effect from database: %s", err)
			c.String(http.StatusInternalServerError, "")
			return
		}

		// Convert int64 fields to strings
		for i, _ := range effects {
			effects[i].Convert()
		}

		c.JSON(http.StatusOK, gin.H{
			"effects": effects,
		})
		return
	})
}
