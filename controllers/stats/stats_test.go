package stats

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/cndy-store/analytics/models/asset_stat"
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
	"time"
)

type HttpTest struct {
	method       string
	url          string
	body         string
	statusCode   int
	bodyContains []string
}

type HttpTestWithEffects struct {
	method        string
	url           string
	body          string
	statusCode    int
	expectedStats []test.Effect
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

	err = test.InsertTestData(tx)
	if err != nil {
		t.Error(err)
	}

	var tests = []HttpTestWithEffects{
		{
			"GET",
			"/stats",
			"",
			http.StatusOK,
			test.Effects,
		},

		// Filter{From}
		{
			"GET",
			fmt.Sprintf("/stats?from=%s", test.Effects[4].CreatedAt.Format(time.RFC3339)),
			"",
			http.StatusOK,
			test.Effects[4:],
		},

		// Filter{To}
		{
			"GET",
			fmt.Sprintf("/stats?to=%s", test.Effects[2].CreatedAt.Format(time.RFC3339)),
			"",
			http.StatusOK,
			test.Effects[:3],
		},

		// Filter{From, To}
		{
			"GET",
			fmt.Sprintf("/stats?from=%s&to=%s", test.Effects[3].CreatedAt.Format(time.RFC3339), test.Effects[6].CreatedAt.Format(time.RFC3339)),
			"",
			http.StatusOK,
			test.Effects[3:7],
		},

		// Invalid Filter{}
		{
			"GET",
			"/stats?from=xxx",
			"",
			http.StatusBadRequest,
			nil,
		},
	}

	router := gin.Default()
	Init(tx, router)

	for _, tt := range tests {
		body := bytes.NewBufferString(tt.body)
		req, _ := http.NewRequest(tt.method, tt.url, body)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		if resp.Code != tt.statusCode {
			t.Errorf("Expected code %v, got %v, for %+v", tt.statusCode, resp.Code, tt)
		}

		type resJson struct {
			Status string
			Stats  []assetStat.AssetStat
		}

		if tt.statusCode == http.StatusOK {
			if !strings.Contains(resp.Body.String(), `"status":"ok"`) {
				t.Errorf("Body did not contain ok status message: %s", resp.Body.String())
			}
		} else {
			if !strings.Contains(resp.Body.String(), `"status":"error"`) {
				t.Errorf("Body did not contain error status message: %s", resp.Body.String())
			}

			// Skip to next test
			continue
		}

		res := resJson{}
		err := json.Unmarshal([]byte(resp.Body.String()), &res)
		if err != nil {
			t.Error(err)
		}

		if len(res.Stats) != len(tt.expectedStats) {
			t.Errorf("Expected %d JSON elements, got %d", len(tt.expectedStats), len(res.Stats))
		}

		for _, e := range tt.expectedStats {
			var s []string
			s = append(s, fmt.Sprintf(`"paging_token":"%s"`, e.PagingToken))
			s = append(s, fmt.Sprintf(`"issued":"%s"`, bigint.ToString(e.Issued)))
			s = append(s, fmt.Sprintf(`"transferred":"%s"`, bigint.ToString(e.Transferred)))
			s = append(s, fmt.Sprintf(`"accounts_with_trustline":%d`, e.AccountsWithTrustline))
			s = append(s, fmt.Sprintf(`"accounts_with_payments":%d`, e.AccountsWithPayments))
			s = append(s, fmt.Sprintf(`"payments":%d`, e.Payments))

			for _, contains := range s {
				if !strings.Contains(resp.Body.String(), contains) {
					t.Errorf("Body did not contain '%s' in '%s'", contains, resp.Body.String())
				}
			}
		}
	}
}

func TestLatestAndCursor(t *testing.T) {
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
			"/stats/latest",
			"",
			http.StatusOK,
			[]string{
				`"status":"ok"`,
				fmt.Sprintf(`"paging_token":"%s"`, latestEffect.PagingToken),
				fmt.Sprintf(`"issued":"%s"`, bigint.ToString(latestEffect.Issued)),
				fmt.Sprintf(`"transferred":"%s"`, bigint.ToString(latestEffect.Transferred)),
				fmt.Sprintf(`"accounts_with_trustline":%d`, latestEffect.AccountsWithTrustline),
				fmt.Sprintf(`"accounts_with_payments":%d`, latestEffect.AccountsWithPayments),
				fmt.Sprintf(`"payments":%d`, latestEffect.Payments),
			},
		},

		{
			"GET",
			"/stats/cursor",
			"",
			http.StatusOK,
			[]string{
				`"status":"ok"`,
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

		if test.statusCode == http.StatusOK {
			if !strings.Contains(resp.Body.String(), `"status":"ok"`) {
				t.Errorf("Body did not contain ok status message: %s", resp.Body.String())
			}
		} else {
			if !strings.Contains(resp.Body.String(), `"status":"error"`) {
				t.Errorf("Body did not contain error status message: %s", resp.Body.String())
			}

			// Skip to next test
			continue
		}

		if len(test.bodyContains) > 0 {
			for _, s := range test.bodyContains {
				if !strings.Contains(resp.Body.String(), s) {
					t.Errorf("Body did not contain '%s' in '%s'", s, resp.Body.String())
				}
			}
		}

		// Check whether JSON ID is hidden (regression test)
		if strings.Contains(resp.Body.String(), `"id":`) || strings.Contains(resp.Body.String(), `"Id":`) {
			t.Errorf("Body did contain JSON ID (should be excluded) in '%s'", resp.Body.String())
		}
	}
}
