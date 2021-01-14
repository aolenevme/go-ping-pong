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

func (w gzipResponseWriter) Write(b []byte) (int, error) {
        return w.Writer.Write(b)
}

func gzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Encoding", "gzip")

		gzw := gzip.NewWriter(w)
		defer gzw.Close()

		next.ServeHTTP(gzipResponseWriter{ResponseWriter: w, Writer: gzw}, r)
	})
}

func cacheControlMiddleware(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                w.Header().Set("Cache-Control", "max-age=3600")
                next.ServeHTTP(w, r)
        })
}
