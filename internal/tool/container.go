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

// Base type for container based tools (container, Podman).
type container struct {
	internal.ExecutableTool
}

func (c *container) ParameterInfo() string {
	return "image name"
}

func (c *container) Detect(_ string) *internal.Environment {
	// No way to detect an environment (yet)
	return nil
}

// CreateEnvironment create container environment where the service is used for running tools.
func (c *container) CreateEnvironment(directory, image string) (internal.Environment, error) {
	internal.Debug(c.Id()+".create_environment", "directory", directory, "image", image)

	if len(image) == 0 {
		return nil, errors.New(c.Id() + " image name is required")
	}

	env := containerEnvironment{
		directory:  directory,
		imageName:  image,
		executable: c,
	}

	return &env, nil
}

type containerEnvironment struct {
	directory   string
	imageName   string
	executable  *container
	containerId string
	autoStop    bool
}

func (c *containerEnvironment) Id() string {
	return c.executable.Id()
}

func (c *containerEnvironment) Start(ctx context.Context) error {
	if c.IsRunning(ctx) {
		return nil
	}

	absPath, err := filepath.Abs(c.directory)
	internal.Assert(err == nil, "failed to determine absolute path")

	internal.Infof("%s start: %v", c.Id(), c.imageName)

	return internal.Trace(ctx, c.Id()+".start", func() error {
		output, err := c.runOutput(ctx, []string{
			"run", "--rm", "-d",
			"-v", fmt.Sprintf("%s:/work", absPath),
			"-w", "/work",
			c.imageName,
			// FIXME: only linux
			"/bin/bash", "-c", "trap : TERM INT; sleep infinity & wait",
		})
		if err == nil {
			c.containerId = output
		}

		return err
	}, "image", c.imageName)
}

func (c *containerEnvironment) IsRunning(ctx context.Context) bool {
	// No container name specified = not running
	if len(c.containerId) == 0 {
		return false
	}

	return internal.Trace(ctx, c.Id()+".is_running", func() bool {
		output, err := c.runOutput(ctx, []string{"inspect", "--format", "json", c.containerId})
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

func (c *containerEnvironment) Stop(ctx context.Context) error {
	if c.containerId == "" {
		return nil
	}

	internal.Infof("%s stop", c.Id())

	return internal.Trace(ctx, c.Id()+".stop", func() error {
		_, err := c.runOutput(ctx, []string{"stop", c.containerId})

		return err
	}, "container", c.containerId)
}

func (c *containerEnvironment) Cleanup(ctx context.Context) error {
	return internal.Trace(ctx, c.Id()+".cleanup", func() error {
		if c.autoStop {
			return c.Stop(ctx)
		}

		return nil
	}, "container", c.containerId)
}

func (c *containerEnvironment) FindExecutable(
	ctx context.Context,
	name string,
) (*internal.Executable, error) {
	if err := c.autoStart(ctx); err != nil {
		return nil, err
	}

	internal.Assert(c.containerId != "", "container ID is not set")

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

			output, err := c.runOutput(ctx, []string{"exec", c.containerId, "which", name})
			if err != nil {
				return nil, nil // nolint: nilerr
			}

			return &internal.Executable{
				Path:    output,
				Runtime: c,
			}, nil
		},
		"container",
		c.containerId,
		"name",
		name,
	)
}

func (c *containerEnvironment) RunExecutable(
	ctx context.Context,
	options internal.RunOptions,
	path string,
	args []string,
) error {
	if err := c.autoStart(ctx); err != nil {
		return err
	}

	internal.Assert(c.containerId != "", "container ID is not set")

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
			append(append(append([]string{"exec"}, envArgs...), c.containerId, path), args...),
		)
	}, "container", c.containerId, "path", path, "args", args, "env", options.Env)
}

func (c *containerEnvironment) runOutput(ctx context.Context, args []string) (string, error) {
	output := bytes.Buffer{}

	err := c.executable.Run(
		ctx,
		internal.RunOptions{
			Directory: c.directory,
			Output:    &output,
			Error:     os.Stderr,
			Silent:    true,
		},
		args,
	)
	if err != nil {
		return "", fmt.Errorf(c.Id()+" run failed: %w", err)
	}

	return strings.TrimSpace(output.String()), nil
}

func (c *containerEnvironment) autoStart(ctx context.Context) error {
	if !c.IsRunning(ctx) {
		c.autoStop = true

		if err := c.Start(ctx); err != nil {
			return err
		}
	}

	return nil
}
