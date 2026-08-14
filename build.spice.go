//go:build spice_config

package spice

import "github.com/spice-framework/spice/project"

var Build = project.Build{
	Kind: project.StarterKind,
	Dependencies: project.Dependencies{
		project.Library("go.opentelemetry.io/otel", "v1.44.0"),
		project.Library("go.opentelemetry.io/otel/metric", "v1.44.0"),
		project.Library("go.opentelemetry.io/otel/sdk", "v1.44.0"),
		project.Library("go.opentelemetry.io/otel/sdk/metric", "v1.44.0"),
		project.Library("go.opentelemetry.io/otel/trace", "v1.44.0"),
		project.BuildTool("github.com/spice-framework/toolchain", "v0.0.0-20260806133530-71211498297c"),
	},
}
