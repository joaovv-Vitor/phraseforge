package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joaovv-Vitor/phraseforge/internal/httpapi"
	"github.com/joaovv-Vitor/phraseforge/internal/storage"
)

func main() {
	categories, err := storage.LoadCategories("data/phrases.json")
	if err != nil {
		log.Fatal(err)
	}

	server := &http.Server{
		Addr:              apiAddress(),
		Handler:           httpapi.NewHandler(categories),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	listener, err := net.Listen("tcp", server.Addr)
	if err != nil {
		log.Fatal(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log.Printf("PhraseForge API listening on %s", server.Addr)
	if err := serve(ctx, server, listener); err != nil {
		log.Fatal(err)
	}
}

func serve(ctx context.Context, server *http.Server, listener net.Listener) error {
	serverErrors := make(chan error, 1)
	go func() {
		serverErrors <- server.Serve(listener)
	}()

	select {
	case err := <-serverErrors:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	case <-ctx.Done():
		log.Print("PhraseForge API shutdown started")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			return err
		}

		if err := <-serverErrors; !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	}
}
