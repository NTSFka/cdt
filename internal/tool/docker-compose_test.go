package tool

import (
	"cdt/internal"
	"cdt/internal/test"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestDockerCompose_NewDockerCompose_NoExecutable(t *testing.T) {
	tool := NewDockerCompose(nil)

	assert.NotNil(t, tool)
	assert.Equal(t, "docker-compose", tool.Id())
	assert.NotEmpty(t, tool.Name())
	assert.NotEmpty(t, tool.Info())
	assert.False(t, tool.IsAvailable())
}

func TestDockerCompose_NewDockerCompose_WithExecutable(t *testing.T) {
	tool := NewDockerCompose(&internal.Executable{Path: "/bin/docker"})

	assert.NotNil(t, tool)
	assert.Equal(t, "docker-compose", tool.Id())
	assert.True(t, tool.IsAvailable())
}

func TestDockerCompose_DetectDockerCompose_NotFound(t *testing.T) {
	environment := test.Environment{}
	environment.OnFindExecutable("docker").
		Return(nil)

	tool := DetectDockerCompose(&environment)
	assert.NotNil(t, tool)
	assert.Equal(t, "docker-compose", tool.Id())
	assert.False(t, tool.IsAvailable())

	environment.AssertExpectations(t)
}

func TestDockerCompose_DetectDockerCompose_Found(t *testing.T) {
	environment := test.Environment{}
	environment.OnFindExecutable("docker").
		Return(environment.MakeExecutable("docker"))

	tool := DetectDockerCompose(&environment)
	assert.NotNil(t, tool)
	assert.Equal(t, "docker-compose", tool.Id())
	assert.True(t, tool.IsAvailable())

	environment.AssertExpectations(t)
}

func TestDockerCompose_CreateEnvironment(t *testing.T) {
	tool := NewDockerCompose(&internal.Executable{Path: "docker"})
	assert.NotNil(t, tool)

	env, err := tool.CreateEnvironment(".", "service1")
	assert.NoError(t, err)
	assert.NotNil(t, env)

	assert.Equal(t, "docker-compose", env.Id())
}

type dockerComposeRunMock struct {
	test.Executable
}

func (m *dockerComposeRunMock) OnCall(args []string) *mock.Call {
	return m.OnRun("docker", args)
}

func (m *dockerComposeRunMock) OnCallOutput(args []string, output string) *mock.Call {
	return m.OnCall(args).
		Run(func(args mock.Arguments) {
			c := args.Get(1).(internal.RunOptions)
			_, _ = c.Output.Write([]byte(output))
		})
}

func (m *dockerComposeRunMock) OnState(service string, result bool) *mock.Call {
	var output string

	if result {
		output = `{"State": "running"}`
	} else {
		output = `{"State": "exited"}`
	}

	return m.OnCallOutput([]string{"compose", "ps", "--format", "json", service}, output)
}

func dockerComposePrepare(t *testing.T, service string) (*dockerComposeRunMock, internal.Environment) {
	runMock := dockerComposeRunMock{}

	tool := NewDockerCompose(runMock.NewExecutable("docker"))
	assert.NotNil(t, tool)

	env, err := tool.CreateEnvironment(".", service)
	assert.NoError(t, err)
	assert.NotNil(t, env)

	return &runMock, env
}

func TestDockerCompose_Environment_Start(t *testing.T) {
	runMock, env := dockerComposePrepare(t, "service1")

	runMock.OnCall([]string{"compose", "up", "-d"}).
		Return(nil)

	err := env.Start(context.Background())
	assert.NoError(t, err)

	runMock.AssertExpectations(t)
}

func TestDockerCompose_Environment_Start_Failed(t *testing.T) {
	runMock, env := dockerComposePrepare(t, "service1")

	runMock.OnCall([]string{"compose", "up", "-d"}).
		Return(errors.New("failed"))

	err := env.Start(context.Background())
	assert.EqualError(t, err, "failed")

	runMock.AssertExpectations(t)
}

func TestDockerCompose_Environment_IsRunning_True(t *testing.T) {
	runMock, env := dockerComposePrepare(t, "service1")

	runMock.OnState("service1", true).
		Return(nil)

	result := env.IsRunning(context.Background())
	assert.True(t, result)

	runMock.AssertExpectations(t)
}

func TestDockerCompose_Environment_IsRunning_False(t *testing.T) {
	runMock, env := dockerComposePrepare(t, "service1")

	runMock.OnState("service1", false).
		Return(nil)

	result := env.IsRunning(context.Background())
	assert.False(t, result)

	runMock.AssertExpectations(t)
}

func TestDockerCompose_Environment_IsRunning_InvalidJson(t *testing.T) {
	runMock, env := dockerComposePrepare(t, "service1")

	runMock.OnCallOutput([]string{"compose", "ps", "--format", "json", "service1"}, `{"invalid": "json"`).
		Return(nil)

	result := env.IsRunning(context.Background())
	assert.False(t, result)

	runMock.AssertExpectations(t)
}

func TestDockerCompose_Environment_IsRunning_Failed(t *testing.T) {
	runMock, env := dockerComposePrepare(t, "service1")

	runMock.OnCall([]string{"compose", "ps", "--format", "json", "service1"}).
		Return(errors.New("failed"))

	result := env.IsRunning(context.Background())
	assert.False(t, result)

	runMock.AssertExpectations(t)
}

func TestDockerCompose_Environment_Stop(t *testing.T) {
	runMock, env := dockerComposePrepare(t, "service1")

	runMock.OnCall([]string{"compose", "stop"}).
		Return(nil)

	err := env.Stop(context.Background())
	assert.NoError(t, err)

	runMock.AssertExpectations(t)
}

func TestDockerCompose_Environment_Stop_Failed(t *testing.T) {
	runMock, env := dockerComposePrepare(t, "service1")

	runMock.OnCall([]string{"compose", "stop"}).
		Return(errors.New("failed"))

	err := env.Stop(context.Background())
	assert.EqualError(t, err, "failed")

	runMock.AssertExpectations(t)
}

func TestDockerCompose_Environment_Cleanup(t *testing.T) {
	runMock, env := dockerComposePrepare(t, "service1")

	err := env.Cleanup(context.Background())
	assert.NoError(t, err)

	runMock.AssertExpectations(t)
}

func TestDockerCompose_Environment_Cleanup_Running(t *testing.T) {
	runMock, env := dockerComposePrepare(t, "service1")
	env.(*dockerComposeEnvironment).autoStop = true

	runMock.OnCall([]string{"compose", "stop"}).
		Return(nil)

	err := env.Cleanup(context.Background())
	assert.NoError(t, err)

	runMock.AssertExpectations(t)
}

func TestDockerCompose_Environment_Cleanup_Running_Failed(t *testing.T) {
	runMock, env := dockerComposePrepare(t, "service1")
	env.(*dockerComposeEnvironment).autoStop = true

	runMock.OnCall([]string{"compose", "stop"}).
		Return(errors.New("failed"))

	err := env.Cleanup(context.Background())
	assert.EqualError(t, err, "failed")

	runMock.AssertExpectations(t)
}

func TestDockerCompose_Environment_FindExecutable(t *testing.T) {
	runMock, env := dockerComposePrepare(t, "service1")

	// Is running
	runMock.OnState("service1", true).
		Return(nil).
		Once()

	runMock.OnCallOutput([]string{"compose", "exec", "service1", "which", "tool1"}, "/usr/bin/tool1").
		Return(nil)

	executable := env.FindExecutable("tool1")
	if assert.NotNil(t, executable) {
		assert.Equal(t, "/usr/bin/tool1", executable.Path)
	}

	runMock.AssertExpectations(t)
}

func TestDockerCompose_Environment_FindExecutable_Failed(t *testing.T) {
	runMock, env := dockerComposePrepare(t, "service1")

	// Is running
	runMock.OnState("service1", true).
		Return(nil).
		Once()

	runMock.OnCallOutput([]string{"compose", "exec", "service1", "which", "tool1"}, "/usr/bin/tool1").
		Return(errors.New("failed"))

	executable := env.FindExecutable("tool1")
	assert.Nil(t, executable)

	runMock.AssertExpectations(t)
}

func TestDockerCompose_Environment_FindExecutable_AutoStart(t *testing.T) {
	runMock, env := dockerComposePrepare(t, "service1")

	// Is running
	runMock.OnState("service1", false).
		Return(nil).
		Once()

	// Start
	runMock.OnCall([]string{"compose", "up", "-d"}).
		Return(nil)

	runMock.OnCallOutput([]string{"compose", "exec", "service1", "which", "tool1"}, "/usr/bin/tool1").
		Return(nil)

	executable := env.FindExecutable("tool1")
	if assert.NotNil(t, executable) {
		assert.Equal(t, "/usr/bin/tool1", executable.Path)
	}

	runMock.AssertExpectations(t)
}

func TestDockerCompose_Environment_FindExecutable_AutoStart_Failed(t *testing.T) {
	runMock, env := dockerComposePrepare(t, "service1")

	// Is running
	runMock.OnState("service1", false).
		Return(nil).
		Once()

	// Start
	runMock.OnCall([]string{"compose", "up", "-d"}).
		Return(errors.New("failed"))

	executable := env.FindExecutable("tool1")
	assert.Nil(t, executable)

	runMock.AssertExpectations(t)
}

func TestDockerCompose_Environment_RunExecutable(t *testing.T) {
	runMock, env := dockerComposePrepare(t, "service1")

	// Is running
	runMock.OnState("service1", true).
		Return(nil).
		Once()

	runMock.OnCall([]string{"compose", "exec", "service1", "tool1", "arg1", "arg2"}).
		Return(nil).
		Once()

	err := env.RunExecutable(context.Background(), internal.RunOptions{}, "tool1", []string{"arg1", "arg2"})
	assert.NoError(t, err)

	runMock.AssertExpectations(t)
}

func TestDockerCompose_Environment_RunExecutable_Failed(t *testing.T) {
	runMock, env := dockerComposePrepare(t, "service1")

	// Is running
	runMock.OnState("service1", true).
		Return(nil).
		Once()

	runMock.OnCall([]string{"compose", "exec", "service1", "tool1", "arg1", "arg2"}).
		Return(errors.New("failed")).
		Once()

	err := env.RunExecutable(context.Background(), internal.RunOptions{}, "tool1", []string{"arg1", "arg2"})
	assert.EqualError(t, err, "failed")

	runMock.AssertExpectations(t)
}

func TestDockerCompose_Environment_RunExecutable_AutoStart(t *testing.T) {
	runMock, env := dockerComposePrepare(t, "service1")

	// Is running
	runMock.OnState("service1", false).
		Return(nil).
		Once()

	// Start
	runMock.OnCall([]string{"compose", "up", "-d"}).
		Return(nil).
		Once()

	runMock.OnCall([]string{"compose", "exec", "service1", "tool1", "arg1", "arg2"}).
		Return(nil).
		Once()

	err := env.RunExecutable(context.Background(), internal.RunOptions{}, "tool1", []string{"arg1", "arg2"})
	assert.NoError(t, err)

	runMock.AssertExpectations(t)
}

func TestDockerCompose_Environment_RunExecutable_AutoStart_Failed(t *testing.T) {
	runMock, env := dockerComposePrepare(t, "service1")

	// Is running
	runMock.OnState("service1", false).
		Return(nil).
		Once()

	// Start
	runMock.OnCall([]string{"compose", "up", "-d"}).Return(errors.New("failed"))

	err := env.RunExecutable(context.Background(), internal.RunOptions{}, "/usr/bin/tool1", []string{"arg1", "arg2"})
	assert.EqualError(t, err, "docker compose start failed: failed")

	runMock.AssertExpectations(t)
}
