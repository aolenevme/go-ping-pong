package handlers

import (
        "net/http"
        "net/http/httptest"
	"testing"
	_ "github.com/eshekak/go-ping-pong/testing"
)

func TestMainPageHandler(t *testing.T) {
	rr := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler := http.HandlerFunc(MainPageHandler)
	handler.ServeHTTP(rr, req)

	expected := "<html>\n</html>\n"
	if rr.Body.String() != expected {
		t.Errorf("Got %v want %v", rr.Body.String(), expected)
	}
}
