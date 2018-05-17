package filter

import (
	"errors"
	"github.com/gin-gonic/gin"
	"time"
)

func Parse(c *gin.Context) (from *time.Time, to *time.Time, err error) {
	if query := c.Query("from"); query != "" {
		t, e := time.Parse(time.RFC3339, query)
		if e != nil {
			err = errors.New("Invalid date in 'from' parameter.")
			return
		}
		from = &t
	}

	if query := c.Query("to"); query != "" {
		t, e := time.Parse(time.RFC3339, query)
		if e != nil {
			err = errors.New("Invalid date in 'to' parameter.")
			return
		}
		to = &t
	}

	return
}
