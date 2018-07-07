package assets

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	expectedBody []string
}

func TestAssets(t *testing.T) {
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
			"POST",
			"/assets",
			`{"code": "TEST", "issuer": "GCJXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"}`,
			http.StatusOK,
			[]string{
				`"status":"ok"`,
				`"code":"TEST"`,
				`"issuer":"GCJXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"`,
			},
		},

		// Check whether duplicates are prevented
		{
			"POST",
			"/assets",
			`{"code": "TEST", "issuer": "GCJXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"}`,
			http.StatusConflict,
			[]string{
				`"status":"error"`,
				`"message":"Asset already exists"`,
			},
		},

		{
			"POST",
			"/assets",
			`{"code": "invalid`,
			http.StatusBadRequest,
			[]string{
				`"status":"error"`,
			},
		},

		// Check whether new asset as well as CNDY asset are present in database
		{
			"GET",
			"/assets",
			"",
			http.StatusOK,
			[]string{
				`"status":"ok"`,
				fmt.Sprintf(`"code":"%s"`, cndy.AssetCode),
				fmt.Sprintf(`"issuer":"%s"`, cndy.AssetIssuer),
				`"code":"TEST"`,
				`"issuer":"GCJXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"`,
			},
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

		for _, contains := range tt.expectedBody {
			if !strings.Contains(resp.Body.String(), contains) {
				t.Errorf("Body did not contain '%s' in '%s'", contains, resp.Body.String())
			}
		}
	}
}
