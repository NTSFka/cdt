package tool

import (
	"cdt/internal"
	"context"
)

const IdDockerCompose = "docker-compose"

// A DockerCompose wraps docker compose tool to manage tools environment.
type DockerCompose struct {
	containerCompose
}

// NewDockerCompose creates a docker compose tool from a custom docker executable.
func NewDockerCompose(detect internal.ExecutableToolDetectFunc) *DockerCompose {
	return &DockerCompose{
		containerCompose: containerCompose{
			ExecutableTool: internal.MakeExecutableTool(
				IdDockerCompose,
				"Docker compose",
				"Define and run multi-container applications with Docker",
				internal.Tags{internal.ToolTagEnvironment},
				detect,
			),
		},
	}
}

// DetectDockerCompose create a docker compose tool with detected docker executable in the given environment.
func DetectDockerCompose(
	ctx context.Context,
	options internal.DetectOptions,
) *DockerCompose {
	return NewDockerCompose(func() (*internal.Executable, error) {
		return options.Environment.FindExecutable(
			ctx,
			options.GetToolPath(IdDockerCompose, "docker"),
		)
	})
}

func (d *DockerCompose) Aliases() []string {
	return []string{"dc"}
}
