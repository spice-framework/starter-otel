# Starter OpenTelemetry implementation contract

This repository owns the independently versioned OpenTelemetry integration for
Spice. Work directly on local `main` in bounded commits. Fetch before editing
and immediately before pushing; never overwrite unexpected remote work.

Go 1.26.5 is mandatory. Every product change must preserve caller-owned
providers and contexts, bounded low-cardinality attributes, module-aware
telemetry, idempotent completion callbacks, payload-free diagnostics, and the
public Spice HTTP and event observer contracts. The starter must never install
global providers, choose exporters, read ambient configuration, or perform
hidden network access.

Add positive and failure-path tests, update public documentation, run
`make verify` on the exact commit tree, and push only a green commit.

Release-parity work must preserve the exact `spice-dev` tool version authorized
by the root `go.mod`, invoke its full package path, and run both central and
retained rehearsals with workspace and network resolution disabled in vendor
mode. The retained repository builder and signed production workflow remain
authoritative until a separately reviewed signing migration; unsigned parity
must never manufacture signatures or key material.
