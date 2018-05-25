package docs

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStats(t *testing.T) {
	router := gin.Default()
	Init(router)

	req, _ := http.NewRequest("GET", "/", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected code %v, got %v", http.StatusOK, resp.Code)
	}

	// Verify that content-type is correct
	if resp.HeaderMap["Content-Type"][0] != "text/html" {
		t.Errorf("Expected Content-Type header text/html, got %s", resp.HeaderMap["Content-Type"][0])
	}
}
