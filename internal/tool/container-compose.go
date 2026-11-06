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

// A container compose tool to manage tools environment.
type containerCompose struct {
	internal.ExecutableTool
}

func (c *containerCompose) ParameterInfo() string {
	return "service name"
}

func (c *containerCompose) Detect(_ string) *internal.Environment {
	// In container compose there might be multiple environments, which one is the one?
	return nil
}

// CreateEnvironment create container compose environment where the service is used for running tools.
func (c *containerCompose) CreateEnvironment(
	directory, service string,
) (internal.Environment, error) {
	internal.Debug(c.Id()+".create_environment", "directory", directory, "service", service)

	if len(service) == 0 {
		return nil, fmt.Errorf("service name is required")
	}

	env := containerComposeEnvironment{
		executable: c,
		directory:  directory,
		service:    service,
	}

	return &env, nil
}

type containerComposeEnvironment struct {
	executable *containerCompose
	directory  string
	service    string
	autoStop   bool
}

func (c *containerComposeEnvironment) Id() string {
	return c.executable.Id()
}

func (c *containerComposeEnvironment) Start(ctx context.Context) error {
	internal.Infof("%s start: %v", c.Id(), c.service)

	return internal.Trace(ctx, c.Id()+".start", func() error {
		// It starts all services
		return c.run(ctx, []string{"up", "-d"})
	}, "service", c.service)
}

func (c *containerComposeEnvironment) IsRunning(ctx context.Context) bool {
	return internal.Trace(ctx, c.Id()+".is_running", func() bool {
		output, err := c.runOutput(ctx, []string{"ps", "--format", "json"})
		if err != nil {
			return false
		}

		var data []struct {
			State  string            `json:"State"`
			Labels map[string]string `json:"Labels"`
		}

		if err := json.Unmarshal([]byte(output), &data); err != nil {
			return false
		}

		for _, info := range data {
			if name, ok := info.Labels["com.docker.compose.service"]; !ok || name != c.service {
				continue
			}

			if info.State == "running" {
				return true
			}
		}

		return false
	}, "service", c.service)
}

func (c *containerComposeEnvironment) Stop(ctx context.Context) error {
	internal.Infof("%s stop: %v", c.Id(), c.service)

	// It stops all services
	return internal.Trace(ctx, c.Id()+"container.stop", func() error {
		return c.run(ctx, []string{"stop"})
	}, "service", c.service)
}

func (c *containerComposeEnvironment) Cleanup(ctx context.Context) error {
	return internal.Trace(ctx, c.Id()+".cleanup", func() error {
		if c.autoStop {
			return c.Stop(ctx)
		}

		return nil
	}, "service", c.service)
}

func (c *containerComposeEnvironment) FindExecutable(
	ctx context.Context,
	name string,
) (*internal.Executable, error) {
	if err := c.autoStart(ctx); err != nil {
		return nil, err
	}

	return internal.TraceErr(
		ctx,
		c.Id()+".find_executable",
		func() (*internal.Executable, error) {
			// Remove the directory from the name
			if filepath.Base(name) != name {
				searchName, err := filepath.Rel(c.directory, name)
				if err == nil {
					name = searchName
				}
			}

			output, err := c.runOutput(ctx, []string{"exec", c.service, "which", name})
			if err != nil {
				return nil, nil // nolint: nilerr
			}

			return &internal.Executable{
				Path:    output,
				Runtime: c,
			}, nil
		},
		"service",
		c.service,
		"name",
		name,
	)
}

func (c *containerComposeEnvironment) RunExecutable(
	ctx context.Context,
	options internal.RunOptions,
	path string,
	args []string,
) error {
	if err := c.autoStart(ctx); err != nil {
		return fmt.Errorf("%s start failed: %w", c.Id(), err)
	}

	var envArgs []string

	for _, env := range options.Env {
		envArgs = append(envArgs, "-e", env)
	}

	return internal.Trace(ctx, c.Id()+".run", func() error {
		opts := options
		opts.Silent = true

		return c.executable.Run(
			ctx,
			opts,
			append(
				append(append([]string{"compose", "exec"}, envArgs...), c.service, path),
				args...,
			),
		)
	}, "path", path, "args", args, "env", options.Env)
}

func (c *containerComposeEnvironment) run(ctx context.Context, args []string) error {
	return c.executable.Run(
		ctx,
		internal.RunOptions{
			Directory: c.directory,
			Input:     os.Stdin,
			Output:    os.Stdout,
			Error:     os.Stderr,
			Silent:    true,
		},
		append([]string{"compose"}, args...),
	)
}

func (c *containerComposeEnvironment) runOutput(
	ctx context.Context,
	args []string,
) (string, error) {
	output := bytes.Buffer{}

	err := c.executable.Run(
		ctx,
		internal.RunOptions{Directory: c.directory, Output: &output, Silent: true},
		append([]string{"compose"}, args...),
	)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(output.String()), nil
}

func (c *containerComposeEnvironment) autoStart(ctx context.Context) error {
	if !c.IsRunning(ctx) {
		c.autoStop = true

		if err := c.Start(ctx); err != nil {
			return err
		}
	}

	return nil
}
