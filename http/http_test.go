package http

import (
        "net/http"
        "net/http/httptest"
	"testing"
)

func TestHelloHandler(t *testing.T) {
	rr := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler := http.HandlerFunc(HelloHandler)
	handler.ServeHTTP(rr, req)

	expected := "Hello, world!\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
