package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var gConfig *config

func main() {
	argsWithoutProgram := os.Args[1:]
	args := parseArguments(argsWithoutProgram)

	c, err := parseConfig(args.configFilePath)
	gConfig = &c

	if err != nil {
		panic(err)
	}

	log.Println("Config =>", c.Port, c.Root, *c.LogRequests)

	var fs http.FileSystem = http.Dir(c.Root)
	handler := http.TimeoutHandler(customFileServer(fs), c.TimeoutMs*time.Millisecond, "Request timeout")

	handler = alwaysMiddleware(handler)

	if c.PathPrefix != "" {
		handler = http.StripPrefix(c.PathPrefix, handler)
	}

	if *c.LogRequests {
		handler = logMiddleware(handler)
	}

	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", c.Port),
		Handler:           handler,
		TLSConfig:         nil,
		ReadTimeout:       c.ReadTimeoutMs * time.Millisecond,
		ReadHeaderTimeout: c.ReadHeaderTimeoutMs * time.Millisecond,
		WriteTimeout:      0,
		IdleTimeout:       0,
		MaxHeaderBytes:    0,
		TLSNextProto:      nil,
		ConnState:         nil,
		ErrorLog:          nil,
		BaseContext:       nil,
		ConnContext:       nil,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			// handle error
			if err.Error() != "http: Server closed" {
				log.Println("ListenAndServe error:", err)
			}
		}
	}()

	// Setting up signal capturing
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// Waiting for SIGINT (kill -2)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Println("Shutdown error:", err)
	}

	// Wait for ListenAndServe goroutine to close.
}
