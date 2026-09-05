package main

import (
	"context"
	"log"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func main() {
	ctx := context.Background()

	shutdownTelemetry, err := initTelemetry(ctx)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := shutdownTelemetry(ctx); err != nil {
			log.Printf("failed to shutdown telemetry: %v", err)
		}
	}()

	mux := http.NewServeMux()

	mux.HandleFunc("/health", health)
	mux.HandleFunc("/api/orders", orders)
	mux.HandleFunc("/api/transaction", transaction)
	mux.HandleFunc("/api/users", users)

	handler := otelhttp.NewHandler(
		mux,
		"HTTP Server",
	)

	server := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	log.Println("API Running on :8080")

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}