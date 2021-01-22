package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"
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

func setCookieMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			rand.Seed(time.Now().UnixNano())
			setCookieValue := fmt.Sprintf("client-id=%d; Secure; HttpOnly; SameSite=Strict", rand.Int())
			w.Header().Set("Set-Cookie", setCookieValue)
		}

		next.ServeHTTP(w, r)
	})
}
