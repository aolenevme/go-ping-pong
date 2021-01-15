package main

import (
	"net/http"
	"time"
)

func main() {
	sse := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Content-Type", "text/event-stream")

		for {
			w.Write([]byte("data: Hey\n\n"))
			w.(http.Flusher).Flush()
			time.Sleep(10 * time.Millisecond)
		}
	})

	http.Handle("/api/v1/sse", sse)

	http.Handle("/", gzipMiddleware(cacheControlMiddleware(http.FileServer(http.Dir("./static")))))

	_ = http.ListenAndServeTLS(":8080", "security/cert.pem", "security/cert.key", nil)
}
