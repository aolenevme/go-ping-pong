package main

import (
	"compress/gzip"
	"io"
	"net/http"
)

type gzipResponseWriter struct {
    io.Writer
    http.ResponseWriter
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
    return w.Writer.Write(b)
}

func main() {
	fs := http.FileServer(http.Dir("./static"))

	_ = http.ListenAndServeTLS(":8080", "security/cert.pem", "security/cert.key", middleware(fs))
}

func middleware(h http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	    w.Header().Set("Cache-Control", "max-age=2592000")
	    w.Header().Set("Content-Encoding", "gzip")

	    gzw := gzip.NewWriter(w)
	    defer gzw.Close()

	    h.ServeHTTP(&gzipResponseWriter{ResponseWriter: w, Writer: gzw}, r)
    })
}


