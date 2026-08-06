# Spice OpenTelemetry starter

`github.com/spice-framework/starter-otel` is the independently versioned,
opt-in OpenTelemetry adapter for Spice applications. It turns generated HTTP
route metadata and typed module-event interactions into traces and bounded
metrics while keeping provider, exporter, shutdown, and transport ownership in
the application.

## Why it is separate

Spice core exposes small observer seams and remains standard-library-first.
Applications that need OpenTelemetry import this module explicitly; applications
that do not need it do not acquire the OpenTelemetry dependency graph.

## Construction

```go
observer, err := otel.NewObserver(otel.Options{
	TracerProvider: tracerProvider,
	MeterProvider:  meterProvider,
})
if err != nil {
	return err
}
```

Both providers are required. The starter never installs global providers,
selects an exporter, reads environment variables, or contacts a collector.
Generated Spice code composes the returned observer directly with the HTTP and
event runtimes.

HTTP telemetry uses route templates and compiler-owned module identities, not
raw URLs, headers, queries, or request bodies. Event telemetry records stable
publisher/subscriber identities and outcomes without event payloads or error
text. Completion callbacks are idempotent.

## Manifest and activation

`Manifest` describes explicit `@otel.Enable` activation and the
`NewHTTPObserver` entrypoint. The application must select the starter and supply
the required providers; importing the package alone does not mutate runtime
state.

## Compatibility and verification

Development and verification require exactly Go 1.26.5. The machine-readable
[`spice-compatibility.json`](spice-compatibility.json) records the provisional
minimum and current supported Spice core lines. The complete gate proves both
lines through isolated modfiles, exact MVS selection, vet, shuffled race tests,
an 85% coverage floor, security analysis, reproducible vendor contents, and
offline tests.

```text
make check
make compatibility
make release-parity
make verify
make verify-release
```

Release parity validates the exact `spice-dev` renderer and
`spice-library-release-verify` verifier authorized by `go.mod`, then runs the
renderer and retained repository builder twice each, entirely from `vendor`
with network and workspace resolution disabled. It requires byte-identical
source archives, regular bounded artifact reads, complete bounded gzip/tar
validation, rejection of hidden data, raw suffixes, and extra members, equivalent SBOM facts,
canonical checksum files, and no rehearsal signatures on Windows or Linux.

See [`docs/dependency-review.md`](docs/dependency-review.md) for the dependency
and security review and [`docs/support.md`](docs/support.md) for the support
policy.

## Releases

Each version tag is an ordinary Go module release. The repository also builds
an exact-commit source archive, committed-graph SPDX 2.3 SBOM, SHA-256
checksums, and an Ed25519 signature/public key without an external release
build system. Production mode requires a clean checkout, exact tag, and
protected signing key; an explicit unsigned rehearsal is available for local
proof. See [`docs/releasing.md`](docs/releasing.md) for the artifact and trust
contract.
The protected central workflow is the release authority. The retained
repository builder remains only an unsigned parity oracle and is held to the
dual-builder contract during the migration.
