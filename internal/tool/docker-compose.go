package tool

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"cdt/internal"
)

const IdDockerCompose = "docker-compose"

// A DockerCompose wraps docker compose tool to manage tools environment.
type DockerCompose struct {
	internal.ExecutableTool
}

// NewDockerCompose creates a docker compose tool from a custom docker executable.
func NewDockerCompose(detect internal.ExecutableToolDetectFunc) *DockerCompose {
	return &DockerCompose{
		internal.MakeExecutableTool(
			IdDockerCompose,
			"Docker compose",
			"Define and run multi-container applications with Docker",
			internal.Tags{internal.ToolTagEnvironment},
			detect,
		),
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

func (d *DockerCompose) ParameterInfo() string {
	return "service name"
}

func (d *DockerCompose) Detect(_ string) *internal.Environment {
	// In docker compose there might be multiple environments, which one is the one?
	return nil
}

// CreateEnvironment create docker compose environment where the service is used for running tools.
func (d *DockerCompose) CreateEnvironment(directory, service string) (internal.Environment, error) {
	internal.Debug("docker-compose.create_environment", "directory", directory, "service", service)

	if len(service) == 0 {
		return nil, fmt.Errorf("service name is required")
	}

	env := dockerComposeEnvironment{
		dockerCompose: *d,
		directory:     directory,
		service:       service,
	}

	return &env, nil
}

type dockerComposeEnvironment struct {
	dockerCompose DockerCompose
	directory     string
	service       string
	autoStop      bool
}

func (d *dockerComposeEnvironment) Id() string {
	return "docker-compose"
}

func (d *dockerComposeEnvironment) Start(ctx context.Context) error {
	internal.Infof("Docker compose start: %v", d.service)

	return internal.Trace(ctx, "docker-compose.start", func() error {
		// It starts all services
		return d.run(ctx, []string{"up", "-d"})
	}, "service", d.service)
}

func (d *dockerComposeEnvironment) IsRunning(ctx context.Context) bool {
	return internal.Trace(ctx, "docker-compose.is_running", func() bool {
		output, err := d.runOutput(ctx, []string{"ps", "--format", "json", d.service})
		if err != nil {
			return false
		}

		var data struct {
			State string `json:"State"`
		}

		if err := json.Unmarshal([]byte(output), &data); err != nil {
			return false
		}

		return data.State == "running"
	}, "service", d.service)
}

func (d *dockerComposeEnvironment) Stop(ctx context.Context) error {
	internal.Infof("Docker compose stop: %v", d.service)

	// It stops all services
	return internal.Trace(ctx, "docker.stop", func() error {
		return d.run(ctx, []string{"stop"})
	}, "service", d.service)
}

func (d *dockerComposeEnvironment) Cleanup(ctx context.Context) error {
	return internal.Trace(ctx, "docker-compose.cleanup", func() error {
		if d.autoStop {
			return d.Stop(ctx)
		}

		return nil
	}, "service", d.service)
}

func (d *dockerComposeEnvironment) FindExecutable(
	ctx context.Context,
	name string,
) (*internal.Executable, error) {
	if err := d.autoStart(ctx); err != nil {
		return nil, err
	}

	return internal.TraceErr(
		ctx,
		"docker-compose.find_executable",
		func() (*internal.Executable, error) {
			// Remove the directory from the name
			if filepath.Base(name) != name {
				searchName, err := filepath.Rel(d.directory, name)
				if err == nil {
					name = searchName
				}
			}

			output, err := d.runOutput(ctx, []string{"exec", d.service, "which", name})
			if err != nil {
				return nil, nil // nolint: nilerr
			}

			return &internal.Executable{
				Path:    output,
				Runtime: d,
			}, nil
		},
		"service",
		d.service,
		"name",
		name,
	)
}

func (d *dockerComposeEnvironment) RunExecutable(
	ctx context.Context,
	options internal.RunOptions,
	path string,
	args []string,
) error {
	if err := d.autoStart(ctx); err != nil {
		return fmt.Errorf("docker compose start failed: %w", err)
	}

	var envArgs []string

	for _, env := range options.Env {
		envArgs = append(envArgs, "-e", env)
	}

	return internal.Trace(ctx, "docker-compose.run", func() error {
		opts := options
		opts.Silent = true

		return d.dockerCompose.Run(
			ctx,
			opts,
			append(
				append(append([]string{"compose", "exec"}, envArgs...), d.service, path),
				args...,
			),
		)
	}, "path", path, "args", args, "env", options.Env)
}

func (d *dockerComposeEnvironment) run(ctx context.Context, args []string) error {
	return d.dockerCompose.Run(
		ctx,
		internal.RunOptions{
			Directory: d.directory,
			Input:     os.Stdin,
			Output:    os.Stdout,
			Error:     os.Stderr,
			Silent:    true,
		},
		append([]string{"compose"}, args...),
	)
}

func (d *dockerComposeEnvironment) runOutput(ctx context.Context, args []string) (string, error) {
	output := bytes.Buffer{}

	err := d.dockerCompose.Run(
		ctx,
		internal.RunOptions{Directory: d.directory, Output: &output, Silent: true},
		append([]string{"compose"}, args...),
	)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(output.String()), nil
}

func (d *dockerComposeEnvironment) autoStart(ctx context.Context) error {
	if !d.IsRunning(ctx) {
		d.autoStop = true

		if err := d.Start(ctx); err != nil {
			return err
		}
	}

	return nil
}
