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

Release-rehearsal work must preserve the exact central release-tool versions
authorized by the root `go.mod`, invoke their full package paths, and render
the same inert plan twice with workspace and network resolution disabled in
vendor mode. The protected central workflow is the sole production builder. An
unsigned rehearsal must never manufacture signatures or key material.
