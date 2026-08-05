package otel

import (
	"slices"
	"testing"

	spicestarter "github.com/spice-framework/spice/annotation/sdk/starter"
)

func TestManifestDeclaresIndependentOpenTelemetryStarter(t *testing.T) {
	t.Parallel()
	manifest := Manifest()
	spec := manifest.Spec()
	if spec.Schema != spicestarter.Schema {
		t.Fatalf("Manifest().Schema = %q", spec.Schema)
	}
	if spec.ID != "github.com/spice-framework/starter-otel" {
		t.Fatalf("Manifest().ID = %q", spec.ID)
	}
	if spec.Module != "github.com/spice-framework/starter-otel" {
		t.Fatalf("Manifest().Module = %q", spec.Module)
	}
	if spec.SpiceAPI != spicestarter.APIVersion {
		t.Fatalf("Manifest().SpiceAPI = %q, want %q", spec.SpiceAPI, spicestarter.APIVersion)
	}
	if spec.Review != "docs/dependency-review.md" {
		t.Fatalf("Manifest().Review = %q", spec.Review)
	}
	if len(spec.Annotations) != 1 || spec.Annotations[0].Name != "otel.Enable" {
		t.Fatalf("Manifest().Annotations = %#v", spec.Annotations)
	}
	if len(spec.ApplicationFeatures) != 1 ||
		spec.ApplicationFeatures[0].EntryPoints[0].Package != "github.com/spice-framework/starter-otel" {
		t.Fatalf("Manifest().ApplicationFeatures = %#v", spec.ApplicationFeatures)
	}
	wantCapabilities := []string{
		"observability.http-server",
		"observability.metrics",
		"observability.module-events",
		"observability.tracing",
	}
	if !slices.Equal(spec.Capabilities, wantCapabilities) {
		t.Fatalf("Manifest().Capabilities = %v, want %v", spec.Capabilities, wantCapabilities)
	}
	if len(spec.Dependencies) != 3 ||
		spec.Dependencies[0].Module != "go.opentelemetry.io/otel" ||
		spec.Dependencies[1].Module != "go.opentelemetry.io/otel/metric" ||
		spec.Dependencies[2].Module != "go.opentelemetry.io/otel/trace" {
		t.Fatalf("Manifest().Dependencies = %#v", spec.Dependencies)
	}
	if err := manifest.Compatible(spicestarter.APIVersion, "go1.26.5"); err != nil {
		t.Fatalf("Compatible() error = %v", err)
	}
	content, err := manifest.JSON()
	if err != nil {
		t.Fatalf("JSON() error = %v", err)
	}
	parsed, err := spicestarter.Parse(content)
	if err != nil {
		t.Fatalf("Parse(JSON()) error = %v", err)
	}
	if parsed.Spec().ID != spec.ID {
		t.Fatalf("parsed ID = %q, want %q", parsed.Spec().ID, spec.ID)
	}
}
