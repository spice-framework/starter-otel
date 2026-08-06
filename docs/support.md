# Support policy

| Contract | Current development support |
|---|---|
| Go | Exactly 1.26.5 for development and release verification |
| Spice minimum/current | Exact versions in [`spice-compatibility.json`](../spice-compatibility.json) |
| OpenTelemetry Go | API and SDK v1.44.0 |
| Operating systems | Windows, Linux, and macOS |
| Architectures | amd64 and arm64 compilation through public Go APIs |
| Exporters | Application-owned; none bundled |
| Central release tools | `github.com/spice-framework/development/cmd/spice-dev` at `v0.0.0-20260806121906-963bb6676069`; `github.com/spice-framework/toolchain/cmd/spice-library-release-verify` at `v0.0.0-20260806054457-a83d9b58034c` |

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
production workflow. Windows and Linux CI still compare unsigned central and
retained outputs under vendor-only offline resolution; the retained command is
only a parity oracle. Production remains disabled until a reviewed
`security/release/ed25519-public.pem`, the per-repository
`SPICE_LIBRARY_RELEASE_SIGNING_KEY`, and protected `release-signing` and
`release-publish` environments are configured.
