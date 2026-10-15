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

const IdPodmanCompose = "podman-compose"

// A PodmanCompose wraps podman compose tool to manage tools environment.
type PodmanCompose struct {
	internal.ExecutableTool
}

// NewPodmanCompose creates a podman compose tool from a custom podman executable.
func NewPodmanCompose(detect internal.ExecutableToolDetectFunc) *PodmanCompose {
	return &PodmanCompose{
		internal.MakeExecutableTool(
			IdPodmanCompose,
			"Podman compose",
			"Define and run multi-container applications with Podman",
			internal.Tags{internal.ToolTagEnvironment},
			detect,
		),
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

func (d *PodmanCompose) ParameterInfo() string {
	return "service name"
}

func (d *PodmanCompose) Detect(_ string) *internal.Environment {
	// In podman compose there might be multiple environments, which one is the one?
	return nil
}

// CreateEnvironment create podman compose environment where the service is used for running tools.
func (d *PodmanCompose) CreateEnvironment(directory, service string) (internal.Environment, error) {
	internal.Debug("podman-compose.create_environment", "directory", directory, "service", service)

	if len(service) == 0 {
		return nil, fmt.Errorf("service name is required")
	}

	env := podmanComposeEnvironment{
		podmanCompose: *d,
		directory:     directory,
		service:       service,
	}

	return &env, nil
}

type podmanComposeEnvironment struct {
	podmanCompose PodmanCompose
	directory     string
	service       string
	autoStop      bool
}

func (d *podmanComposeEnvironment) Id() string {
	return "podman-compose"
}

func (d *podmanComposeEnvironment) Start(ctx context.Context) error {
	internal.Infof("Podman compose start: %v", d.service)

	return internal.Trace(ctx, "podman-compose.start", func() error {
		// It starts all services
		return d.run(ctx, []string{"up", "-d"})
	}, "service", d.service)
}

func (d *podmanComposeEnvironment) IsRunning(ctx context.Context) bool {
	return internal.Trace(ctx, "podman-compose.is_running", func() bool {
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

func (d *podmanComposeEnvironment) Stop(ctx context.Context) error {
	internal.Infof("Podman compose stop: %v", d.service)

	// It stops all services
	return internal.Trace(ctx, "podman.stop", func() error {
		return d.run(ctx, []string{"stop"})
	}, "service", d.service)
}

func (d *podmanComposeEnvironment) Cleanup(ctx context.Context) error {
	return internal.Trace(ctx, "podman-compose.cleanup", func() error {
		if d.autoStop {
			return d.Stop(ctx)
		}

		return nil
	}, "service", d.service)
}

func (d *podmanComposeEnvironment) FindExecutable(
	ctx context.Context,
	name string,
) (*internal.Executable, error) {
	if err := d.autoStart(ctx); err != nil {
		return nil, err
	}

	return internal.TraceErr(
		ctx,
		"podman-compose.find_executable",
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

func (d *podmanComposeEnvironment) RunExecutable(
	ctx context.Context,
	options internal.RunOptions,
	path string,
	args []string,
) error {
	if err := d.autoStart(ctx); err != nil {
		return fmt.Errorf("podman compose start failed: %w", err)
	}

	var envArgs []string

	for _, env := range options.Env {
		envArgs = append(envArgs, "-e", env)
	}

	return internal.Trace(ctx, "podman-compose.run", func() error {
		opts := options
		opts.Silent = true

		return d.podmanCompose.Run(
			ctx,
			opts,
			append(
				append(append([]string{"compose", "exec"}, envArgs...), d.service, path),
				args...,
			),
		)
	}, "path", path, "args", args, "env", options.Env)
}

func (d *podmanComposeEnvironment) run(ctx context.Context, args []string) error {
	return d.podmanCompose.Run(
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

func (d *podmanComposeEnvironment) runOutput(ctx context.Context, args []string) (string, error) {
	output := bytes.Buffer{}

	err := d.podmanCompose.Run(
		ctx,
		internal.RunOptions{Directory: d.directory, Output: &output, Silent: true},
		append([]string{"compose"}, args...),
	)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(output.String()), nil
}

func (d *podmanComposeEnvironment) autoStart(ctx context.Context) error {
	if !d.IsRunning(ctx) {
		d.autoStop = true

		if err := d.Start(ctx); err != nil {
			return err
		}
	}

	return nil
}
