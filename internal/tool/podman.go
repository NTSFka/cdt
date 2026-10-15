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

const IdPodman = "podman"

// A Podman is a tool that wraps podman executable.
type Podman struct {
	internal.ExecutableTool
}

// NewPodman creates a podman tool from a custom executable.
func NewPodman(detect internal.ExecutableToolDetectFunc) *Podman {
	return &Podman{internal.MakeExecutableTool(
		IdPodman,
		"Podman",
		"Manage containers, pods, and images with Podman.",
		internal.Tags{internal.ToolTagEnvironment},
		detect,
	)}
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

func (d *Podman) ParameterInfo() string {
	return "image name"
}

func (d *Podman) Detect(_ string) *internal.Environment {
	// No way to detect an environment (yet)
	return nil
}

// CreateEnvironment create podman environment where the service is used for running tools.
func (d *Podman) CreateEnvironment(directory, image string) (internal.Environment, error) {
	internal.Debug("podman.create_environment", "directory", directory, "image", image)

	if len(image) == 0 {
		return nil, errors.New("podman image name is required")
	}

	env := podmanEnvironment{
		directory: directory,
		imageName: image,
		podman:    d,
	}

	return &env, nil
}

type podmanEnvironment struct {
	directory   string
	imageName   string
	podman      *Podman
	containerId string
	autoStop    bool
}

func (d *podmanEnvironment) Id() string {
	return "podman"
}

func (d *podmanEnvironment) Start(ctx context.Context) error {
	if d.IsRunning(ctx) {
		return nil
	}

	absPath, err := filepath.Abs(d.directory)
	internal.Assert(err == nil, "failed to determine absolute path")

	internal.Infof("Podman start: %v", d.imageName)

	return internal.Trace(ctx, "podman.start", func() error {
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

func (d *podmanEnvironment) IsRunning(ctx context.Context) bool {
	// No container name specified = not running
	if len(d.containerId) == 0 {
		return false
	}

	return internal.Trace(ctx, "podman.is_running", func() bool {
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

func (d *podmanEnvironment) Stop(ctx context.Context) error {
	if d.containerId == "" {
		return nil
	}

	internal.Infof("Podman stop")

	return internal.Trace(ctx, "podman.stop", func() error {
		_, err := d.runOutput(ctx, []string{"stop", d.containerId})

		return err
	}, "container", d.containerId)
}

func (d *podmanEnvironment) Cleanup(ctx context.Context) error {
	return internal.Trace(ctx, "podman.cleanup", func() error {
		if d.autoStop {
			return d.Stop(ctx)
		}

		return nil
	}, "container", d.containerId)
}

func (d *podmanEnvironment) FindExecutable(
	ctx context.Context,
	name string,
) (*internal.Executable, error) {
	if err := d.autoStart(ctx); err != nil {
		return nil, err
	}

	internal.Assert(d.containerId != "", "container ID is not set")

	return internal.TraceErr(ctx, "podman.find_executable", func() (*internal.Executable, error) {
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

func (d *podmanEnvironment) RunExecutable(
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

	return internal.Trace(ctx, "podman.run", func() error {
		opts := options
		opts.Silent = true

		return d.podman.Run(
			ctx,
			opts,
			append(append(append([]string{"exec"}, envArgs...), d.containerId, path), args...),
		)
	}, "container", d.containerId, "path", path, "args", args, "env", options.Env)
}

func (d *podmanEnvironment) runOutput(ctx context.Context, args []string) (string, error) {
	output := bytes.Buffer{}

	err := d.podman.Run(
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
		return "", fmt.Errorf("podman run failed: %w", err)
	}

	return strings.TrimSpace(output.String()), nil
}

func (d *podmanEnvironment) autoStart(ctx context.Context) error {
	if !d.IsRunning(ctx) {
		d.autoStop = true

		if err := d.Start(ctx); err != nil {
			return err
		}
	}

	return nil
}
