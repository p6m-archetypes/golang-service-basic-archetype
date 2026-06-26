module {{ module_path }}

go 1.23

require (
	github.com/go-chi/chi/v5 v5.2.1
	github.com/prometheus/client_golang v1.20.5
	go.opentelemetry.io/otel v1.33.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.33.0
	go.opentelemetry.io/otel/sdk v1.33.0
)
