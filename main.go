package main

import "net/http"

func main() {
	 _ = http.ListenAndServeTLS(":8080", "security/cert.pem", "security/cert.key", http.FileServer(http.Dir("./static")))
}
