package history

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/cndy-store/analytics/models/asset_stat"
	"github.com/cndy-store/analytics/utils/bigint"
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
	method        string
	url           string
	body          string
	statusCode    int
	expectedStats []test.Effect
}

func TestHistory(t *testing.T) {
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

	var tests = []HttpTest{
		{
			"GET",
			"/history",
			"",
			http.StatusOK,
			test.Effects,
		},

		// Filter{From}
		{
			"GET",
			fmt.Sprintf("/history?from=%s", test.Effects[4].CreatedAt.Format(time.RFC3339)),
			"",
			http.StatusOK,
			test.Effects[4:],
		},

		// Filter{To}
		{
			"GET",
			fmt.Sprintf("/history?to=%s", test.Effects[2].CreatedAt.Format(time.RFC3339)),
			"",
			http.StatusOK,
			test.Effects[:3],
		},

		// Filter{From, To}
		{
			"GET",
			fmt.Sprintf("/history?from=%s&to=%s", test.Effects[3].CreatedAt.Format(time.RFC3339), test.Effects[6].CreatedAt.Format(time.RFC3339)),
			"",
			http.StatusOK,
			test.Effects[3:7],
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

		history := make(map[string][]assetStat.AssetStat)
		err := json.Unmarshal([]byte(resp.Body.String()), &history)
		if err != nil {
			t.Error(err)
		}

		_, ok := history["history"]
		if !ok {
			t.Error(`Expected element "history" in JSON response`)
		}

		if len(history["history"]) != len(tt.expectedStats) {
			t.Errorf("Expected %d JSON elements, got %d", len(tt.expectedStats), len(history["history"]))
		}

		for _, e := range tt.expectedStats {
			var s []string
			s = append(s, fmt.Sprintf(`"paging_token":"%s"`, e.PagingToken))
			s = append(s, fmt.Sprintf(`"total_amount":"%s"`, bigint.ToString(e.TotalAmount)))
			s = append(s, fmt.Sprintf(`"num_accounts":%d`, e.NumAccounts))
			s = append(s, fmt.Sprintf(`"payments":%d`, e.Payments))

			for _, contains := range s {
				if !strings.Contains(resp.Body.String(), contains) {
					t.Errorf("Body did not contain '%s' in '%s'", contains, resp.Body.String())
				}
			}
		}
	}
}
