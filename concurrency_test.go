package otel

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	globalotel "go.opentelemetry.io/otel"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"

	spiceevent "github.com/spice-framework/spice/event"
	"github.com/spice-framework/spice/web"
)

func TestHTTPCompletionIsConcurrentAndIdempotent(t *testing.T) {
	t.Parallel()
	exporter, reader, tracerProvider, meterProvider := newInMemoryProviders(t)
	observer, err := NewHTTPObserver(Options{
		TracerProvider: tracerProvider,
		MeterProvider:  meterProvider,
	})
	if err != nil {
		t.Fatalf("NewHTTPObserver() error = %v", err)
	}
	_, finish := observer.BeginHTTP(context.Background(), web.RouteMetadata{
		ID:      "orders.create",
		Module:  "example.com/shop/orders",
		Method:  "POST",
		Pattern: "/orders",
	})
	var wait sync.WaitGroup
	for range 64 {
		wait.Go(func() {
			finish(web.HTTPResult{Status: 201, Bytes: 32, Duration: 10 * time.Millisecond})
		})
	}
	wait.Wait()
	if spans := exporter.GetSpans(); len(spans) != 1 {
		t.Fatalf("completed spans = %d, want 1", len(spans))
	}
	metrics := collectMetrics(t, reader)
	if got := int64Sum(t, metrics, "http.server.request.count"); got != 1 {
		t.Fatalf("request count = %d, want 1", got)
	}
	if got := int64Sum(t, metrics, "http.server.active_requests"); got != 0 {
		t.Fatalf("active requests = %d, want 0", got)
	}
	if got := histogramCount[int64](t, metrics, "http.server.response.body.size"); got != 1 {
		t.Fatalf("response size count = %d, want 1", got)
	}
}

func TestEventCompletionIsConcurrentIdempotentAndPayloadFree(t *testing.T) {
	t.Parallel()
	exporter, reader, tracerProvider, meterProvider := newInMemoryProviders(t)
	observer, err := NewEventObserver(Options{
		TracerProvider: tracerProvider,
		MeterProvider:  meterProvider,
	})
	if err != nil {
		t.Fatalf("NewEventObserver() error = %v", err)
	}
	interaction := spiceevent.Interaction{
		Event: spiceevent.Definition{ID: "order.placed", Module: "orders"},
		Subscriber: spiceevent.SubscriberMetadata{
			ID:     "inventory.reserve",
			Module: "inventory",
		},
	}
	_, finish := observer.BeginEvent(context.Background(), interaction)
	const secret = "customer-card-secret"
	var wait sync.WaitGroup
	for range 64 {
		wait.Go(func() {
			finish(spiceevent.Result{
				Interaction: interaction,
				Duration:    5 * time.Millisecond,
				Err:         errors.New(secret),
			})
		})
	}
	wait.Wait()
	spans := exporter.GetSpans()
	if len(spans) != 1 {
		t.Fatalf("completed spans = %d, want 1", len(spans))
	}
	for _, attribute := range spans[0].Attributes {
		if fmt.Sprint(attribute.Value.AsInterface()) == secret {
			t.Fatalf("span attribute exposed failure text: %#v", attribute)
		}
	}
	metrics := collectMetrics(t, reader)
	if got := int64Sum(t, metrics, "spice.event.delivery.count"); got != 1 {
		t.Fatalf("delivery count = %d, want 1", got)
	}
	if got := int64Sum(t, metrics, "spice.event.active_deliveries"); got != 0 {
		t.Fatalf("active deliveries = %d, want 0", got)
	}
	if got := histogramCount[float64](t, metrics, "spice.event.delivery.duration"); got != 1 {
		t.Fatalf("delivery duration count = %d, want 1", got)
	}
}

func TestConstructionDoesNotReplaceGlobalProviders(t *testing.T) {
	tracerGlobal := globalotel.GetTracerProvider()
	meterGlobal := globalotel.GetMeterProvider()
	_, _, tracerProvider, meterProvider := newInMemoryProviders(t)
	if _, err := NewObserver(Options{
		TracerProvider: tracerProvider,
		MeterProvider:  meterProvider,
	}); err != nil {
		t.Fatalf("NewObserver() error = %v", err)
	}
	if globalotel.GetTracerProvider() != tracerGlobal ||
		globalotel.GetMeterProvider() != meterGlobal {
		t.Fatal("NewObserver() replaced a global OpenTelemetry provider")
	}
}

func newInMemoryProviders(
	t *testing.T,
) (*tracetest.InMemoryExporter, *sdkmetric.ManualReader, *sdktrace.TracerProvider, *sdkmetric.MeterProvider) {
	t.Helper()
	exporter := tracetest.NewInMemoryExporter()
	reader := sdkmetric.NewManualReader()
	tracerProvider := sdktrace.NewTracerProvider(sdktrace.WithSyncer(exporter))
	meterProvider := sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader))
	t.Cleanup(func() {
		if err := tracerProvider.Shutdown(context.Background()); err != nil {
			t.Errorf("TracerProvider.Shutdown() error = %v", err)
		}
		if err := meterProvider.Shutdown(context.Background()); err != nil {
			t.Errorf("MeterProvider.Shutdown() error = %v", err)
		}
	})
	return exporter, reader, tracerProvider, meterProvider
}

func collectMetrics(t *testing.T, reader *sdkmetric.ManualReader) metricdata.ResourceMetrics {
	t.Helper()
	var metrics metricdata.ResourceMetrics
	if err := reader.Collect(context.Background(), &metrics); err != nil {
		t.Fatalf("Collect() error = %v", err)
	}
	return metrics
}

func int64Sum(t *testing.T, metrics metricdata.ResourceMetrics, name string) int64 {
	t.Helper()
	data := metricData(t, metrics, name)
	sum, ok := data.(metricdata.Sum[int64])
	if !ok {
		t.Fatalf("metric %q data = %T, want metricdata.Sum[int64]", name, data)
	}
	var result int64
	for _, point := range sum.DataPoints {
		result += point.Value
	}
	return result
}

func histogramCount[N int64 | float64](
	t *testing.T,
	metrics metricdata.ResourceMetrics,
	name string,
) uint64 {
	t.Helper()
	data := metricData(t, metrics, name)
	histogram, ok := data.(metricdata.Histogram[N])
	if !ok {
		t.Fatalf("metric %q data = %T, want histogram", name, data)
	}
	var result uint64
	for _, point := range histogram.DataPoints {
		result += point.Count
	}
	return result
}

func metricData(
	t *testing.T,
	metrics metricdata.ResourceMetrics,
	name string,
) metricdata.Aggregation {
	t.Helper()
	for _, scope := range metrics.ScopeMetrics {
		for _, item := range scope.Metrics {
			if item.Name == name {
				return item.Data
			}
		}
	}
	t.Fatalf("metrics do not contain %q", name)
	return nil
}
