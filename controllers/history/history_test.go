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
	expectedStats []test.AssetStat
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

	err = test.InsertAssetStats(tx)
	if err != nil {
		t.Error(err)
	}

	var tests = []HttpTest{
		{
			"GET",
			"/history",
			"",
			http.StatusOK,
			test.AssetStats,
		},

		// Filter{From}
		{
			"GET",
			fmt.Sprintf("/history?from=%s", test.AssetStats[2].CreatedAt.Format(time.RFC3339)),
			"",
			http.StatusOK,
			test.AssetStats[2:],
		},

		// Filter{To}
		{
			"GET",
			fmt.Sprintf("/history?to=%s", test.AssetStats[1].CreatedAt.Format(time.RFC3339)),
			"",
			http.StatusOK,
			test.AssetStats[:2],
		},

		// Filter{From, To}
		{
			"GET",
			fmt.Sprintf("/history?from=%s&to=%s", test.AssetStats[1].CreatedAt.Format(time.RFC3339), test.AssetStats[2].CreatedAt.Format(time.RFC3339)),
			"",
			http.StatusOK,
			test.AssetStats[1:3],
		},
	}

	router := gin.Default()
	Init(tx, router)

	for _, test := range tests {
		body := bytes.NewBufferString(test.body)
		req, _ := http.NewRequest(test.method, test.url, body)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		if resp.Code != test.statusCode {
			t.Errorf("Expected code %v, got %v, for %+v", test.statusCode, resp.Code, test)
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

		if len(history["history"]) != len(test.expectedStats) {
			t.Errorf("Expected %d JSON elements, got %d", len(test.expectedStats), len(history["history"]))
		}

		for _, e := range test.expectedStats {
			var s []string
			s = append(s, fmt.Sprintf(`"paging_token":"%s"`, e.PagingToken))
			s = append(s, fmt.Sprintf(`"total_amount":"%s"`, bigint.ToString(e.TotalAmount)))
			s = append(s, fmt.Sprintf(`"num_accounts":%d`, e.NumAccounts))
			s = append(s, fmt.Sprintf(`"num_effects":%d`, e.NumEffects))

			for _, contains := range s {
				if !strings.Contains(resp.Body.String(), contains) {
					t.Errorf("Body did not contain '%s' in '%s'", contains, resp.Body.String())
				}
			}
		}
	}
}
