package tool

import (
	"cdt/internal"
	"context"
)

const IdDocker = "docker"

// A Docker is a tool that wraps docker executable.
type Docker struct {
	container
}

// NewDocker creates a docker tool from a custom executable.
func NewDocker(detect internal.ExecutableToolDetectFunc) *Docker {
	return &Docker{container: container{
		ExecutableTool: internal.MakeExecutableTool(
			IdDocker,
			"Docker",
			"Docker image and container command line interface.",
			internal.Tags{internal.ToolTagEnvironment},
			detect,
		),
	}}
}

// DetectDocker create a docker tool with a detected docker executable in the given environment.
func DetectDocker(
	ctx context.Context,
	options internal.DetectOptions,
) *Docker {
	return NewDocker(func() (*internal.Executable, error) {
		return options.Environment.FindExecutable(ctx, options.GetToolPath(IdDocker, "docker"))
	})
}

func (d *Docker) Aliases() []string {
	return []string{"d"}
}
