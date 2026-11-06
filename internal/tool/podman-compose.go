package tool

import (
	"cdt/internal"
	"context"
)

const IdPodmanCompose = "podman-compose"

// A PodmanCompose wraps podman compose tool to manage tools environment.
type PodmanCompose struct {
	containerCompose
}

// NewPodmanCompose creates a podman compose tool from a custom podman executable.
func NewPodmanCompose(detect internal.ExecutableToolDetectFunc) *PodmanCompose {
	return &PodmanCompose{
		containerCompose{
			internal.MakeExecutableTool(
				IdPodmanCompose,
				"Podman compose",
				"Define and run multi-container applications with Podman",
				internal.Tags{internal.ToolTagEnvironment},
				detect,
			),
		},
	}
}

// DetectPodmanCompose create a podman compose tool with detected podman executable in the given environment.
func DetectPodmanCompose(
	ctx context.Context,
	options internal.DetectOptions,
) *PodmanCompose {
	return NewPodmanCompose(func() (*internal.Executable, error) {
		return options.Environment.FindExecutable(
			ctx,
			options.GetToolPath(IdPodmanCompose, "podman"),
		)
	})
}

func (d *PodmanCompose) Aliases() []string {
	return []string{"dc"}
}
