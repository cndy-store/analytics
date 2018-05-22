package stats

import (
	"bytes"
	"fmt"
	"github.com/cndy-store/analytics/models/cursor"
	"github.com/cndy-store/analytics/utils/cndy"
	"github.com/cndy-store/analytics/utils/sql"
	"github.com/cndy-store/analytics/utils/test"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type HttpTest struct {
	method       string
	url          string
	body         string
	statusCode   int
	bodyContains []string
}

func TestStats(t *testing.T) {
	db, err := sql.OpenAndMigrate("../..")
	if err != nil {
		t.Error(err)
	}

	tx, err := db.Beginx()
	if err != nil {
		t.Error(err)
	}
	defer tx.Rollback()

	// Insert test data
	test.InsertEffects(tx)

	var tests = []HttpTest{
		{
			"GET",
			"/stats",
			"",
			http.StatusOK,
			[]string{
				fmt.Sprintf(`"asset_code":"%s"`, cndy.AssetCode),
				`"effect_count":8`,
				`"accounts_involved":3`,
				`"amount_issued":"1100.0000000"`,
				`"trustlines_created":2`,
				`"amount_transferred":"1115.0000000"`,
				fmt.Sprintf(`"current_cursor":"%s"`, cndy.GenesisCursor),
			},
		},

		// Filter{From}
		{
			"GET",
			fmt.Sprintf("/stats?from=%s", test.Effects[4].CreatedAt.Format(time.RFC3339)),
			"",
			http.StatusOK,
			[]string{
				fmt.Sprintf(`"asset_code":"%s"`, cndy.AssetCode),
				`"effect_count":4`,
				`"accounts_involved":3`,
				`"amount_issued":"100.0000000"`,
				`"trustlines_created":0`,
				`"amount_transferred":"115.0000000"`,
				fmt.Sprintf(`"current_cursor":"%s"`, cndy.GenesisCursor),
			},
		},

		// Filter{To}
		{
			"GET",
			fmt.Sprintf("/stats?to=%s", test.Effects[4].CreatedAt.Format(time.RFC3339)),
			"",
			http.StatusOK,
			[]string{
				fmt.Sprintf(`"asset_code":"%s"`, cndy.AssetCode),
				`"effect_count":5`,
				`"accounts_involved":3`,
				`"amount_issued":"1000.0000000"`,
				`"trustlines_created":2`,
				`"amount_transferred":"1015.0000000"`,
				fmt.Sprintf(`"current_cursor":"%s"`, cndy.GenesisCursor),
			},
		},

		// Filter{From, To}
		{
			"GET",
			fmt.Sprintf("/stats?from=%s&to=%s", test.Effects[1].CreatedAt.Format(time.RFC3339), test.Effects[3].CreatedAt.Format(time.RFC3339)),
			"",
			http.StatusOK,
			[]string{
				fmt.Sprintf(`"asset_code":"%s"`, cndy.AssetCode),
				`"effect_count":3`,
				`"accounts_involved":3`,
				`"amount_issued":"1000.0000000"`,
				`"trustlines_created":1`,
				`"amount_transferred":"1000.0000000"`,
				fmt.Sprintf(`"current_cursor":"%s"`, cndy.GenesisCursor),
			},
		},
	}

	router := gin.Default()
	err = cursor.LoadLatest(db)
	if err != nil {
		t.Errorf("Couldn't get latest cursor from database: %s", err)
	}
	Init(tx, router)

	for _, test := range tests {
		body := bytes.NewBufferString(test.body)
		req, _ := http.NewRequest(test.method, test.url, body)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		if resp.Code != test.statusCode {
			t.Errorf("Expected code %v, got %v, for %+v", test.statusCode, resp.Code, test)
		}

		if len(test.bodyContains) > 0 {
			for _, s := range test.bodyContains {
				if !strings.Contains(resp.Body.String(), s) {
					t.Errorf("Body did not contain '%s' in '%s'", s, resp.Body.String())
				}
			}
		}
	}
}
