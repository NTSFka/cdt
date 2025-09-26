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
	"github.com/stretchr/testify/require"
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
	require.NoError(t, err)
	assert.NotNil(t, env)

	assert.Equal(t, "docker", env.Id())
}

func TestDocker_CreateEnvironment_NoImage(t *testing.T) {
	tool := NewDocker(func() *internal.Executable { return &internal.Executable{Path: "docker"} })
	assert.NotNil(t, tool)

	env, err := tool.CreateEnvironment(".", "")
	require.EqualError(t, err, "docker image name is required")
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
	require.NoError(t, err)
	assert.NotNil(t, env)

	return &runMock, env
}

func TestDocker_Environment_Start(t *testing.T) {
	runMock, env := dockerPrepare(t, "image1")

	runMock.OnStart("image1", "3961248ad455").
		Return(nil)

	err := env.Start(context.Background())
	require.NoError(t, err)

	runMock.AssertExpectations(t)
}

func TestDocker_Environment_Start_Failed(t *testing.T) {
	runMock, env := dockerPrepare(t, "image2")

	runMock.OnStart("image2", "").
		Return(errors.New("failed"))

	err := env.Start(context.Background())
	require.EqualError(t, err, "docker run failed: failed")

	runMock.AssertExpectations(t)
}

func TestDocker_Environment_Start_AlreadyRunning(t *testing.T) {
	runMock, env := dockerPrepare(t, "image3")
	env.(*dockerEnvironment).containerId = "89054a3c5ff8"

	runMock.OnInspectResult("89054a3c5ff8", true).
		Return(nil)

	err := env.Start(context.Background())
	require.NoError(t, err)

	runMock.AssertExpectations(t)
}

func TestDocker_Environment_IsRunning_True(t *testing.T) {
	runMock, env := dockerPrepare(t, "image4")
	env.(*dockerEnvironment).containerId = "3d2219139f14"

	runMock.OnInspectResult("3d2219139f14", true).
		Return(nil)

	result := env.IsRunning(context.Background())
	assert.True(t, result)

	runMock.AssertExpectations(t)
}

func TestDocker_Environment_IsRunning_False(t *testing.T) {
	runMock, env := dockerPrepare(t, "image5")
	env.(*dockerEnvironment).containerId = "e24aa1130c46"

	runMock.OnInspectResult("e24aa1130c46", false).
		Return(nil)

	result := env.IsRunning(context.Background())
	assert.False(t, result)

	runMock.AssertExpectations(t)
}

func TestDocker_Environment_IsRunning_InvalidJson(t *testing.T) {
	runMock, env := dockerPrepare(t, "image6")
	env.(*dockerEnvironment).containerId = "8ae4358f08b6"

	runMock.OnInspect("8ae4358f08b6", `{"invalid": "json"`).
		Return(nil)

	result := env.IsRunning(context.Background())
	assert.False(t, result)

	runMock.AssertExpectations(t)
}

func TestDocker_Environment_IsRunning_EmptyJsonArray(t *testing.T) {
	runMock, env := dockerPrepare(t, "image7")
	env.(*dockerEnvironment).containerId = "827e5611389a"

	runMock.OnInspect("827e5611389a", `[]`).
		Return(nil)

	result := env.IsRunning(context.Background())
	assert.False(t, result)

	runMock.AssertExpectations(t)
}

func TestDocker_Environment_IsRunning_Failed(t *testing.T) {
	runMock, env := dockerPrepare(t, "image8")
	env.(*dockerEnvironment).containerId = "bd49ab6a6a00"

	runMock.OnInspect("bd49ab6a6a00", "").
		Return(errors.New("failed"))

	result := env.IsRunning(context.Background())
	assert.False(t, result)

	runMock.AssertExpectations(t)
}

func TestDocker_Environment_Stop(t *testing.T) {
	runMock, env := dockerPrepare(t, "image9")
	env.(*dockerEnvironment).containerId = "93cff8e00c49"

	runMock.OnCall([]string{"stop", "93cff8e00c49"}).
		Return(nil)

	err := env.Stop(context.Background())
	require.NoError(t, err)

	runMock.AssertExpectations(t)
}

func TestDocker_Environment_Stop_Failed(t *testing.T) {
	runMock, env := dockerPrepare(t, "image10")
	env.(*dockerEnvironment).containerId = "62473370e7ee"

	runMock.OnCall([]string{"stop", "62473370e7ee"}).
		Return(errors.New("failed"))

	err := env.Stop(context.Background())
	require.EqualError(t, err, "docker run failed: failed")

	runMock.AssertExpectations(t)
}

func TestDocker_Environment_Cleanup(t *testing.T) {
	runMock, env := dockerPrepare(t, "image11")

	err := env.Cleanup(context.Background())
	require.NoError(t, err)

	runMock.AssertExpectations(t)
}

func TestDocker_Environment_Cleanup_Running(t *testing.T) {
	runMock, env := dockerPrepare(t, "image12")
	env.(*dockerEnvironment).autoStop = true
	env.(*dockerEnvironment).containerId = "c31b2ec8b325"

	runMock.OnCall([]string{"stop", "c31b2ec8b325"}).
		Return(nil)

	err := env.Cleanup(context.Background())
	require.NoError(t, err)

	runMock.AssertExpectations(t)
}

func TestDocker_Environment_Cleanup_Running_Failed(t *testing.T) {
	runMock, env := dockerPrepare(t, "image13")
	env.(*dockerEnvironment).autoStop = true
	env.(*dockerEnvironment).containerId = "e01c76bb1351"

	runMock.OnCall([]string{"stop", "e01c76bb1351"}).
		Return(errors.New("failed"))

	err := env.Cleanup(context.Background())
	require.EqualError(t, err, "docker run failed: failed")

	runMock.AssertExpectations(t)
}

func TestDocker_Environment_FindExecutable(t *testing.T) {
	runMock, env := dockerPrepare(t, "image14")
	env.(*dockerEnvironment).containerId = "15dc587d6d92"

	// Is running
	runMock.OnInspectResult("15dc587d6d92", true).
		Return(nil).
		Once()

	runMock.OnCallOutput([]string{"exec", "15dc587d6d92", "which", "tool1"}, "/usr/bin/tool1").
		Return(nil)

	executable := env.FindExecutable(context.Background(), "tool1")
	require.NotNil(t, executable)
	assert.Equal(t, "/usr/bin/tool1", executable.Path)

	runMock.AssertExpectations(t)
}

func TestDocker_Environment_FindExecutable_Failed(t *testing.T) {
	runMock, env := dockerPrepare(t, "image15")
	env.(*dockerEnvironment).containerId = "15dc587d6d92"

	// Is running
	runMock.OnInspectResult("15dc587d6d92", true).
		Return(nil).
		Once()

	runMock.OnCallOutput([]string{"exec", "15dc587d6d92", "which", "tool1"}, "/usr/bin/tool1").
		Return(errors.New("failed"))

	executable := env.FindExecutable(context.Background(), "tool1")
	assert.Nil(t, executable)

	runMock.AssertExpectations(t)
}

func TestDocker_Environment_FindExecutable_AutoStart(t *testing.T) {
	runMock, env := dockerPrepare(t, "image16")

	// Start
	runMock.OnStart("image16", "e1ff62268161").
		Return(nil)

	runMock.OnCallOutput([]string{"exec", "e1ff62268161", "which", "tool1"}, "/usr/bin/tool1").
		Return(nil)

	executable := env.FindExecutable(context.Background(), "tool1")
	require.NotNil(t, executable)
	assert.Equal(t, "/usr/bin/tool1", executable.Path)

	runMock.AssertExpectations(t)
}

func TestDocker_Environment_FindExecutable_AutoStart_Failed(t *testing.T) {
	runMock, env := dockerPrepare(t, "image17")

	// Start
	runMock.OnStart("image17", "db0ac83ce405").
		Return(errors.New("failed"))

	executable := env.FindExecutable(context.Background(), "tool1")
	assert.Nil(t, executable)

	runMock.AssertExpectations(t)
}

func TestDocker_Environment_RunExecutable(t *testing.T) {
	runMock, env := dockerPrepare(t, "image18")
	env.(*dockerEnvironment).containerId = "33910ef9319e"

	// Is running
	runMock.OnInspectResult("33910ef9319e", true).
		Return(nil).
		Once()

	runMock.OnCall([]string{"exec", "33910ef9319e", "tool1", "arg1", "arg2"}).
		Return(nil).
		Once()

	err := env.RunExecutable(context.Background(), internal.RunOptions{}, "tool1", []string{"arg1", "arg2"})
	require.NoError(t, err)

	runMock.AssertExpectations(t)
}

func TestDocker_Environment_RunExecutable_Failed(t *testing.T) {
	runMock, env := dockerPrepare(t, "image19")
	env.(*dockerEnvironment).containerId = "db0ac83ce405"

	// Is running
	runMock.OnInspectResult("db0ac83ce405", true).
		Return(nil).
		Once()

	runMock.OnCall([]string{"exec", "db0ac83ce405", "tool1", "arg1", "arg2"}).
		Return(errors.New("failed")).
		Once()

	err := env.RunExecutable(context.Background(), internal.RunOptions{}, "tool1", []string{"arg1", "arg2"})
	require.EqualError(t, err, "failed")

	runMock.AssertExpectations(t)
}

func TestDocker_Environment_RunExecutable_AutoStart(t *testing.T) {
	runMock, env := dockerPrepare(t, "image20")

	// Start
	runMock.OnStart("image20", "3961248ad455").
		Return(nil).
		Once()

	runMock.OnCall([]string{"exec", "3961248ad455", "tool1", "arg1", "arg2"}).
		Return(nil).
		Once()

	err := env.RunExecutable(context.Background(), internal.RunOptions{}, "tool1", []string{"arg1", "arg2"})
	require.NoError(t, err)

	runMock.AssertExpectations(t)
}

func TestDocker_Environment_RunExecutable_AutoStart_Failed(t *testing.T) {
	runMock, env := dockerPrepare(t, "image21")

	// Start
	runMock.OnStart("image21", "3961248ad455").
		Return(errors.New("failed")).
		Once()

	err := env.RunExecutable(context.Background(), internal.RunOptions{}, "/usr/bin/tool1", []string{"arg1", "arg2"})
	require.EqualError(t, err, "docker run failed: failed")

	runMock.AssertExpectations(t)
}
