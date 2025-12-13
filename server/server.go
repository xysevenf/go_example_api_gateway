package main

import (
	"api-gateway/config"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Printf("config load error")
		return
	}
	server := &http.Server{Addr: ":3000", Handler: NewRouter(cfg)}
	go func() {
		fmt.Println("Starting listening on :3000")
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server stopped with error: %v", err)
		}
	}()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	<-signalCh

	log.Printf("Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	server.Shutdown(ctx)
	log.Printf("Server stopped")
}
