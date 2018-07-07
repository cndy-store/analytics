package filter

import (
	"errors"
	"github.com/cndy-store/analytics/utils/cndy"
	"github.com/gin-gonic/gin"
	"time"
)

type Filter struct {
	From        *time.Time
	To          *time.Time
	AssetCode   string
	AssetIssuer string
}

func Parse(c *gin.Context) (filter Filter, err error) {
	if query := c.Query("from"); query != "" {
		t, e := time.Parse(time.RFC3339, query)
		if e != nil {
			err = errors.New("Invalid date in 'from' parameter.")
			return
		}
		filter.From = &t
	}

	if query := c.Query("to"); query != "" {
		t, e := time.Parse(time.RFC3339, query)
		if e != nil {
			err = errors.New("Invalid date in 'to' parameter.")
			return
		}
		filter.To = &t
	}

	if query := c.Query("asset_code"); query != "" {
		filter.AssetCode = query
	} else {
		err = errors.New("Missing 'asset_code' parameter")
		return
	}

	if query := c.Query("asset_issuer"); query != "" {
		filter.AssetIssuer = query
	} else {
		err = errors.New("Missing 'asset_issuer' parameter")
		return
	}

	return
}

func (f *Filter) Defaults() {
	if f.From == nil {
		t := time.Unix(0, 0)
		f.From = &t
	}

	if f.To == nil {
		t := time.Now()
		f.To = &t
	}
}

// Returns a filter object with pre-filled AssetIssuer and AssetCode
// for the CNDY coin used in testing
func NewCNDYFilter(from *time.Time, to *time.Time) Filter {
	if from == nil {
		t := time.Unix(0, 0)
		from = &t
	}

	if to == nil {
		t := time.Now()
		to = &t
	}

	return Filter{
		From:        from,
		To:          to,
		AssetCode:   cndy.AssetCode,
		AssetIssuer: cndy.AssetIssuer,
	}
}
