package effects

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/cndy-store/analytics/models/effect"
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
	method        string
	url           string
	body          string
	statusCode    int
	expectedStats []test.Effect
}

func TestEffects(t *testing.T) {
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
			fmt.Sprintf("/effects?asset_code=%s&asset_issuer=%s", cndy.AssetCode, cndy.AssetIssuer),
			"",
			http.StatusOK,
			test.Effects,
		},

		// Filter{From}
		{
			"GET",
			fmt.Sprintf("/effects?asset_code=%s&asset_issuer=%s&from=%s", cndy.AssetCode, cndy.AssetIssuer, test.Effects[5].CreatedAt.Format(time.RFC3339)),
			"",
			http.StatusOK,
			test.Effects[5:],
		},

		// Filter{To}
		{
			"GET",
			fmt.Sprintf("/effects?asset_code=%s&asset_issuer=%s&to=%s", cndy.AssetCode, cndy.AssetIssuer, test.Effects[2].CreatedAt.Format(time.RFC3339)),
			"",
			http.StatusOK,
			test.Effects[:3],
		},

		// Filter{From, To}
		{
			"GET",
			fmt.Sprintf("/effects?asset_code=%s&asset_issuer=%s&from=%s&to=%s", cndy.AssetCode, cndy.AssetIssuer, test.Effects[3].CreatedAt.Format(time.RFC3339), test.Effects[4].CreatedAt.Format(time.RFC3339)),
			"",
			http.StatusOK,
			test.Effects[3:5],
		},

		// Invalid Filter{}
		{
			"GET",
			fmt.Sprintf("/effects?asset_code=%s&asset_issuer=%s&from=xxx", cndy.AssetCode, cndy.AssetIssuer),
			"",
			http.StatusBadRequest,
			nil,
		},

		// Missing asset_code and asset_issuer
		{
			"GET",
			"/effects",
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
			Status  string
			Effects []effect.Effect
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

		if len(res.Effects) != len(tt.expectedStats) {
			t.Errorf("%s %s: Expected %d JSON elements, got %d", tt.method, tt.url, len(tt.expectedStats), len(res.Effects))
		}

		for _, e := range tt.expectedStats {
			var s []string
			s = append(s, fmt.Sprintf(`"paging_token":"%s"`, e.PagingToken))
			s = append(s, fmt.Sprintf(`"account":"%s"`, e.Account))
			s = append(s, fmt.Sprintf(`"type":"%s"`, e.Type))

			if e.Amount == "" {
				s = append(s, `"amount":"0.0000000"`)
			} else {
				s = append(s, fmt.Sprintf(`"amount":"%s"`, e.Amount))
			}

			for _, contains := range s {
				if !strings.Contains(resp.Body.String(), contains) {
					t.Errorf("Body did not contain '%s' in '%s'", contains, resp.Body.String())
				}
			}
		}
	}
}
