package assets

import (
	"fmt"
	"github.com/cndy-store/analytics/models/asset"
	"github.com/cndy-store/analytics/utils/sql"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func Init(db sql.Database, router *gin.Engine) {
	// TODO: Disable asset creation until issues with verification/ authentication are solved
	/*
		router.POST("/assets", func(c *gin.Context) {
			// Read JSON body and parse it into asset struct
			body := asset.Asset{}
			err := c.BindJSON(&body)
			if err != nil {
				jsonErrorMsg := fmt.Sprintf("Invalid JSON body: %s", err)
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  "error",
					"message": jsonErrorMsg,
				})
				return
			}

			exists, err := asset.Exists(db, body)
			if err != nil {
				log.Printf("[ERROR] POST /assets: %s", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"status":  "error",
					"message": "Internal server error",
				})
				return
			}
			if exists {
				c.JSON(http.StatusConflict, gin.H{
					"status":  "error",
					"message": "Asset already exists",
				})
				return
			}

			newAsset, err := asset.New(db, body)
			if err != nil {
				log.Printf("[ERROR] POST /assets: %s", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"status":  "error",
					"message": "Internal server error",
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"status": "ok",
				"asset":  newAsset,
			})
		})
	*/

	router.GET("/assets", func(c *gin.Context) {
		assets, err := asset.Get(db)
		if err != nil {
			log.Printf("[ERROR] Couldn't get assets from database: %s", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Internal server error",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"assets": assets,
		})
		return
	})
}
