package main

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
)

func initTelemetry(ctx context.Context) (func(context.Context) error, error) {
	res, err := resource.New(
		ctx,
		resource.WithAttributes(
			semconv.ServiceName("observability-api"),
		),
	)
	if err != nil {
		return nil, err
	}

	// Traces
	traceExporter, err := otlptracegrpc.New(
		ctx,
		otlptracegrpc.WithEndpoint("otel-collector:4317"),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	tracerProvider := trace.NewTracerProvider(
		trace.WithBatcher(traceExporter),
		trace.WithResource(res),
	)

	otel.SetTracerProvider(tracerProvider)

	// Metrics
	metricExporter, err := otlpmetricgrpc.New(
		ctx,
		otlpmetricgrpc.WithEndpoint("otel-collector:4317"),
		otlpmetricgrpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	meterProvider := metric.NewMeterProvider(
		metric.WithReader(
			metric.NewPeriodicReader(metricExporter),
		),
		metric.WithResource(res),
	)

	otel.SetMeterProvider(meterProvider)

	// Logs
	logExporter, err := otlploggrpc.New(
		ctx,
		otlploggrpc.WithEndpoint("otel-collector:4317"),
		otlploggrpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	logProcessor := log.NewBatchProcessor(logExporter)

	loggerProvider := log.NewLoggerProvider(
		log.WithProcessor(logProcessor),
		log.WithResource(res),
	)

	global.SetLoggerProvider(loggerProvider)

	// Shutdown
	shutdown := func(ctx context.Context) error {
		if err := tracerProvider.Shutdown(ctx); err != nil {
			return err
		}

		if err := meterProvider.Shutdown(ctx); err != nil {
			return err
		}

		if err := loggerProvider.Shutdown(ctx); err != nil {
			return err
		}

		return nil
	}

	return shutdown, nil
}