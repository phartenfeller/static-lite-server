package main

import (
	"net/http"
	"path/filepath"
)

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	if gConfig.NotFoundPage != "" {
		w.WriteHeader(http.StatusNotFound)
		http.ServeFile(w, r, filepath.Join(gConfig.Root, gConfig.NotFoundPage))
	} else {
		http.Error(w, "Not Found (404)", http.StatusNotFound)
	}
	return
}
