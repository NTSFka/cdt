package tool_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"cdt/internal"
	"cdt/internal/test"
	"cdt/internal/tool"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestPodmanCompose_NewPodmanCompose_NoExecutable(t *testing.T) {
	podmanCompose := tool.NewPodmanCompose(test.LazyExecutableNil)

	assert.NotNil(t, podmanCompose)
	assert.Equal(t, tool.IdPodmanCompose, podmanCompose.Id())
	assert.NotEmpty(t, podmanCompose.Name())
	assert.NotEmpty(t, podmanCompose.Info())
	assert.Equal(t, []string{"dc"}, podmanCompose.Aliases())
	assert.False(t, podmanCompose.IsAvailable())
}

func TestPodmanCompose_NewPodmanCompose_WithExecutable(t *testing.T) {
	podmanCompose := tool.NewPodmanCompose(test.LazyExecutable("/bin/podman"))

	assert.NotNil(t, podmanCompose)
	assert.Equal(t, tool.IdPodmanCompose, podmanCompose.Id())
	assert.True(t, podmanCompose.IsAvailable())
}

func TestPodmanCompose_DetectPodmanCompose_NotFound(t *testing.T) {
	podmanCompose := test.RunDetectToolNotFound(t, podmanComposeDetect, []string{"podman"})

	assert.Equal(t, tool.IdPodmanCompose, podmanCompose.Id())
}

func TestPodmanCompose_DetectPodmanCompose_Found(t *testing.T) {
	podmanCompose := test.RunDetectToolLastFound(t, podmanComposeDetect, []string{"podman"})

	assert.Equal(t, tool.IdPodmanCompose, podmanCompose.Id())
}

func TestPodmanCompose_DetectPodmanCompose_Config(t *testing.T) {
	podmanCompose := test.RunDetectToolConfig(t, podmanComposeDetect, "podman-compose")

	assert.Equal(t, tool.IdPodmanCompose, podmanCompose.Id())
}

func TestPodmanCompose_Detect(t *testing.T) {
	podmanCompose := tool.NewPodmanCompose(test.LazyExecutable("podman"))
	assert.NotNil(t, podmanCompose)

	env := podmanCompose.Detect(".")
	assert.Nil(t, env)
}

func TestPodmanCompose_CreateEnvironment(t *testing.T) {
	podmanCompose := tool.NewPodmanCompose(test.LazyExecutable("podman"))
	assert.NotNil(t, podmanCompose)

	env, err := podmanCompose.CreateEnvironment(".", "service1")
	require.NoError(t, err)
	assert.NotNil(t, env)

	assert.Equal(t, "podman-compose", env.Id())
}

func TestPodmanCompose_CreateEnvironment_NoService(t *testing.T) {
	podmanCompose := tool.NewPodmanCompose(test.LazyExecutable("podman"))
	assert.NotNil(t, podmanCompose)

	env, err := podmanCompose.CreateEnvironment(".", "")
	require.EqualError(t, err, "service name is required")
	assert.Nil(t, env)
}

type podmanComposeRunMock struct {
	test.Executable
}

func (m *podmanComposeRunMock) OnCall(args []string) *mock.Call {
	return m.OnRun("podman", args)
}

func (m *podmanComposeRunMock) OnCallOutput(args []string, output string) *mock.Call {
	return m.OnCall(args).
		Run(func(args mock.Arguments) {
			c := args.Get(1).(internal.RunOptions)
			_, _ = c.Output.Write([]byte(output))
		})
}

func (m *podmanComposeRunMock) OnState(service string, result bool) *mock.Call {
	var output string

	if result {
		output = `{"State": "running"}`
	} else {
		output = `{"State": "exited"}`
	}

	return m.OnCallOutput([]string{"compose", "ps", "--format", "json", service}, output)
}

func podmanComposePrepare(
	t *testing.T,
	service string,
) (*podmanComposeRunMock, internal.Environment) {
	runMock := podmanComposeRunMock{}

	podmanCompose := tool.NewPodmanCompose(runMock.LazyExecutable("podman"))
	assert.NotNil(t, podmanCompose)

	env, err := podmanCompose.CreateEnvironment(".", service)
	require.NoError(t, err)
	assert.NotNil(t, env)

	return &runMock, env
}

func TestPodmanCompose_Environment_Start(t *testing.T) {
	runMock, env := podmanComposePrepare(t, "service1")

	runMock.OnCall([]string{"compose", "up", "-d"}).
		Return(nil)

	err := env.Start(t.Context())
	require.NoError(t, err)

	runMock.AssertExpectations(t)
}

func TestPodmanCompose_Environment_Start_Failed(t *testing.T) {
	runMock, env := podmanComposePrepare(t, "service2")

	runMock.OnCall([]string{"compose", "up", "-d"}).
		Return(errors.New("failed"))

	err := env.Start(t.Context())
	require.EqualError(t, err, "failed")

	runMock.AssertExpectations(t)
}

func TestPodmanCompose_Environment_IsRunning_True(t *testing.T) {
	runMock, env := podmanComposePrepare(t, "service3")

	runMock.OnState("service3", true).
		Return(nil)

	result := env.IsRunning(t.Context())
	assert.True(t, result)

	runMock.AssertExpectations(t)
}

func TestPodmanCompose_Environment_IsRunning_False(t *testing.T) {
	runMock, env := podmanComposePrepare(t, "service4")

	runMock.OnState("service4", false).
		Return(nil)

	result := env.IsRunning(t.Context())
	assert.False(t, result)

	runMock.AssertExpectations(t)
}

func TestPodmanCompose_Environment_IsRunning_InvalidJson(t *testing.T) {
	runMock, env := podmanComposePrepare(t, "service5")

	runMock.OnCallOutput([]string{"compose", "ps", "--format", "json", "service5"}, `{"invalid": "json"`).
		Return(nil)

	result := env.IsRunning(t.Context())
	assert.False(t, result)

	runMock.AssertExpectations(t)
}

func TestPodmanCompose_Environment_IsRunning_Failed(t *testing.T) {
	runMock, env := podmanComposePrepare(t, "service6")

	runMock.OnCall([]string{"compose", "ps", "--format", "json", "service6"}).
		Return(errors.New("failed"))

	result := env.IsRunning(t.Context())
	assert.False(t, result)

	runMock.AssertExpectations(t)
}

func TestPodmanCompose_Environment_Stop(t *testing.T) {
	runMock, env := podmanComposePrepare(t, "service7")

	runMock.OnCall([]string{"compose", "stop"}).
		Return(nil)

	err := env.Stop(t.Context())
	require.NoError(t, err)

	runMock.AssertExpectations(t)
}

func TestPodmanCompose_Environment_Stop_Failed(t *testing.T) {
	runMock, env := podmanComposePrepare(t, "service8")

	runMock.OnCall([]string{"compose", "stop"}).
		Return(errors.New("failed"))

	err := env.Stop(t.Context())
	require.EqualError(t, err, "failed")

	runMock.AssertExpectations(t)
}

func TestPodmanCompose_Environment_Cleanup(t *testing.T) {
	runMock, env := podmanComposePrepare(t, "service9")

	err := env.Cleanup(t.Context())
	require.NoError(t, err)

	runMock.AssertExpectations(t)
}

func TestPodmanCompose_Environment_Cleanup_Running(t *testing.T) {
	runMock, env := podmanComposePrepare(t, "service10")

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

func TestPodmanCompose_Environment_Cleanup_Running_Failed(t *testing.T) {
	runMock, env := podmanComposePrepare(t, "service11")

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

func TestPodmanCompose_Environment_FindExecutable(t *testing.T) {
	runMock, env := podmanComposePrepare(t, "service12")

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

func TestPodmanCompose_Environment_FindExecutable_WorkingDirectory(t *testing.T) {
	// FIXME: path conversion based on podman mode
	if runtime.GOOS == "windows" { // nolint: goconst
		t.Skip("Problem with path conversion between Windows and podman")
	}

	dir := t.TempDir()
	binDir := filepath.Join(dir, "bin")
	require.NoError(t, os.MkdirAll(binDir, os.ModePerm))
	require.NoError(t, os.WriteFile(filepath.Join(binDir, "echo"), []byte(""), os.ModePerm))

	t.Chdir(dir)

	runMock := podmanComposeRunMock{}

	podmanCompose := tool.NewPodmanCompose(runMock.LazyExecutable("podman"))
	assert.NotNil(t, podmanCompose)

	env, err := podmanCompose.CreateEnvironment(dir, "service12")
	require.NoError(t, err)
	assert.NotNil(t, env)

	// Is running
	runMock.OnState("service12", true).
		Return(nil).
		Once()

	runMock.OnCallOutput([]string{"compose", "exec", "service12", "which", "bin/echo"}, "/usr/bin/tool1").
		Return(nil)

	executable, err := env.FindExecutable(t.Context(), filepath.Join(dir, "bin", "echo"))
	require.NotNil(t, executable)
	require.NoError(t, err)
	assert.Equal(t, "/usr/bin/tool1", executable.Path)

	runMock.AssertExpectations(t)
}

func TestPodmanCompose_Environment_FindExecutable_NotFound(t *testing.T) {
	runMock, env := podmanComposePrepare(t, "service13")

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

func TestPodmanCompose_Environment_FindExecutable_AutoStart(t *testing.T) {
	runMock, env := podmanComposePrepare(t, "service14")

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

func TestPodmanCompose_Environment_FindExecutable_AutoStart_Failed(t *testing.T) {
	runMock, env := podmanComposePrepare(t, "service15")

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

func TestPodmanCompose_Environment_RunExecutable(t *testing.T) {
	runMock, env := podmanComposePrepare(t, "service16")

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

func TestPodmanCompose_Environment_RunExecutable_Failed(t *testing.T) {
	runMock, env := podmanComposePrepare(t, "service17")

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

func TestPodmanCompose_Environment_RunExecutable_Env(t *testing.T) {
	runMock, env := podmanComposePrepare(t, "service16")

	// Is running
	runMock.OnState("service16", true).
		Return(nil).
		Once()

	runMock.OnCall([]string{
		"compose",
		"exec",
		"-e",
		"VAR1=value1",
		"-e",
		"VAR2=value2",
		"service16",
		"tool1",
		"arg1",
		"arg2",
	}).
		Return(nil).
		Once()

	err := env.RunExecutable(t.Context(), internal.RunOptions{
		Env: []string{"VAR1=value1", "VAR2=value2"},
	}, "tool1", []string{"arg1", "arg2"})
	require.NoError(t, err)

	runMock.AssertExpectations(t)
}

func TestPodmanCompose_Environment_RunExecutable_AutoStart(t *testing.T) {
	runMock, env := podmanComposePrepare(t, "service18")

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

func TestPodmanCompose_Environment_RunExecutable_AutoStart_Failed(t *testing.T) {
	runMock, env := podmanComposePrepare(t, "service19")

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
	require.EqualError(t, err, "podman compose start failed: failed")

	runMock.AssertExpectations(t)
}

func podmanComposeDetect(ctx context.Context, options internal.DetectOptions) internal.Tool {
	return tool.DetectPodmanCompose(ctx, options)
}
