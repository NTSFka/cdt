package tool

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"cdt/internal"
)

const IdDocker = "docker"

// A Docker is a tool that wraps docker executable.
type Docker struct {
	internal.ExecutableTool
}

// NewDocker creates a docker tool from a custom executable.
func NewDocker(detect internal.ExecutableToolDetectFunc) *Docker {
	return &Docker{internal.MakeExecutableTool(
		IdDocker,
		"Docker",
		"Docker image and container command line interface.",
		internal.Tags{internal.ToolTagEnvironment},
		detect,
	)}
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

func (d *Docker) ParameterInfo() string {
	return "image name"
}

func (d *Docker) Detect(_ string) *internal.Environment {
	// No way to detect an environment (yet)
	return nil
}

// CreateEnvironment create docker environment where the service is used for running tools.
func (d *Docker) CreateEnvironment(directory, image string) (internal.Environment, error) {
	internal.Debug("docker.create_environment", "directory", directory, "image", image)

	if len(image) == 0 {
		return nil, errors.New("docker image name is required")
	}

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

func (d *dockerEnvironment) Start(ctx context.Context) error {
	if d.IsRunning(ctx) {
		return nil
	}

	absPath, err := filepath.Abs(d.directory)
	internal.Assert(err == nil, "failed to determine absolute path")

	internal.Infof("Docker start: %v", d.imageName)

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
	if d.containerId == "" {
		return nil
	}

	internal.Infof("Docker stop")

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

func (d *dockerEnvironment) FindExecutable(
	ctx context.Context,
	name string,
) (*internal.Executable, error) {
	if err := d.autoStart(ctx); err != nil {
		return nil, err
	}

	internal.Assert(d.containerId != "", "container ID is not set")

	return internal.TraceErr(ctx, "docker.find_executable", func() (*internal.Executable, error) {
		// Remove the directory from the name
		if filepath.Base(name) != name {
			searchName, err := filepath.Rel(d.directory, name)
			if err == nil {
				name = searchName
			}
		}

		output, err := d.runOutput(ctx, []string{"exec", d.containerId, "which", name})
		if err != nil {
			return nil, nil // nolint: nilerr
		}

		return &internal.Executable{
			Path:    output,
			Runtime: d,
		}, nil
	}, "container", d.containerId, "name", name)
}

func (d *dockerEnvironment) RunExecutable(
	ctx context.Context,
	options internal.RunOptions,
	path string,
	args []string,
) error {
	if err := d.autoStart(ctx); err != nil {
		return err
	}

	internal.Assert(d.containerId != "", "container ID is not set")

	var envArgs []string

	for _, env := range options.Env {
		envArgs = append(envArgs, "-e", env)
	}

	return internal.Trace(ctx, "docker.run", func() error {
		opts := options
		opts.Silent = true

		return d.docker.Run(
			ctx,
			opts,
			append(append(append([]string{"exec"}, envArgs...), d.containerId, path), args...),
		)
	}, "container", d.containerId, "path", path, "args", args, "env", options.Env)
}

func (d *dockerEnvironment) runOutput(ctx context.Context, args []string) (string, error) {
	output := bytes.Buffer{}

	err := d.docker.Run(
		ctx,
		internal.RunOptions{
			Directory: d.directory,
			Output:    &output,
			Error:     os.Stderr,
			Silent:    true,
		},
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
