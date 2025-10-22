package observability

import (
	"context"
	"net/url"
	"os"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

func InitTracer(ctx context.Context, service string) (*sdktrace.TracerProvider, error) {
	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	var exp sdktrace.SpanExporter
	var err error

	if endpoint != "" {
		// Try to decode the URL if it's double-encoded
		if decoded, decodeErr := url.QueryUnescape(endpoint); decodeErr == nil && decoded != endpoint {
			endpoint = decoded
		}

		// Parse and clean the endpoint URL
		opts := []otlptracehttp.Option{}

		// Handle different URL formats
		if strings.HasPrefix(endpoint, "http://") {
			endpoint = strings.TrimPrefix(endpoint, "http://")
			opts = append(opts, otlptracehttp.WithInsecure())
		} else if strings.HasPrefix(endpoint, "https://") {
			endpoint = strings.TrimPrefix(endpoint, "https://")
			// no WithInsecure() for https
		}

		// Validate the cleaned endpoint
		if _, err := url.Parse("http://" + endpoint); err != nil {
			// If URL parsing fails, treat as no endpoint configured
			endpoint = ""
		}

		if endpoint != "" {
			opts = append(opts, otlptracehttp.WithEndpoint(endpoint))

			// Add headers if configured
			if h := os.Getenv("OTEL_EXPORTER_OTLP_HEADERS"); h != "" {
				headers := map[string]string{}
				for _, kv := range strings.Split(h, ",") {
					if kv == "" {
						continue
					}
					parts := strings.SplitN(kv, "=", 2)
					if len(parts) == 2 {
						headers[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
					}
				}
				opts = append(opts, otlptracehttp.WithHeaders(headers))
			}

			exp, err = otlptracehttp.New(ctx, opts...)
			if err != nil {
				return nil, err
			}
		}
	} else {
		// No exporter configured -> no-op provider (cheap)
		provider := sdktrace.NewTracerProvider()
		otel.SetTracerProvider(provider)
		return provider, nil
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp, sdktrace.WithMaxExportBatchSize(2048)),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(service),
		)),
	)
	otel.SetTracerProvider(tp)
	return tp, nil
}

func ShutdownTracer(ctx context.Context, tp *sdktrace.TracerProvider) {
	if tp == nil {
		return
	}
	_ = tp.Shutdown(ctx)
}
