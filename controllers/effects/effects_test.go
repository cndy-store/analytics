package effects

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/cndy-store/analytics/models/effect"
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

	err = test.InsertEffects(tx)
	if err != nil {
		t.Error(err)
	}

	var tests = []HttpTest{
		{
			"GET",
			"/effects",
			"",
			http.StatusOK,
			test.Effects,
		},

		// Filter{From}
		{
			"GET",
			fmt.Sprintf("/effects?from=%s", test.Effects[5].CreatedAt.Format(time.RFC3339)),
			"",
			http.StatusOK,
			test.Effects[5:],
		},

		// Filter{To}
		{
			"GET",
			fmt.Sprintf("/effects?to=%s", test.Effects[2].CreatedAt.Format(time.RFC3339)),
			"",
			http.StatusOK,
			test.Effects[:3],
		},

		// Filter{From, To}
		{
			"GET",
			fmt.Sprintf("/effects?from=%s&to=%s", test.Effects[3].CreatedAt.Format(time.RFC3339), test.Effects[4].CreatedAt.Format(time.RFC3339)),
			"",
			http.StatusOK,
			test.Effects[3:5],
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

		effects := make(map[string][]effect.Effect)
		err := json.Unmarshal([]byte(resp.Body.String()), &effects)
		if err != nil {
			t.Error(err)
		}

		_, ok := effects["effects"]
		if !ok {
			t.Error(`Expected element "effects" in JSON response`)
		}

		if len(effects["effects"]) != len(test.expectedStats) {
			t.Errorf("Expected %d JSON elements, got %d", len(test.expectedStats), len(effects["effects"]))
		}

		for _, e := range test.expectedStats {
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
