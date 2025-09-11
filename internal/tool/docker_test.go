package tool

import (
	"cdt/internal"
	"cdt/internal/test"
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDocker_NewDocker_NoExecutable(t *testing.T) {
	tool := NewDocker(func() *internal.Executable { return nil })

	assert.NotNil(t, tool)
	assert.Equal(t, "docker", tool.Id())
	assert.NotEmpty(t, tool.Name())
	assert.NotEmpty(t, tool.Info())
	assert.Equal(t, []string{"d"}, tool.Aliases())
	assert.False(t, tool.IsAvailable())
}

func TestDocker_NewDocker_WithExecutable(t *testing.T) {
	tool := NewDocker(func() *internal.Executable { return &internal.Executable{Path: "/bin/docker"} })

	assert.NotNil(t, tool)
	assert.Equal(t, "docker", tool.Id())
	assert.True(t, tool.IsAvailable())
}

func TestDocker_DetectDocker_NotFound(t *testing.T) {
	env := test.NewEnvironment(t)
	env.OnFindExecutable("docker").
		Return(nil)

	tool := DetectDocker(context.Background(), env)
	assert.NotNil(t, tool)
	assert.Equal(t, "docker", tool.Id())
	assert.False(t, tool.IsAvailable())

	env.AssertExpectations(t)
}

func TestDocker_DetectDocker_Found(t *testing.T) {
	env := test.NewEnvironment(t)
	env.OnFindExecutable("docker").
		Return(env.NewExecutable("docker"))

	tool := DetectDocker(context.Background(), env)
	assert.NotNil(t, tool)
	assert.Equal(t, "docker", tool.Id())
	assert.True(t, tool.IsAvailable())

	env.AssertExpectations(t)
}

func TestDocker_Detect(t *testing.T) {
	tool := NewDocker(func() *internal.Executable { return &internal.Executable{Path: "docker"} })
	assert.NotNil(t, tool)

	env := tool.Detect(".")
	assert.Nil(t, env)
}

func TestDocker_CreateEnvironment(t *testing.T) {
	tool := NewDocker(func() *internal.Executable { return &internal.Executable{Path: "docker"} })
	assert.NotNil(t, tool)

	env, err := tool.CreateEnvironment(".", "image1")
	assert.NoError(t, err)
	assert.NotNil(t, env)

	assert.Equal(t, "docker", env.Id())
}

func TestDocker_CreateEnvironment_NoImage(t *testing.T) {
	tool := NewDocker(func() *internal.Executable { return &internal.Executable{Path: "docker"} })
	assert.NotNil(t, tool)

	env, err := tool.CreateEnvironment(".", "")
	assert.EqualError(t, err, "docker image name is required")
	assert.Nil(t, env)
}

type dockerRunMock struct {
	test.Executable
}

func (m *dockerRunMock) OnCall(args []string) *mock.Call {
	return m.OnRun("docker", args)
}

func (m *dockerRunMock) OnCallOutput(args []string, output string) *mock.Call {
	return m.OnCall(args).
		Run(func(args mock.Arguments) {
			c := args.Get(1).(internal.RunOptions)
			_, _ = c.Output.Write([]byte(output))
		})
}

func (m *dockerRunMock) OnStart(image string, containerId string) *mock.Call {
	absPath, _ := filepath.Abs(".")

	return m.OnCallOutput([]string{
		"run", "--rm", "-d",
		"-v", fmt.Sprintf("%s:/work", absPath),
		"-w", "/work",
		image,
		"/bin/bash", "-c", "trap : TERM INT; sleep infinity & wait",
	}, containerId)
}

func (m *dockerRunMock) OnInspect(containerId string, output string) *mock.Call {
	return m.OnCallOutput([]string{"inspect", "--format", "json", containerId}, output)
}

func (m *dockerRunMock) OnInspectResult(containerId string, result bool) *mock.Call {
	return m.OnInspect(containerId, fmt.Sprintf(`[{"State": {"Running": %v}}]`, result))
}

func dockerPrepare(t *testing.T, image string) (*dockerRunMock, internal.Environment) {
	runMock := dockerRunMock{}

	tool := NewDocker(runMock.LazyExecutable("docker"))
	assert.NotNil(t, tool)

	env, err := tool.CreateEnvironment(".", image)
	assert.NoError(t, err)
	assert.NotNil(t, env)

	return &runMock, env
}

func TestDocker_Environment_Start(t *testing.T) {
	runMock, env := dockerPrepare(t, "image1")

	runMock.OnStart("image1", "3961248ad455").
		Return(nil)

	err := env.Start(context.Background())
	assert.NoError(t, err)

	runMock.AssertExpectations(t)
}

func TestDocker_Environment_Start_Failed(t *testing.T) {
	runMock, env := dockerPrepare(t, "image1")

	runMock.OnStart("image1", "").
		Return(errors.New("failed"))

	err := env.Start(context.Background())
	assert.EqualError(t, err, "docker run failed: failed")

	runMock.AssertExpectations(t)
}

func TestDocker_Environment_Start_AlreadyRunning(t *testing.T) {
	runMock, env := dockerPrepare(t, "image1")
	env.(*dockerEnvironment).containerId = "db0ac83ce405"

	runMock.OnInspectResult("db0ac83ce405", true).
		Return(nil)

	err := env.Start(context.Background())
	assert.NoError(t, err)

	runMock.AssertExpectations(t)
}

func TestDocker_Environment_IsRunning_True(t *testing.T) {
	runMock, env := dockerPrepare(t, "image1")
	env.(*dockerEnvironment).containerId = "db0ac83ce405"

	runMock.OnInspectResult("db0ac83ce405", true).
		Return(nil)

	result := env.IsRunning(context.Background())
	assert.True(t, result)

	runMock.AssertExpectations(t)
}

func TestDocker_Environment_IsRunning_False(t *testing.T) {
	runMock, env := dockerPrepare(t, "image1")
	env.(*dockerEnvironment).containerId = "db0ac83ce405"

	runMock.OnInspectResult("db0ac83ce405", false).
		Return(nil)

	result := env.IsRunning(context.Background())
	assert.False(t, result)

	runMock.AssertExpectations(t)
}

func TestDocker_Environment_IsRunning_InvalidJson(t *testing.T) {
	runMock, env := dockerPrepare(t, "image1")
	env.(*dockerEnvironment).containerId = "db0ac83ce405"

	runMock.OnInspect("db0ac83ce405", `{"invalid": "json"`).
		Return(nil)

	result := env.IsRunning(context.Background())
	assert.False(t, result)

	runMock.AssertExpectations(t)
}

func TestDocker_Environment_IsRunning_EmptyJsonArray(t *testing.T) {
	runMock, env := dockerPrepare(t, "image1")
	env.(*dockerEnvironment).containerId = "db0ac83ce405"

	runMock.OnInspect("db0ac83ce405", `[]`).
		Return(nil)

	result := env.IsRunning(context.Background())
	assert.False(t, result)

	runMock.AssertExpectations(t)
}

func TestDocker_Environment_IsRunning_Failed(t *testing.T) {
	runMock, env := dockerPrepare(t, "image1")
	env.(*dockerEnvironment).containerId = "db0ac83ce405"

	runMock.OnInspect("db0ac83ce405", "").
		Return(errors.New("failed"))

	result := env.IsRunning(context.Background())
	assert.False(t, result)

	runMock.AssertExpectations(t)
}

func TestDocker_Environment_Stop(t *testing.T) {
	runMock, env := dockerPrepare(t, "image1")
	env.(*dockerEnvironment).containerId = "db0ac83ce405"

	runMock.OnCall([]string{"stop", "db0ac83ce405"}).
		Return(nil)

	err := env.Stop(context.Background())
	assert.NoError(t, err)

	runMock.AssertExpectations(t)
}

func TestDocker_Environment_Stop_Failed(t *testing.T) {
	runMock, env := dockerPrepare(t, "image1")
	env.(*dockerEnvironment).containerId = "db0ac83ce405"

	runMock.OnCall([]string{"stop", "db0ac83ce405"}).
		Return(errors.New("failed"))

	err := env.Stop(context.Background())
	assert.EqualError(t, err, "docker run failed: failed")

	runMock.AssertExpectations(t)
}

func TestDocker_Environment_Cleanup(t *testing.T) {
	runMock, env := dockerPrepare(t, "image1")

	err := env.Cleanup(context.Background())
	assert.NoError(t, err)

	runMock.AssertExpectations(t)
}

func TestDocker_Environment_Cleanup_Running(t *testing.T) {
	runMock, env := dockerPrepare(t, "image1")
	env.(*dockerEnvironment).autoStop = true
	env.(*dockerEnvironment).containerId = "db0ac83ce405"

	runMock.OnCall([]string{"stop", "db0ac83ce405"}).
		Return(nil)

	err := env.Cleanup(context.Background())
	assert.NoError(t, err)

	runMock.AssertExpectations(t)
}

func TestDocker_Environment_Cleanup_Running_Failed(t *testing.T) {
	runMock, env := dockerPrepare(t, "image1")
	env.(*dockerEnvironment).autoStop = true
	env.(*dockerEnvironment).containerId = "db0ac83ce405"

	runMock.OnCall([]string{"stop", "db0ac83ce405"}).
		Return(errors.New("failed"))

	err := env.Cleanup(context.Background())
	assert.EqualError(t, err, "docker run failed: failed")

	runMock.AssertExpectations(t)
}

func TestDocker_Environment_FindExecutable(t *testing.T) {
	runMock, env := dockerPrepare(t, "image1")
	env.(*dockerEnvironment).containerId = "db0ac83ce405"

	// Is running
	runMock.OnInspectResult("db0ac83ce405", true).
		Return(nil).
		Once()

	runMock.OnCallOutput([]string{"exec", "db0ac83ce405", "which", "tool1"}, "/usr/bin/tool1").
		Return(nil)

	executable := env.FindExecutable(context.Background(), "tool1")
	if assert.NotNil(t, executable) {
		assert.Equal(t, "/usr/bin/tool1", executable.Path)
	}

	runMock.AssertExpectations(t)
}

func TestDocker_Environment_FindExecutable_Failed(t *testing.T) {
	runMock, env := dockerPrepare(t, "image1")
	env.(*dockerEnvironment).containerId = "db0ac83ce405"

	// Is running
	runMock.OnInspectResult("db0ac83ce405", true).
		Return(nil).
		Once()

	runMock.OnCallOutput([]string{"exec", "db0ac83ce405", "which", "tool1"}, "/usr/bin/tool1").
		Return(errors.New("failed"))

	executable := env.FindExecutable(context.Background(), "tool1")
	assert.Nil(t, executable)

	runMock.AssertExpectations(t)
}

func TestDocker_Environment_FindExecutable_AutoStart(t *testing.T) {
	runMock, env := dockerPrepare(t, "image1")

	// Start
	runMock.OnStart("image1", "db0ac83ce405").
		Return(nil)

	runMock.OnCallOutput([]string{"exec", "db0ac83ce405", "which", "tool1"}, "/usr/bin/tool1").
		Return(nil)

	executable := env.FindExecutable(context.Background(), "tool1")
	if assert.NotNil(t, executable) {
		assert.Equal(t, "/usr/bin/tool1", executable.Path)
	}

	runMock.AssertExpectations(t)
}

func TestDocker_Environment_FindExecutable_AutoStart_Failed(t *testing.T) {
	runMock, env := dockerPrepare(t, "image1")

	// Start
	runMock.OnStart("image1", "db0ac83ce405").
		Return(errors.New("failed"))

	executable := env.FindExecutable(context.Background(), "tool1")
	assert.Nil(t, executable)

	runMock.AssertExpectations(t)
}

func TestDocker_Environment_RunExecutable(t *testing.T) {
	runMock, env := dockerPrepare(t, "image1")
	env.(*dockerEnvironment).containerId = "db0ac83ce405"

	// Is running
	runMock.OnInspectResult("db0ac83ce405", true).
		Return(nil).
		Once()

	runMock.OnCall([]string{"exec", "db0ac83ce405", "tool1", "arg1", "arg2"}).
		Return(nil).
		Once()

	err := env.RunExecutable(context.Background(), internal.RunOptions{}, "tool1", []string{"arg1", "arg2"})
	assert.NoError(t, err)

	runMock.AssertExpectations(t)
}

func TestDocker_Environment_RunExecutable_Failed(t *testing.T) {
	runMock, env := dockerPrepare(t, "image1")
	env.(*dockerEnvironment).containerId = "db0ac83ce405"

	// Is running
	runMock.OnInspectResult("db0ac83ce405", true).
		Return(nil).
		Once()

	runMock.OnCall([]string{"exec", "db0ac83ce405", "tool1", "arg1", "arg2"}).
		Return(errors.New("failed")).
		Once()

	err := env.RunExecutable(context.Background(), internal.RunOptions{}, "tool1", []string{"arg1", "arg2"})
	assert.EqualError(t, err, "failed")

	runMock.AssertExpectations(t)
}

func TestDocker_Environment_RunExecutable_AutoStart(t *testing.T) {
	runMock, env := dockerPrepare(t, "image1")

	// Start
	runMock.OnStart("image1", "3961248ad455").
		Return(nil).
		Once()

	runMock.OnCall([]string{"exec", "3961248ad455", "tool1", "arg1", "arg2"}).
		Return(nil).
		Once()

	err := env.RunExecutable(context.Background(), internal.RunOptions{}, "tool1", []string{"arg1", "arg2"})
	assert.NoError(t, err)

	runMock.AssertExpectations(t)
}

func TestDocker_Environment_RunExecutable_AutoStart_Failed(t *testing.T) {
	runMock, env := dockerPrepare(t, "image1")

	// Start
	runMock.OnStart("image1", "3961248ad455").
		Return(errors.New("failed")).
		Once()

	err := env.RunExecutable(context.Background(), internal.RunOptions{}, "/usr/bin/tool1", []string{"arg1", "arg2"})
	assert.EqualError(t, err, "docker run failed: failed")

	runMock.AssertExpectations(t)
}
