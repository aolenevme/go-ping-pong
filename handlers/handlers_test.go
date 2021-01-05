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
		<head>
			<style>
				body{margin:0;width:100vw;height:100vh}
				canvas{width:100%;height:100%;}
			</style>
		</head>
		<body>
			<canvas>Classic 2D Ping Pong</canvas>
		</body>
	</html>`)

	if got != expected {
		t.Errorf("Got %v want %v", got, expected)
	}
}

func removeAllSpaces(str string) string {
	intermediate_str := strings.ReplaceAll(str, "\n", "")
	intermediate_str = strings.ReplaceAll(intermediate_str, " ", "")

	return strings.ReplaceAll(intermediate_str, "\t", "")
}
