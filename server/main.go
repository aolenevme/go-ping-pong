package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func main() {
	lastEventId := 1
	sse := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Last-Event-ID", strconv.Itoa(lastEventId))

		for {
			fmt.Fprintf(w, "id: %d\ndata: Hey\n\n", lastEventId)
			w.(http.Flusher).Flush()
			time.Sleep(10 * time.Millisecond)
			lastEventId++
		}
	})

	http.Handle("/api/v1/sse", sse)

	http.Handle("/", gzipMiddleware(cacheControlMiddleware(http.FileServer(http.Dir("./static")))))

	_ = http.ListenAndServeTLS(":8080", "security/cert.pem", "security/cert.key", nil)
}
