package tool

import (
	"cdt/internal"
	"context"
)

const IdPodman = "podman"

// A Podman is a tool that wraps podman executable.
type Podman struct {
	container
}

// NewPodman creates a podman tool from a custom executable.
func NewPodman(detect internal.ExecutableToolDetectFunc) *Podman {
	return &Podman{container: container{
		ExecutableTool: internal.MakeExecutableTool(
			IdPodman,
			"Podman",
			"Manage containers, pods, and images with Podman.",
			internal.Tags{internal.ToolTagEnvironment},
			detect,
		),
	}}
}

// DetectPodman create a podman tool with a detected podman executable in the given environment.
func DetectPodman(
	ctx context.Context,
	options internal.DetectOptions,
) *Podman {
	return NewPodman(func() (*internal.Executable, error) {
		return options.Environment.FindExecutable(ctx, options.GetToolPath(IdPodman, "podman"))
	})
}

func (d *Podman) Aliases() []string {
	return []string{"p"}
}
