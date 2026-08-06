# Dependency review: OpenTelemetry Go

- Decision: approved for the independently versioned
  `github.com/spice-framework/starter-otel` module.
- Version: `go.opentelemetry.io/otel` and SDK v1.44.0.
- Upstream: <https://github.com/open-telemetry/opentelemetry-go>.
- License: Apache-2.0; retained in vendored module license files.
- Maintenance: vendor-neutral CNCF/OpenTelemetry project with signed immutable
  releases and an explicit supported-Go policy.
- Compatibility: v1.44.0 explicitly supports Go 1.26. Traces and metrics are
  stable; the beta logs API is deliberately not adopted.
- Security: exporter packages are not included by this starter. Applications
  choose and configure transport/export separately. `govulncheck` and `gosec`
  remain mandatory for the reachable repository graph.
- Cancellation: all span and measurement operations use the request context.
  Provider shutdown and exporter deadlines remain application-owned.
- Observability: the adapter emits server spans and four bounded metrics using
  generated route templates, stable symbol IDs, module paths, method, and
  status. It never labels telemetry with raw paths, queries, headers, or body
  values.
- Configuration: providers are explicit constructor inputs; the starter does
  not read environment variables, install globals, contact a collector, or
  select an exporter.
- Activation: the committed manifest must be explicitly selected and
  `@otel.Enable` must occur on the application marker. The compiler validates
  the exact HTTP-observer output and reachable mux capability before generated
  code directly composes the constructor result.
- Transitive scope: OpenTelemetry API/SDK support libraries plus their small
  logging, UUID, and platform support graph only. No collector,
  network exporter, protobuf, gRPC, or Prometheus dependency is accepted by
  this slice.

Primary references:

- <https://github.com/open-telemetry/opentelemetry-go/releases>
- <https://opentelemetry.io/docs/languages/go/>

## Build-only dependencies: central release tools

- Decision: approved only as repository-authorized release tooling.
- Renderer: `github.com/spice-framework/development/cmd/spice-dev` from
  `github.com/spice-framework/development`
  `v0.0.0-20260806052122-9025218a91c0`.
- Independent verifier:
  `github.com/spice-framework/toolchain/cmd/spice-library-release-verify` from
  `github.com/spice-framework/toolchain`
  `v0.0.0-20260806054457-a83d9b58034c`.
- Tool registration: both commands use standard Go `tool` directives and all
  invocations use their full package paths.
- License: Apache-2.0, with its notice retained in `vendor`.
- Runtime scope: none. Product packages import neither tool module, and
  released applications acquire no runtime dependency on them.
- Dependency graph: both tools participate in normal Go minimal-version
  selection. That build-time coupling is accepted and visible in `go.mod`,
  `go.sum`, and `vendor/modules.txt`; no parallel tool registry is introduced.
- Integrity and network behavior: both exact pseudo-versions are pinned and
  checksummed. Release parity runs with `GOWORK=off`, `GOPROXY=off`,
  `GOTOOLCHAIN=local`, and `GOFLAGS=-mod=vendor`, so it cannot select an ambient
  checkout, upgrade itself, or download dependencies.
- Security: the trusted native renderer reads the exact committed Git graph
  and writes only to caller-supplied temporary output directories. The
  verifier independently checks release artifacts without signing them. The
  rehearsal emits no signatures or signing material.
- Maintenance: the retained local builder and production signing workflow stay
  in place. A dual-builder gate detects central renderer regressions before any
  future authority migration.
