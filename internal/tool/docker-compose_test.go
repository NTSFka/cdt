package tool

import (
	"cdt/internal"
	"cdt/internal/test"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestDockerCompose_NewDockerCompose_NoExecutable(t *testing.T) {
	tool := NewDockerCompose(func() *internal.Executable { return nil })

	assert.NotNil(t, tool)
	assert.Equal(t, "docker-compose", tool.Id())
	assert.NotEmpty(t, tool.Name())
	assert.NotEmpty(t, tool.Info())
	assert.Equal(t, []string{"dc"}, tool.Aliases())
	assert.False(t, tool.IsAvailable())
}

func TestDockerCompose_NewDockerCompose_WithExecutable(t *testing.T) {
	tool := NewDockerCompose(func() *internal.Executable { return &internal.Executable{Path: "/bin/docker"} })

	assert.NotNil(t, tool)
	assert.Equal(t, "docker-compose", tool.Id())
	assert.True(t, tool.IsAvailable())
}

func TestDockerCompose_DetectDockerCompose_NotFound(t *testing.T) {
	env := test.NewEnvironment(t)
	env.OnFindExecutable("docker").
		Return(nil)

	tool := DetectDockerCompose(context.Background(), env)
	assert.NotNil(t, tool)
	assert.Equal(t, "docker-compose", tool.Id())
	assert.False(t, tool.IsAvailable())

	env.AssertExpectations(t)
}

func TestDockerCompose_DetectDockerCompose_Found(t *testing.T) {
	env := test.NewEnvironment(t)
	env.OnFindExecutable("docker").
		Return(env.NewExecutable("docker"))

	tool := DetectDockerCompose(context.Background(), env)
	assert.NotNil(t, tool)
	assert.Equal(t, "docker-compose", tool.Id())
	assert.True(t, tool.IsAvailable())

	env.AssertExpectations(t)
}

func TestDockerCompose_Detect(t *testing.T) {
	tool := NewDockerCompose(func() *internal.Executable { return &internal.Executable{Path: "docker"} })
	assert.NotNil(t, tool)

	env := tool.Detect(".")
	assert.Nil(t, env)
}

func TestDockerCompose_CreateEnvironment(t *testing.T) {
	tool := NewDockerCompose(func() *internal.Executable { return &internal.Executable{Path: "docker"} })
	assert.NotNil(t, tool)

	env, err := tool.CreateEnvironment(".", "service1")
	require.NoError(t, err)
	assert.NotNil(t, env)

	assert.Equal(t, "docker-compose", env.Id())
}

func TestDockerCompose_CreateEnvironment_NoService(t *testing.T) {
	tool := NewDockerCompose(func() *internal.Executable { return &internal.Executable{Path: "docker"} })
	assert.NotNil(t, tool)

	env, err := tool.CreateEnvironment(".", "")
	require.EqualError(t, err, "service name is required")
	assert.Nil(t, env)
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

	tool := NewDockerCompose(runMock.LazyExecutable("docker"))
	assert.NotNil(t, tool)

	env, err := tool.CreateEnvironment(".", service)
	require.NoError(t, err)
	assert.NotNil(t, env)

	return &runMock, env
}

func TestDockerCompose_Environment_Start(t *testing.T) {
	runMock, env := dockerComposePrepare(t, "service1")

	runMock.OnCall([]string{"compose", "up", "-d"}).
		Return(nil)

	err := env.Start(context.Background())
	require.NoError(t, err)

	runMock.AssertExpectations(t)
}

func TestDockerCompose_Environment_Start_Failed(t *testing.T) {
	runMock, env := dockerComposePrepare(t, "service2")

	runMock.OnCall([]string{"compose", "up", "-d"}).
		Return(errors.New("failed"))

	err := env.Start(context.Background())
	require.EqualError(t, err, "failed")

	runMock.AssertExpectations(t)
}

func TestDockerCompose_Environment_IsRunning_True(t *testing.T) {
	runMock, env := dockerComposePrepare(t, "service3")

	runMock.OnState("service3", true).
		Return(nil)

	result := env.IsRunning(context.Background())
	assert.True(t, result)

	runMock.AssertExpectations(t)
}

func TestDockerCompose_Environment_IsRunning_False(t *testing.T) {
	runMock, env := dockerComposePrepare(t, "service4")

	runMock.OnState("service4", false).
		Return(nil)

	result := env.IsRunning(context.Background())
	assert.False(t, result)

	runMock.AssertExpectations(t)
}

func TestDockerCompose_Environment_IsRunning_InvalidJson(t *testing.T) {
	runMock, env := dockerComposePrepare(t, "service5")

	runMock.OnCallOutput([]string{"compose", "ps", "--format", "json", "service5"}, `{"invalid": "json"`).
		Return(nil)

	result := env.IsRunning(context.Background())
	assert.False(t, result)

	runMock.AssertExpectations(t)
}

func TestDockerCompose_Environment_IsRunning_Failed(t *testing.T) {
	runMock, env := dockerComposePrepare(t, "service6")

	runMock.OnCall([]string{"compose", "ps", "--format", "json", "service6"}).
		Return(errors.New("failed"))

	result := env.IsRunning(context.Background())
	assert.False(t, result)

	runMock.AssertExpectations(t)
}

func TestDockerCompose_Environment_Stop(t *testing.T) {
	runMock, env := dockerComposePrepare(t, "service7")

	runMock.OnCall([]string{"compose", "stop"}).
		Return(nil)

	err := env.Stop(context.Background())
	require.NoError(t, err)

	runMock.AssertExpectations(t)
}

func TestDockerCompose_Environment_Stop_Failed(t *testing.T) {
	runMock, env := dockerComposePrepare(t, "service8")

	runMock.OnCall([]string{"compose", "stop"}).
		Return(errors.New("failed"))

	err := env.Stop(context.Background())
	require.EqualError(t, err, "failed")

	runMock.AssertExpectations(t)
}

func TestDockerCompose_Environment_Cleanup(t *testing.T) {
	runMock, env := dockerComposePrepare(t, "service9")

	err := env.Cleanup(context.Background())
	require.NoError(t, err)

	runMock.AssertExpectations(t)
}

func TestDockerCompose_Environment_Cleanup_Running(t *testing.T) {
	runMock, env := dockerComposePrepare(t, "service10")
	env.(*dockerComposeEnvironment).autoStop = true

	runMock.OnCall([]string{"compose", "stop"}).
		Return(nil)

	err := env.Cleanup(context.Background())
	require.NoError(t, err)

	runMock.AssertExpectations(t)
}

func TestDockerCompose_Environment_Cleanup_Running_Failed(t *testing.T) {
	runMock, env := dockerComposePrepare(t, "service11")
	env.(*dockerComposeEnvironment).autoStop = true

	runMock.OnCall([]string{"compose", "stop"}).
		Return(errors.New("failed"))

	err := env.Cleanup(context.Background())
	require.EqualError(t, err, "failed")

	runMock.AssertExpectations(t)
}

func TestDockerCompose_Environment_FindExecutable(t *testing.T) {
	runMock, env := dockerComposePrepare(t, "service12")

	// Is running
	runMock.OnState("service12", true).
		Return(nil).
		Once()

	runMock.OnCallOutput([]string{"compose", "exec", "service12", "which", "tool1"}, "/usr/bin/tool1").
		Return(nil)

	executable := env.FindExecutable(context.Background(), "tool1")
	require.NotNil(t, executable)
	assert.Equal(t, "/usr/bin/tool1", executable.Path)

	runMock.AssertExpectations(t)
}

func TestDockerCompose_Environment_FindExecutable_Failed(t *testing.T) {
	runMock, env := dockerComposePrepare(t, "service13")

	// Is running
	runMock.OnState("service13", true).
		Return(nil).
		Once()

	runMock.OnCallOutput([]string{"compose", "exec", "service13", "which", "tool1"}, "/usr/bin/tool1").
		Return(errors.New("failed"))

	executable := env.FindExecutable(context.Background(), "tool1")
	assert.Nil(t, executable)

	runMock.AssertExpectations(t)
}

func TestDockerCompose_Environment_FindExecutable_AutoStart(t *testing.T) {
	runMock, env := dockerComposePrepare(t, "service14")

	// Is running
	runMock.OnState("service14", false).
		Return(nil).
		Once()

	// Start
	runMock.OnCall([]string{"compose", "up", "-d"}).
		Return(nil)

	runMock.OnCallOutput([]string{"compose", "exec", "service14", "which", "tool1"}, "/usr/bin/tool1").
		Return(nil)

	executable := env.FindExecutable(context.Background(), "tool1")
	require.NotNil(t, executable)
	assert.Equal(t, "/usr/bin/tool1", executable.Path)

	runMock.AssertExpectations(t)
}

func TestDockerCompose_Environment_FindExecutable_AutoStart_Failed(t *testing.T) {
	runMock, env := dockerComposePrepare(t, "service15")

	// Is running
	runMock.OnState("service15", false).
		Return(nil).
		Once()

	// Start
	runMock.OnCall([]string{"compose", "up", "-d"}).
		Return(errors.New("failed"))

	executable := env.FindExecutable(context.Background(), "tool1")
	assert.Nil(t, executable)

	runMock.AssertExpectations(t)
}

func TestDockerCompose_Environment_RunExecutable(t *testing.T) {
	runMock, env := dockerComposePrepare(t, "service16")

	// Is running
	runMock.OnState("service16", true).
		Return(nil).
		Once()

	runMock.OnCall([]string{"compose", "exec", "service16", "tool1", "arg1", "arg2"}).
		Return(nil).
		Once()

	err := env.RunExecutable(context.Background(), internal.RunOptions{}, "tool1", []string{"arg1", "arg2"})
	require.NoError(t, err)

	runMock.AssertExpectations(t)
}

func TestDockerCompose_Environment_RunExecutable_Failed(t *testing.T) {
	runMock, env := dockerComposePrepare(t, "service17")

	// Is running
	runMock.OnState("service17", true).
		Return(nil).
		Once()

	runMock.OnCall([]string{"compose", "exec", "service17", "tool1", "arg1", "arg2"}).
		Return(errors.New("failed")).
		Once()

	err := env.RunExecutable(context.Background(), internal.RunOptions{}, "tool1", []string{"arg1", "arg2"})
	require.EqualError(t, err, "failed")

	runMock.AssertExpectations(t)
}

func TestDockerCompose_Environment_RunExecutable_AutoStart(t *testing.T) {
	runMock, env := dockerComposePrepare(t, "service18")

	// Is running
	runMock.OnState("service18", false).
		Return(nil).
		Once()

	// Start
	runMock.OnCall([]string{"compose", "up", "-d"}).
		Return(nil).
		Once()

	runMock.OnCall([]string{"compose", "exec", "service18", "tool1", "arg1", "arg2"}).
		Return(nil).
		Once()

	err := env.RunExecutable(context.Background(), internal.RunOptions{}, "tool1", []string{"arg1", "arg2"})
	require.NoError(t, err)

	runMock.AssertExpectations(t)
}

func TestDockerCompose_Environment_RunExecutable_AutoStart_Failed(t *testing.T) {
	runMock, env := dockerComposePrepare(t, "service19")

	// Is running
	runMock.OnState("service19", false).
		Return(nil).
		Once()

	// Start
	runMock.OnCall([]string{"compose", "up", "-d"}).Return(errors.New("failed"))

	err := env.RunExecutable(context.Background(), internal.RunOptions{}, "/usr/bin/tool1", []string{"arg1", "arg2"})
	require.EqualError(t, err, "docker compose start failed: failed")

	runMock.AssertExpectations(t)
}
