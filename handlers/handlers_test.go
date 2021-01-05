package handlers

import (
        "net/http"
        "net/http/httptest"
	"strings"
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

	got := removeAllSpaces(rr.Body.String())
	expected := removeAllSpaces(`
	<html>
		<head></head>
		<body>
			<canvas>Classic 2D Ping Pong</canvas>
		</body>
	</html>`)

	if got != expected {
		t.Errorf("Got %v want %v", got, expected)
	}
}

func removeAllSpaces(str string) string {
	str_without_new_lines := strings.ReplaceAll(str, "\n", "")

	return strings.ReplaceAll(str_without_new_lines, "\t", "")
}
