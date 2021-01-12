package main

import (
	"log"
	"net/http"
)

func main() {
	if err := http.ListenAndServeTLS(":8080", "security/cert.pem", "security/cert.key", http.FileServer(http.Dir("./static"))); err != nil {
		log.Fatalf("%v", err)
	}
}
