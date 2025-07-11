package tool

import (
	"bytes"
	"cdt/internal"
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// A DockerCompose wraps docker compose tool to manage tools environment.
type DockerCompose struct {
	docker *internal.Executable
}

func (d *DockerCompose) Id() string {
	return "docker-compose"
}

func (d *DockerCompose) Name() string {
	return "Docker compose"
}

func (d *DockerCompose) Info() string {
	return "Define and run multi-container applications with Docker"
}

func (d *DockerCompose) IsAvailable() bool {
	return d.docker != nil
}

// NewDockerCompose creates a docker compose tool from a custom docker executable
func NewDockerCompose(docker *internal.Executable) *DockerCompose {
	return &DockerCompose{
		docker: docker,
	}
}

// DetectDockerCompose create a docker compose tool with detected docker executable in the given environment.
func DetectDockerCompose(environment internal.Environment) *DockerCompose {
	return NewDockerCompose(environment.FindExecutable("docker"))
}

// CreateEnvironment create docker compose environment where the service is used for running tools
func (d *DockerCompose) CreateEnvironment(directory, service string) (internal.Environment, error) {
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

func (d *dockerComposeEnvironment) run(ctx context.Context, args []string) error {
	return d.dockerCompose.docker.Run(
		internal.NewRunContext(d.directory),
		append([]string{"compose"}, args...),
	)
}

func (d *dockerComposeEnvironment) runOutput(ctx context.Context, args []string) (string, error) {
	output := bytes.Buffer{}
	err := d.dockerCompose.docker.Run(
		internal.RunContext{Directory: d.directory, Output: &output},
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

func (d *dockerComposeEnvironment) Id() string {
	return "docker-compose"
}

func (d *dockerComposeEnvironment) Start(ctx context.Context) error {
	// It starts all services
	return d.run(ctx, []string{"up", "-d"})
}

func (d *dockerComposeEnvironment) IsRunning(ctx context.Context) bool {
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
}

func (d *dockerComposeEnvironment) Stop(ctx context.Context) error {
	// It stops all services
	return d.run(ctx, []string{"stop"})
}

func (d *dockerComposeEnvironment) Cleanup(ctx context.Context) error {
	if d.autoStop {
		return d.Stop(ctx)
	}

	return nil
}

func (d *dockerComposeEnvironment) FindExecutable(name string) *internal.Executable {
	ctx := context.Background()

	if err := d.autoStart(ctx); err != nil {
		return nil
	}

	output, err := d.runOutput(ctx, []string{"exec", d.service, "which", name})

	if err != nil {
		return nil
	}

	return &internal.Executable{
		Path:    output,
		RunFunc: d.RunExecutable,
	}
}

func (d *dockerComposeEnvironment) RunExecutable(ctx internal.RunContext, path string, args []string) error {
	c := context.Background()

	if err := d.autoStart(c); err != nil {
		return fmt.Errorf("docker compose start failed: %w", err)
	}

	return d.dockerCompose.docker.Run(
		ctx,
		append([]string{"compose", "exec", d.service, path}, args...),
	)
}
