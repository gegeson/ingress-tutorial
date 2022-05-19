package main

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

const shutdownTimeout = 5 * time.Second

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log.SetFlags(0)
	log.Println("GOMAXPROCS:", runtime.GOMAXPROCS(0))

	// setup http server
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)

	srv := &http.Server{
		Addr:    ":" + os.Getenv("PORT"),
		Handler: mux,
	}

	go func(ctx context.Context) {
		if err := srv.ListenAndServe(); err != nil {
			if errors.Is(ctx.Err(), context.Canceled) {
				return
			}
			log.Println("ListenAndServe error:", err)
		}
	}(ctx)

	log.Println("Starting server on", srv.Addr)
	<-ctx.Done()

	// graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Println("Shutdown error:", err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	env := os.Getenv("GOMAXPROCS")
	if env == "1" {
		io.WriteString(w, `{"message": "111"}`)
	} else if env == "2" {
		io.WriteString(w, `{"message": "222"}`)
	}
}
