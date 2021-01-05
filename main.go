package main

import (
	"github.com/eshekak/go-ping-pong/handlers"
	"net/http"
)

func main() {
	http.HandleFunc("/", handlers.MainPageHandler)
	http.ListenAndServe(":8080", nil)
}
