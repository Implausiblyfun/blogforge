// Package tracing attempts to put a minimal otel on jeager setup just in case we want to dig into whats going on in the blog.
package tracing

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

// defaulTraceURL is the initial endpoint so locally running (or composed running we can have our jeager collector)
const defaultTraceURL = "http://localhost:14268/api/traces"

// NewBaseTracer is meant to be a convience wrapping to set up a basic jeager pipeline.
func NewBaseTracer(reportURL, sectionName string) (*tracesdk.TracerProvider, error) {
	if reportURL == "" {
		reportURL = defaultTraceURL
	}

	// Create the Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(reportURL)))
	if err != nil {
		return nil, err
	}
	tp := tracesdk.NewTracerProvider(
		tracesdk.WithSampler(tracesdk.AlwaysSample()),
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(sectionName),
		)),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp, nil
}
