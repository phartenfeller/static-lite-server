package main

import (
	"fmt"
	"log"
	"net/http"
)

func alwaysMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Server", gConfig.HeaderServerName)
		next.ServeHTTP(w, r)
	})
}

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(fmt.Sprintf("Method=%s Url=%s RemoteAddr=%s UserAgent=%s Referrer=%s", r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent(), r.Referer()))
		next.ServeHTTP(w, r)
	})
}
