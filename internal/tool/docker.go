package tool

import (
	"bytes"
	"cdt/internal"
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"path/filepath"
	"strings"
)

// A Docker is a tool that wraps docker executable.
type Docker struct {
	client *client.Client
}

func (d *Docker) Id() string {
	return "docker"
}

func (d *Docker) Name() string {
	return "Docker"
}

func (d *Docker) Info() string {
	return "Docker image and container command line interface."
}

func (d *Docker) IsAvailable() bool {
	return d.client != nil
}

// NewDocker creates a docker tool from a custom executable
func NewDocker(client *client.Client) *Docker {
	return &Docker{
		client: client,
	}
}

// DetectDocker create docker tool can be used in the project
func DetectDocker() *Docker {
	c, err := client.NewClientWithOpts(client.FromEnv)

	if err != nil {
		return nil
	}

	return NewDocker(c)
}

func (d *Docker) CreateEnvironment(directory, image string) (internal.Environment, error) {
	env := DockerEnvironment{
		Directory: directory,
		Image:     image,
		Docker:    d,
	}

	return &env, nil
}

type DockerEnvironment struct {
	Directory     string
	Image         string
	Docker        *Docker
	ContainerName string
	ContainerId   string
	AutoStop      bool
}

func (d *DockerEnvironment) Id() string {
	return "docker"
}

func (d *DockerEnvironment) Start(ctx context.Context) error {
	if d.IsRunning(ctx) {
		return nil
	}

	absPath, err := filepath.Abs(d.Directory)

	if err != nil {
		return err
	}

	d.ContainerName = fmt.Sprintf("cdt-%s-%s", filepath.Base(absPath), strings.ReplaceAll(d.Image, ":", "_"))

	resp, err := d.Docker.client.ContainerCreate(ctx, &container.Config{
		Image: d.Image,
		// FIXME: linux only
		Cmd:        []string{"/bin/bash", "-c", "trap : TERM INT; sleep infinity & wait"},
		WorkingDir: "/work",
	}, &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: absPath,
				Target: "/work",
			},
		},
	}, nil, nil, d.ContainerName)

	if err != nil {
		return err
	}

	d.ContainerId = resp.ID

	return d.Docker.client.ContainerStart(ctx, d.ContainerId, container.StartOptions{})
}

func (d *DockerEnvironment) IsRunning(ctx context.Context) bool {
	resp, err := d.Docker.client.ContainerInspect(ctx, d.ContainerId)

	if err != nil {
		return false
	}

	return resp.State.Running
}

func (d *DockerEnvironment) Stop(ctx context.Context) error {
	if !d.IsRunning(ctx) {
		return nil
	}

	err := d.Docker.client.ContainerStop(ctx, d.ContainerId, container.StopOptions{})
	if err != nil {
		return err
	}

	return d.Docker.client.ContainerRemove(ctx, d.ContainerId, container.RemoveOptions{})
}

func (d *DockerEnvironment) Cleanup(ctx context.Context) error {
	if d.AutoStop {
		return d.Stop(ctx)
	}

	return nil
}

func (d *DockerEnvironment) FindExecutable(name string) *internal.Executable {
	ctx := context.Background()

	if !d.IsRunning(ctx) {
		d.AutoStop = true

		if err := d.Start(ctx); err != nil {
			return nil
		}
	}

	execResp, err := d.Docker.client.ContainerExecCreate(ctx, d.ContainerId, container.ExecOptions{
		Cmd:          []string{"which", name},
		AttachStdout: true,
	})

	if err != nil {
		fmt.Printf("Error executing the docker container: %v\n", err)
		return nil
	}

	attachResp, err := d.Docker.client.ContainerExecAttach(ctx, execResp.ID, container.ExecAttachOptions{})

	if err != nil {
		fmt.Printf("Error attaching to docker container: %v\n", err)
		return nil
	}

	buffer := bytes.Buffer{}
	_, _ = stdcopy.StdCopy(&buffer, nil, attachResp.Reader)

	if buffer.Available() == 0 {
		return nil
	}

	return &internal.Executable{
		Path:    strings.TrimSpace(buffer.String()),
		RunFunc: d.RunExecutable,
	}
}

func (d *DockerEnvironment) RunExecutable(ctx internal.RunContext, path string, args []string) error {
	c := context.Background()

	if !d.IsRunning(c) {
		d.AutoStop = true

		if err := d.Start(c); err != nil {
			return err
		}
	}

	execResp, err := d.Docker.client.ContainerExecCreate(c, d.ContainerId, container.ExecOptions{
		Cmd:          append([]string{path}, args...),
		AttachStdout: true,
		AttachStderr: true,
	})

	if err != nil {
		return err
	}

	attachResp, err := d.Docker.client.ContainerExecAttach(c, execResp.ID, container.ExecAttachOptions{})

	if err != nil {
		return err
	}

	_, _ = stdcopy.StdCopy(ctx.Output, ctx.Error, attachResp.Reader)

	return nil
}
