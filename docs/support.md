# Support policy

| Contract | Current development support |
|---|---|
| Go | Exactly 1.26.5 for development and release verification |
| Spice minimum/current | Exact versions in [`spice-compatibility.json`](../spice-compatibility.json) |
| OpenTelemetry Go | API and SDK v1.44.0 |
| Operating systems | Windows, Linux, and macOS |
| Architectures | amd64 and arm64 compilation through public Go APIs |
| Exporters | Application-owned; none bundled |
| Central release tools | `github.com/spice-framework/development/cmd/spice-dev` at `v0.0.0-20260806132124-4c308d1b9fda`; `github.com/spice-framework/toolchain/cmd/spice-library-release-verify` at `v0.0.0-20260806133530-71211498297c` |

`spice-compatibility.json` is the sole compatibility boundary source. Its
minimum must equal the exact direct Spice requirement in `go.mod`; its current
value is a forward-compatibility endpoint rather than a moving branch. The
repository gate verifies both boundaries using isolated alternate modfiles,
exact MVS selection, vet, and shuffled race tests without modifying product or
module files. A release may raise the minimum only through an intentional
module and compatibility-contract change with green minimum/current evidence.

Traces and metrics are supported. The OpenTelemetry logs API and all exporters
are deliberately outside this starter until they receive independent API,
dependency, maintenance, cancellation, security, and operability review.

Release artifacts are produced only from an exact tagged commit under the
contract in [`releasing.md`](releasing.md). A compromised or missing signing
secret fails a production release; it never falls back to unsigned output.

The pinned central signer and independent verifier power the protected reusable
production workflow. The reviewed repository-specific trust anchor is
`security/release/ed25519-public.pem` (SHA-256 fingerprint
`dd6fc8e44c6051d387d026481bb3871be7e24685aa6a37313d55e4e2ac82528b`).
Its private key exists only in the protected `release-signing` environment;
`release-publish` contains no signing secret. Windows and Linux CI still compare
unsigned central and retained outputs under vendor-only offline resolution; the
retained command is only a parity oracle until the first signed cutover passes.
