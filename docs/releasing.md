# Release contract

starter-otel releases are ordinary Go module tags plus a small, independently
verifiable artifact set. The repository owns the release contract while the
organization-owned reusable workflow performs the common build, signing,
independent verification, and publication phases. No mutable workspace snapshot
or network-resolved package list participates in artifact construction.

For `v1.2.3`, the release builder produces:

| Artifact | Contract |
|---|---|
| `starter-otel_1.2.3_source.tar.gz` | Exact tagged Git commit, under one versioned directory |
| `starter-otel_1.2.3_sbom.spdx.json` | SPDX 2.3 packages from the consistent committed `go.mod`, `go.sum`, and `vendor/modules.txt` graph |
| `checksums.txt` | SHA-256 of the source archive and SBOM, sorted by filename |
| `checksums.txt.sig` | Raw Ed25519 signature of the exact checksum bytes |
| `checksums.txt.pem` | X.509 SubjectPublicKeyInfo PEM for signature verification |

The source archive is reconstructed from the full commit's `git ls-tree`
identity and exact object bytes read through `git cat-file --batch`. It never
uses checkout filters or `git archive`, so `core.autocrlf` and host line-ending
settings cannot alter an artifact. Every tar and gzip timestamp is the source
commit epoch; paths are relative, ownership is zeroed, executable modes and
validated symlinks are preserved, and gzip output is deterministic. Gitlinks
and unsupported modes fail closed. Dirty or untracked workspace files cannot
enter the archive. The SBOM creation time uses the same epoch and contains no
absolute checkout path. Construction fails when committed module selection,
checksums, and vendored versions or replacements disagree; the builder does
not rely on an earlier verifier to detect a stale dependency graph.

## Production ceremony

Production releases call the organization-owned reusable workflow at an
immutable commit. Before any release tag is created, a release owner must:

1. generate a user-owned Ed25519 private key dedicated to this repository;
2. review and commit its public key as
   `security/release/ed25519-public.pem` (SHA-256 fingerprint
   `dd6fc8e44c6051d387d026481bb3871be7e24685aa6a37313d55e4e2ac82528b`);
3. store the private key as `SPICE_LIBRARY_RELEASE_SIGNING_KEY` only in the
   protected `release-signing` environment; and
4. configure protected `release-signing` and `release-publish` environments
   with the required human reviewers.

Do not create or push a release tag until all four controls exist. The caller
maps no secrets. The reusable workflow obtains the signing key only from its
`release-signing` job, validates the exact tag and public trust anchor, signs
with the centrally pinned tool, independently verifies with the separately
pinned verifier, and publishes only through `release-publish`. A missing key,
anchor, environment, review, or verification result fails closed.

## Unsigned dual-builder rehearsal

The application module authorizes an exact central renderer through its
`go.mod` tool directive. `make release-parity` runs that fully qualified tool
and the retained repository builder twice each with `GOWORK=off`,
`GOPROXY=off`, `GOTOOLCHAIN=local`, and `GOFLAGS=-mod=vendor`. It first asks
the central tool for a read-only plan, then renders the plan without resolving
an ambient workspace or downloading a module.

The central renderer and signer are the production implementation. The retained
repository builder remains only an unsigned parity oracle during the migration:

```text
make release-parity
```

Both rehearsals are unsigned and always archive `HEAD`, never working-tree
contents. Their source archives must be byte-identical and each builder must be
byte-deterministic across two runs. Artifact inputs must be regular files and
are size-checked before bounded reads. Parity decodes, bounds, and completely
drains both PAX/gzip streams; hidden decompressed data, raw compressed suffixes,
an additional gzip member, unsafe roots, duplicate entries, unsupported
metadata, oversized entries, or oversized aggregate content fail closed. Each
canonical checksum file must verify its own archive and SBOM, and neither output
may contain a signature or public key.

The SBOMs must be semantically identical except for these intentional builder
provenance fields:

- document name (`starter-otel VERSION` centrally and
  `Spice OpenTelemetry starter VERSION` in the retained builder);
- document namespace shape (the central renderer includes `spdx/v1/`);
- organization and tool creators identifying the actual renderer.

The central renderer uses `Organization: Spice Framework`; the retained
builder uses its existing `Organization: Spice Authors` identity. Package
facts, relationships including `DESCRIBES`, creation time, SPDX contract, and
every other decoded field must match exactly. Because the SBOM bytes differ,
the checksum files differ only in the SBOM digest; the source archive checksum
is identical. The parity gate fails closed on any extra artifact, dependency
drift, malformed or noncanonical checksum, or undocumented SBOM difference.

`make verify-release` executes this dual-builder rehearsal after the complete
repository verification contract. The retained builder is not removed by this
cutover and never receives production signing authority; removal requires a
separate reviewed change after the central signed path has durable evidence.
## Consumer verification

With OpenSSL 3 and GNU-compatible checksum tooling:

```text
sha256sum -c checksums.txt
openssl pkeyutl -verify -pubin -inkey checksums.txt.pem \
  -rawin -in checksums.txt -sigfile checksums.txt.sig
```

Consumers must authenticate `checksums.txt.sig` against the reviewed
`security/release/ed25519-public.pem` from the exact tagged source, not against a
public key supplied only beside release assets. Publishing remains fail-closed
if the anchor, signing secret, protected environments, or immutable tag rules
are absent.
