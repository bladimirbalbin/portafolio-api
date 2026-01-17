package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bladimirbalbin/portafolio-api/internal/config"
	aphttp "github.com/bladimirbalbin/portafolio-api/internal/http"
)

func main() {
	cfg := config.Load()

	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           aphttp.NewRouter(cfg),
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("server listening on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("shutting down server")
}
