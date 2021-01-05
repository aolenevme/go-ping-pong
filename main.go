package main

import (
	"net/http"
	"github.com/eshekak/go-ping-pong/handlers"
)

func main() {
	http.HandleFunc("/", handlers.MainPageHandler)
	http.ListenAndServe(":8080", nil)
}
