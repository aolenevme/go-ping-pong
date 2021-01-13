package main

import "net/http"

// TODO
// 1. Serve with cache
// 2. Compress with gzip
// 3. Make it work with tinygo

func main() {
	 _ = http.ListenAndServeTLS(":8080", "security/cert.pem", "security/cert.key", http.FileServer(http.Dir("./static")))
}
