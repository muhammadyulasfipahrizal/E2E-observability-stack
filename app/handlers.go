package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/global"
)

var logger = global.Logger("observability-api")

func emitLog(
	ctx context.Context,
	severity log.Severity,
	message string,
	attrs ...attribute.KeyValue,
) {
	var record log.Record

	record.SetTimestamp(time.Now())
	record.SetSeverity(severity)
	record.SetSeverityText(severity.String())
	record.SetBody(attribute.StringValue(message))
	record.AddAttributes(attrs...)

	logger.Emit(ctx, record)
}

type Order struct {
	ID    int     `json:"id"`
	Item  string  `json:"item"`
	Price float64 `json:"price"`
}

func health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	emitLog(
		r.Context(),
		log.SeverityInfo,
		"health check",
		attribute.String("endpoint", "/health"),
	)

	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}

func orders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	orders := []Order{
		{ID: 1, Item: "Laptop", Price: 1200},
		{ID: 2, Item: "Keyboard", Price: 100},
		{ID: 3, Item: "Monitor", Price: 300},
	}

	emitLog(
		r.Context(),
		log.SeverityInfo,
		"orders fetched successfully",
		attribute.String("endpoint", "/api/orders"),
		attribute.Int("order_count", len(orders)),
	)

	json.NewEncoder(w).Encode(orders)
}

func transaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	start := time.Now()

	emitLog(
		r.Context(),
		log.SeverityInfo,
		"transaction started",
		attribute.String("endpoint", "/api/transaction"),
	)

	time.Sleep(3 * time.Second)

	duration := time.Since(start)

	emitLog(
		r.Context(),
		log.SeverityInfo,
		"transaction completed",
		attribute.String("endpoint", "/api/transaction"),
		attribute.Int64("duration_ms", duration.Milliseconds()),
	)

	json.NewEncoder(w).Encode(map[string]string{
		"status":  "completed",
		"message": "transaction completed",
	})
}

func users(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	emitLog(
		r.Context(),
		log.SeverityError,
		"failed to fetch users",
		attribute.String("endpoint", "/api/users"),
		attribute.Int("status_code", http.StatusInternalServerError),
	)

	http.Error(
		w,
		`{"error":"failed to fetch users"}`,
		http.StatusInternalServerError,
	)
}