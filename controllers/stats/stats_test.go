package stats

import (
	"bytes"
	"fmt"
	"github.com/cndy-store/analytics/models/cursor"
	"github.com/cndy-store/analytics/utils/bigint"
	"github.com/cndy-store/analytics/utils/cndy"
	"github.com/cndy-store/analytics/utils/sql"
	"github.com/cndy-store/analytics/utils/test"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
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
	test.InsertTestData(tx)

	latestEffect := test.Effects[len(test.Effects)-1]
	var tests = []HttpTest{
		{
			"GET",
			"/stats",
			"",
			http.StatusOK,
			[]string{
				fmt.Sprintf(`"total_amount":"%s"`, bigint.ToString(latestEffect.TotalAmount)),
				fmt.Sprintf(`"num_accounts":%d`, latestEffect.NumAccounts),
				fmt.Sprintf(`"payments":%d`, latestEffect.Payments),
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
