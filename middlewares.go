package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func alwaysMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Server", gConfig.HeaderServerName)
		next.ServeHTTP(w, r)
	})
}

func customFileServer(root http.FileSystem) http.Handler {
	log.Println("root", root)
	fs := http.FileServer(root)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rPath := r.URL.Path
		f, err := root.Open(rPath)
		if err != nil {
			notFoundHandler(w, r)
			return
		}

		s, err := f.Stat()
		if err != nil {
			log.Println("error stating file", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if s.IsDir() {
			index := filepath.Join(rPath, "index.html")
			_, err := root.Open(index)

			if err == nil {
				log.Println("path", filepath.Join(gConfig.Root, index))
				http.ServeFile(w, r, filepath.Join(gConfig.Root, index))
				return
			} else {
				notFoundHandler(w, r)
				return
			}
		}

		// source: https://gist.github.com/lummie/91cd1c18b2e32fa9f316862221a6fd5c
		if gConfig.NotFoundPage != "" {
			log.Println("here?")
			requestPath := r.URL.Path

			// make sure that url starts with /
			if !strings.HasPrefix(requestPath, "/") {
				requestPath = "/" + requestPath
				r.URL.Path = requestPath
			}
			requestPath = path.Clean(requestPath)

			// attempt to open the file via the http.FileSystem
			f, err := root.Open(requestPath)
			if err != nil {
				if os.IsNotExist(err) {
					w.WriteHeader(http.StatusNotFound)
					http.ServeFile(w, r, gConfig.NotFoundPage)
					return
				}
			}

			// close if successfully opened
			if err == nil {
				err := f.Close()
				if err != nil {
					log.Println("Error closing file", err)
					return
				}
			}
		}

		log.Println("fileServer", r.URL.Path)

		// default option
		fs.ServeHTTP(w, r)
	})

}

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(fmt.Sprintf("Method=%s Url=%s RemoteAddr=%s UserAgent=%s Referrer=%s", r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent(), r.Referer()))
		next.ServeHTTP(w, r)
	})
}
