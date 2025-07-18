package tool

import (
	"bytes"
	"cdt/internal"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"path/filepath"
	"strings"
)

// A Docker is a tool that wraps docker executable.
type Docker struct {
	internal.ExecutableTool
}

// NewDocker creates a docker tool from a custom executable
func NewDocker(detect func() *internal.Executable) *Docker {
	return &Docker{internal.MakeExecutableTool(
		"docker",
		"Docker",
		"Docker image and container command line interface.",
		detect,
	)}
}

// DetectDocker create a docker tool with a detected docker executable in the given environment.
func DetectDocker(environment internal.Environment) *Docker {
	return NewDocker(func() *internal.Executable {
		return environment.FindExecutable("docker")
	})
}

// CreateEnvironment create docker environment where the service is used for running tools
func (d *Docker) CreateEnvironment(directory, image string) (internal.Environment, error) {
	slog.Debug("docker.create_environment", "directory", directory, "image", image)

	env := dockerEnvironment{
		directory: directory,
		imageName: image,
		docker:    d,
	}

	return &env, nil
}

type dockerEnvironment struct {
	directory   string
	imageName   string
	docker      *Docker
	containerId string
	autoStop    bool
}

func (d *dockerEnvironment) Id() string {
	return "docker"
}

func (d *dockerEnvironment) runOutput(ctx context.Context, args []string) (string, error) {
	output := bytes.Buffer{}
	err := d.docker.Run(
		ctx,
		internal.RunOptions{Directory: d.directory, Output: &output},
		args,
	)

	if err != nil {
		return "", fmt.Errorf("docker run failed: %w", err)
	}

	return strings.TrimSpace(output.String()), nil
}

func (d *dockerEnvironment) autoStart(ctx context.Context) error {
	if !d.IsRunning(ctx) {
		d.autoStop = true

		if err := d.Start(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (d *dockerEnvironment) Start(ctx context.Context) error {
	if d.IsRunning(ctx) {
		return nil
	}

	absPath, err := filepath.Abs(d.directory)
	internal.Assert(err == nil, "failed to determine absolute path")

	slog.Info("Docker start", "image", d.imageName)

	return internal.Trace(ctx, "docker.start", func() error {
		output, err := d.runOutput(ctx, []string{
			"run", "--rm", "-d",
			"-v", fmt.Sprintf("%s:/work", absPath),
			"-w", "/work",
			d.imageName,
			// FIXME: only linux
			"/bin/bash", "-c", "trap : TERM INT; sleep infinity & wait",
		})

		if err == nil {
			d.containerId = output
		}

		return err
	}, "image", d.imageName)
}

func (d *dockerEnvironment) IsRunning(ctx context.Context) bool {
	// No container name specified = not running
	if len(d.containerId) == 0 {
		return false
	}

	return internal.Trace(ctx, "docker.is_running", func() bool {
		output, err := d.runOutput(ctx, []string{"inspect", "--format", "json", d.containerId})

		if err != nil {
			return false
		}

		var data []struct {
			State struct {
				Running bool `json:"Running"`
			} `json:"State"`
		}
		if err := json.Unmarshal([]byte(output), &data); err != nil {
			return false
		}

		if len(data) == 0 {
			return false
		}

		return data[0].State.Running
	})
}

func (d *dockerEnvironment) Stop(ctx context.Context) error {
	internal.Assert(d.containerId != "", "container ID is not set")

	slog.Info("Docker stop")

	return internal.Trace(ctx, "docker.stop", func() error {
		_, err := d.runOutput(ctx, []string{"stop", d.containerId})
		return err
	}, "container", d.containerId)
}

func (d *dockerEnvironment) Cleanup(ctx context.Context) error {
	return internal.Trace(ctx, "docker.cleanup", func() error {
		if d.autoStop {
			return d.Stop(ctx)
		}

		return nil
	}, "container", d.containerId)
}

func (d *dockerEnvironment) FindExecutable(name string) *internal.Executable {
	ctx := context.Background()

	if err := d.autoStart(ctx); err != nil {
		return nil
	}

	internal.Assert(d.containerId != "", "container ID is not set")

	return internal.Trace(ctx, "docker.find_executable", func() *internal.Executable {
		output, err := d.runOutput(ctx, []string{"exec", d.containerId, "which", name})

		if err != nil {
			return nil
		}

		return &internal.Executable{
			Path:    output,
			RunFunc: d.RunExecutable,
		}
	}, "container", d.containerId, "name", name)
}

func (d *dockerEnvironment) RunExecutable(ctx context.Context, options internal.RunOptions, path string, args []string) error {
	c := context.Background()

	if err := d.autoStart(c); err != nil {
		return err
	}

	internal.Assert(d.containerId != "", "container ID is not set")

	return internal.Trace(ctx, "docker.run", func() error {
		return d.docker.Run(ctx, options, append([]string{"exec", d.containerId, path}, args...))
	}, "container", d.containerId, "path", path, "args", args)
}
