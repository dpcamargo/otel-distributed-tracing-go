package main

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

// BenchmarkSpanOverhead measures the cost of the Start/End span pair this
// service does per request/outgoing-call (tracer.Start + span.End), using an
// in-memory exporter so the number reflects pure OpenTelemetry SDK overhead
// — not network latency to the collector. This is the overhead "tax" paid
// on every request regardless of what the OTel collector/backend does with
// the data afterwards.
//
// Run with:
//
//	go test ./cmd/service-a/... -run=^$ -bench=. -benchmem
func BenchmarkSpanOverhead(b *testing.B) {
	exporter := tracetest.NewInMemoryExporter()
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithSyncer(exporter),
	)
	otel.SetTracerProvider(tp)
	tracer := otel.Tracer("benchmark-tracer")

	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, span := tracer.Start(ctx, "outgoing request to service-b")
		span.End()
	}
}

// BenchmarkSpanOverhead_NoSampling isolates the case where the sampler drops
// the span (as it would with a probabilistic sampler under high load),
// which is cheaper than AlwaysSample because no attributes/events get
// recorded or exported.
func BenchmarkSpanOverhead_NoSampling(b *testing.B) {
	exporter := tracetest.NewInMemoryExporter()
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.NeverSample()),
		sdktrace.WithSyncer(exporter),
	)
	otel.SetTracerProvider(tp)
	tracer := otel.Tracer("benchmark-tracer")

	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, span := tracer.Start(ctx, "outgoing request to service-b")
		span.End()
	}
}
