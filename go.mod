module github.com/spice-framework/starter-otel

go 1.26.0

toolchain go1.26.5

require (
	github.com/spice-framework/spice v0.0.0-20260805175412-383c17744300
	go.opentelemetry.io/otel v1.44.0
	go.opentelemetry.io/otel/metric v1.44.0
	go.opentelemetry.io/otel/sdk v1.44.0
	go.opentelemetry.io/otel/sdk/metric v1.44.0
	go.opentelemetry.io/otel/trace v1.44.0
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/spice-framework/development v0.0.0-20260806034648-1856466df09d // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	golang.org/x/sys v0.47.0 // indirect
)

tool github.com/spice-framework/development/cmd/spice-dev
