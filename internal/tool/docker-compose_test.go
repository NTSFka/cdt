package tool_test

import (
	"errors"
	"testing"

	"cdt/internal"
	"cdt/internal/test"
	"cdt/internal/tool"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestDockerCompose_NewDockerCompose_NoExecutable(t *testing.T) {
	dockerCompose := tool.NewDockerCompose(test.LazyExecutableNil)

	assert.NotNil(t, dockerCompose)
	assert.Equal(t, "docker-compose", dockerCompose.Id())
	assert.NotEmpty(t, dockerCompose.Name())
	assert.NotEmpty(t, dockerCompose.Info())
	assert.Equal(t, []string{"dc"}, dockerCompose.Aliases())
	assert.False(t, dockerCompose.IsAvailable())
}

func TestDockerCompose_NewDockerCompose_WithExecutable(t *testing.T) {
	dockerCompose := tool.NewDockerCompose(test.LazyExecutable("/bin/docker"))

	assert.NotNil(t, dockerCompose)
	assert.Equal(t, "docker-compose", dockerCompose.Id())
	assert.True(t, dockerCompose.IsAvailable())
}

func TestDockerCompose_DetectDockerCompose_NotFound(t *testing.T) {
	env := test.NewEnvironment(t)
	env.OnFindExecutable("docker").
		Return(nil, nil)

	dockerCompose := tool.DetectDockerCompose(t.Context(), tool.DetectOptions{Environment: env})
	assert.NotNil(t, dockerCompose)
	assert.Equal(t, "docker-compose", dockerCompose.Id())
	assert.False(t, dockerCompose.IsAvailable())

	env.AssertExpectations(t)
}

func TestDockerCompose_DetectDockerCompose_Found(t *testing.T) {
	env := test.NewEnvironment(t)
	env.OnFindExecutable("docker").
		Return(env.NewExecutable("docker"), nil, nil)

	dockerCompose := tool.DetectDockerCompose(t.Context(), tool.DetectOptions{Environment: env})
	assert.NotNil(t, dockerCompose)
	assert.Equal(t, "docker-compose", dockerCompose.Id())
	assert.True(t, dockerCompose.IsAvailable())

	env.AssertExpectations(t)
}

func TestDockerCompose_DetectDockerCompose_Config(t *testing.T) {
	env := test.NewEnvironment(t)
	env.OnFindExecutable("docker-2").
		Return(env.NewExecutable("docker"), nil, nil)

	dockerCompose := tool.DetectDockerCompose(t.Context(), tool.DetectOptions{
		Environment: env,
		ToolsPaths:  map[string]string{"docker-compose": "docker-2"},
	})
	assert.NotNil(t, dockerCompose)
	assert.Equal(t, "docker-compose", dockerCompose.Id())
	assert.True(t, dockerCompose.IsAvailable())

	env.AssertExpectations(t)
}

func TestDockerCompose_Detect(t *testing.T) {
	dockerCompose := tool.NewDockerCompose(test.LazyExecutable("docker"))
	assert.NotNil(t, dockerCompose)

	env := dockerCompose.Detect(".")
	assert.Nil(t, env)
}

func TestDockerCompose_CreateEnvironment(t *testing.T) {
	dockerCompose := tool.NewDockerCompose(test.LazyExecutable("docker"))
	assert.NotNil(t, dockerCompose)

	env, err := dockerCompose.CreateEnvironment(".", "service1")
	require.NoError(t, err)
	assert.NotNil(t, env)

	assert.Equal(t, "docker-compose", env.Id())
}

func TestDockerCompose_CreateEnvironment_NoService(t *testing.T) {
	dockerCompose := tool.NewDockerCompose(test.LazyExecutable("docker"))
	assert.NotNil(t, dockerCompose)

	env, err := dockerCompose.CreateEnvironment(".", "")
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

func dockerComposePrepare(
	t *testing.T,
	service string,
) (*dockerComposeRunMock, internal.Environment) {
	runMock := dockerComposeRunMock{}

	dockerCompose := tool.NewDockerCompose(runMock.LazyExecutable("docker"))
	assert.NotNil(t, dockerCompose)

	env, err := dockerCompose.CreateEnvironment(".", service)
	require.NoError(t, err)
	assert.NotNil(t, env)

	return &runMock, env
}

func TestDockerCompose_Environment_Start(t *testing.T) {
	runMock, env := dockerComposePrepare(t, "service1")

	runMock.OnCall([]string{"compose", "up", "-d"}).
		Return(nil)

	err := env.Start(t.Context())
	require.NoError(t, err)

	runMock.AssertExpectations(t)
}

func TestDockerCompose_Environment_Start_Failed(t *testing.T) {
	runMock, env := dockerComposePrepare(t, "service2")

	runMock.OnCall([]string{"compose", "up", "-d"}).
		Return(errors.New("failed"))

	err := env.Start(t.Context())
	require.EqualError(t, err, "failed")

	runMock.AssertExpectations(t)
}

func TestDockerCompose_Environment_IsRunning_True(t *testing.T) {
	runMock, env := dockerComposePrepare(t, "service3")

	runMock.OnState("service3", true).
		Return(nil)

	result := env.IsRunning(t.Context())
	assert.True(t, result)

	runMock.AssertExpectations(t)
}

func TestDockerCompose_Environment_IsRunning_False(t *testing.T) {
	runMock, env := dockerComposePrepare(t, "service4")

	runMock.OnState("service4", false).
		Return(nil)

	result := env.IsRunning(t.Context())
	assert.False(t, result)

	runMock.AssertExpectations(t)
}

func TestDockerCompose_Environment_IsRunning_InvalidJson(t *testing.T) {
	runMock, env := dockerComposePrepare(t, "service5")

	runMock.OnCallOutput([]string{"compose", "ps", "--format", "json", "service5"}, `{"invalid": "json"`).
		Return(nil)

	result := env.IsRunning(t.Context())
	assert.False(t, result)

	runMock.AssertExpectations(t)
}

func TestDockerCompose_Environment_IsRunning_Failed(t *testing.T) {
	runMock, env := dockerComposePrepare(t, "service6")

	runMock.OnCall([]string{"compose", "ps", "--format", "json", "service6"}).
		Return(errors.New("failed"))

	result := env.IsRunning(t.Context())
	assert.False(t, result)

	runMock.AssertExpectations(t)
}

func TestDockerCompose_Environment_Stop(t *testing.T) {
	runMock, env := dockerComposePrepare(t, "service7")

	runMock.OnCall([]string{"compose", "stop"}).
		Return(nil)

	err := env.Stop(t.Context())
	require.NoError(t, err)

	runMock.AssertExpectations(t)
}

func TestDockerCompose_Environment_Stop_Failed(t *testing.T) {
	runMock, env := dockerComposePrepare(t, "service8")

	runMock.OnCall([]string{"compose", "stop"}).
		Return(errors.New("failed"))

	err := env.Stop(t.Context())
	require.EqualError(t, err, "failed")

	runMock.AssertExpectations(t)
}

func TestDockerCompose_Environment_Cleanup(t *testing.T) {
	runMock, env := dockerComposePrepare(t, "service9")

	err := env.Cleanup(t.Context())
	require.NoError(t, err)

	runMock.AssertExpectations(t)
}

func TestDockerCompose_Environment_Cleanup_Running(t *testing.T) {
	runMock, env := dockerComposePrepare(t, "service10")

	// Force autoStart
	runMock.OnState("service10", false).Return(nil)
	runMock.OnCall([]string{"compose", "up", "-d"}).Return(nil)
	runMock.OnCall([]string{"compose", "exec", "service10", "test"}).Return(nil)

	err := env.RunExecutable(t.Context(), internal.RunOptions{}, "test", nil)
	require.NoError(t, err)

	runMock.OnCall([]string{"compose", "stop"}).
		Return(nil)

	err = env.Cleanup(t.Context())
	require.NoError(t, err)

	runMock.AssertExpectations(t)
}

func TestDockerCompose_Environment_Cleanup_Running_Failed(t *testing.T) {
	runMock, env := dockerComposePrepare(t, "service11")

	// Force autoStart
	runMock.OnState("service11", false).Return(nil)
	runMock.OnCall([]string{"compose", "up", "-d"}).Return(nil)
	runMock.OnCall([]string{"compose", "exec", "service11", "test"}).Return(nil)

	err := env.RunExecutable(t.Context(), internal.RunOptions{}, "test", nil)
	require.NoError(t, err)

	runMock.OnCall([]string{"compose", "stop"}).
		Return(errors.New("failed"))

	err = env.Cleanup(t.Context())
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

	executable, err := env.FindExecutable(t.Context(), "tool1")
	require.NotNil(t, executable)
	require.NoError(t, err)
	assert.Equal(t, "/usr/bin/tool1", executable.Path)

	runMock.AssertExpectations(t)
}

func TestDockerCompose_Environment_FindExecutable_NotFound(t *testing.T) {
	runMock, env := dockerComposePrepare(t, "service13")

	// Is running
	runMock.OnState("service13", true).
		Return(nil).
		Once()

	runMock.OnCallOutput([]string{"compose", "exec", "service13", "which", "tool1"}, "/usr/bin/tool1").
		Return(errors.New("failed"))

	executable, err := env.FindExecutable(t.Context(), "tool1")
	assert.Nil(t, executable)
	require.NoError(t, err)

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

	executable, err := env.FindExecutable(t.Context(), "tool1")
	require.NotNil(t, executable)
	require.NoError(t, err)
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

	executable, err := env.FindExecutable(t.Context(), "tool1")
	assert.Nil(t, executable)
	require.EqualError(t, err, "failed")

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

	err := env.RunExecutable(t.Context(), internal.RunOptions{}, "tool1", []string{"arg1", "arg2"})
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

	err := env.RunExecutable(t.Context(), internal.RunOptions{}, "tool1", []string{"arg1", "arg2"})
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

	err := env.RunExecutable(t.Context(), internal.RunOptions{}, "tool1", []string{"arg1", "arg2"})
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

	err := env.RunExecutable(
		t.Context(),
		internal.RunOptions{},
		"/usr/bin/tool1",
		[]string{"arg1", "arg2"},
	)
	require.EqualError(t, err, "docker compose start failed: failed")

	runMock.AssertExpectations(t)
}
