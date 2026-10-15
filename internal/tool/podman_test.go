package tool_test

import (
	"context"
	"errors"
	"fmt"
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

func TestPodman_NewPodman_NoExecutable(t *testing.T) {
	podman := tool.NewPodman(test.LazyExecutableNil)

	assert.NotNil(t, podman)
	assert.Equal(t, tool.IdPodman, podman.Id())
	assert.NotEmpty(t, podman.Name())
	assert.NotEmpty(t, podman.Info())
	assert.Equal(t, []string{"p"}, podman.Aliases())
	assert.False(t, podman.IsAvailable())
}

func TestPodman_NewPodman_WithExecutable(t *testing.T) {
	podman := tool.NewPodman(test.LazyExecutable("/bin/podman"))

	assert.NotNil(t, podman)
	assert.Equal(t, tool.IdPodman, podman.Id())
	assert.True(t, podman.IsAvailable())
}

func TestPodman_DetectPodman_NotFound(t *testing.T) {
	podman := test.RunDetectToolNotFound(t, podmanDetect, []string{"podman"})

	assert.Equal(t, tool.IdPodman, podman.Id())
}

func TestPodman_DetectPodman_Found(t *testing.T) {
	podman := test.RunDetectToolLastFound(t, podmanDetect, []string{"podman"})

	assert.Equal(t, tool.IdPodman, podman.Id())
}

func TestPodman_DetectPodman_Config(t *testing.T) {
	podman := test.RunDetectToolConfig(t, podmanDetect, "podman")

	assert.Equal(t, tool.IdPodman, podman.Id())
}

func TestPodman_Detect(t *testing.T) {
	podman := tool.NewPodman(test.LazyExecutable("podman"))
	assert.NotNil(t, podman)

	env := podman.Detect(".")
	assert.Nil(t, env)
}

func TestPodman_CreateEnvironment(t *testing.T) {
	podman := tool.NewPodman(test.LazyExecutable("podman"))
	assert.NotNil(t, podman)

	env, err := podman.CreateEnvironment(".", "image1")
	require.NoError(t, err)
	assert.NotNil(t, env)

	assert.Equal(t, "podman", env.Id())
}

func TestPodman_CreateEnvironment_NoImage(t *testing.T) {
	podman := tool.NewPodman(test.LazyExecutable("podman"))
	assert.NotNil(t, podman)

	env, err := podman.CreateEnvironment(".", "")
	require.EqualError(t, err, "podman image name is required")
	assert.Nil(t, env)
}

type podmanRunMock struct {
	test.Executable
}

func (m *podmanRunMock) OnCall(args []string) *mock.Call {
	return m.OnRun("podman", args)
}

func (m *podmanRunMock) OnCallOutput(args []string, output string) *mock.Call {
	return m.OnCall(args).
		Run(func(args mock.Arguments) {
			c := args.Get(1).(internal.RunOptions)
			_, _ = c.Output.Write([]byte(output))
		})
}

func (m *podmanRunMock) OnStart(image string, containerId string) *mock.Call {
	absPath, _ := filepath.Abs(".")

	return m.OnCallOutput([]string{
		"run", "--rm", "-d",
		"-v", fmt.Sprintf("%s:/work", absPath),
		"-w", "/work",
		image,
		"/bin/bash", "-c", "trap : TERM INT; sleep infinity & wait",
	}, containerId)
}

func (m *podmanRunMock) OnInspect(containerId string, output string) *mock.Call {
	return m.OnCallOutput([]string{"inspect", "--format", "json", containerId}, output)
}

func (m *podmanRunMock) OnInspectResult(containerId string, result bool) *mock.Call {
	return m.OnInspect(containerId, fmt.Sprintf(`[{"State": {"Running": %v}}]`, result))
}

func podmanPrepare(t *testing.T, image string) (*podmanRunMock, internal.Environment) {
	runMock := podmanRunMock{}
	runMock.Test(t)

	podman := tool.NewPodman(runMock.LazyExecutable("podman"))
	assert.NotNil(t, podman)

	env, err := podman.CreateEnvironment(".", image)
	require.NoError(t, err)
	assert.NotNil(t, env)

	return &runMock, env
}

func TestPodman_Environment_Start(t *testing.T) {
	runMock, env := podmanPrepare(t, "image1")

	runMock.OnStart("image1", "3961248ad455").
		Return(nil)

	err := env.Start(t.Context())
	require.NoError(t, err)

	runMock.AssertExpectations(t)
}

func TestPodman_Environment_Start_Failed(t *testing.T) {
	runMock, env := podmanPrepare(t, "image2")

	runMock.OnStart("image2", "").
		Return(errors.New("failed"))

	err := env.Start(t.Context())
	require.EqualError(t, err, "podman run failed: failed")

	runMock.AssertExpectations(t)
}

func TestPodman_Environment_Start_AlreadyRunning(t *testing.T) {
	runMock, env := podmanPrepare(t, "image3")

	runMock.OnStart("image3", "89054a3c5ff8").
		Return(nil)

	err := env.Start(t.Context())
	require.NoError(t, err)

	runMock.OnInspectResult("89054a3c5ff8", true).
		Return(nil)

	err = env.Start(t.Context())
	require.NoError(t, err)

	runMock.AssertExpectations(t)
}

func TestPodman_Environment_IsRunning_True(t *testing.T) {
	runMock, env := podmanPrepare(t, "image4")

	runMock.OnStart("image4", "3d2219139f14").
		Return(nil)

	err := env.Start(t.Context())
	require.NoError(t, err)

	runMock.OnInspectResult("3d2219139f14", true).
		Return(nil)

	result := env.IsRunning(t.Context())
	assert.True(t, result)

	runMock.AssertExpectations(t)
}

func TestPodman_Environment_IsRunning_False(t *testing.T) {
	runMock, env := podmanPrepare(t, "image5")
	// env.(*tool.PodmanEnvironment).containerId = "e24aa1130c46"

	runMock.OnInspectResult("e24aa1130c46", false).
		Return(nil)

	result := env.IsRunning(t.Context())
	assert.False(t, result)

	runMock.AssertExpectations(t)
}

func TestPodman_Environment_IsRunning_InvalidJson(t *testing.T) {
	runMock, env := podmanPrepare(t, "image6")
	// env.(*tool.PodmanEnvironment).containerId = "8ae4358f08b6"

	runMock.OnInspect("8ae4358f08b6", `{"invalid": "json"`).
		Return(nil)

	result := env.IsRunning(t.Context())
	assert.False(t, result)

	runMock.AssertExpectations(t)
}

func TestPodman_Environment_IsRunning_EmptyJsonArray(t *testing.T) {
	runMock, env := podmanPrepare(t, "image7")
	// env.(*tool.PodmanEnvironment).containerId = "827e5611389a"

	runMock.OnInspect("827e5611389a", `[]`).
		Return(nil)

	result := env.IsRunning(t.Context())
	assert.False(t, result)

	runMock.AssertExpectations(t)
}

func TestPodman_Environment_IsRunning_Failed(t *testing.T) {
	runMock, env := podmanPrepare(t, "image8")
	// env.(*tool.PodmanEnvironment).containerId = "bd49ab6a6a00"

	runMock.OnInspect("bd49ab6a6a00", "").
		Return(errors.New("failed"))

	result := env.IsRunning(t.Context())
	assert.False(t, result)

	runMock.AssertExpectations(t)
}

func TestPodman_Environment_Stop(t *testing.T) {
	runMock, env := podmanPrepare(t, "image9")

	runMock.OnStart("image9", "93cff8e00c49").
		Return(nil)

	err := env.Start(t.Context())
	require.NoError(t, err)

	runMock.OnCall([]string{"stop", "93cff8e00c49"}).
		Return(nil)

	err = env.Stop(t.Context())
	require.NoError(t, err)

	runMock.AssertExpectations(t)
}

func TestPodman_Environment_Stop_Failed(t *testing.T) {
	runMock, env := podmanPrepare(t, "image10")

	runMock.OnStart("image10", "62473370e7ee").
		Return(nil)

	err := env.Start(t.Context())
	require.NoError(t, err)

	runMock.OnCall([]string{"stop", "62473370e7ee"}).
		Return(errors.New("failed"))

	err = env.Stop(t.Context())
	require.EqualError(t, err, "podman run failed: failed")

	runMock.AssertExpectations(t)
}

func TestPodman_Environment_Cleanup(t *testing.T) {
	runMock, env := podmanPrepare(t, "image11")

	err := env.Cleanup(t.Context())
	require.NoError(t, err)

	runMock.AssertExpectations(t)
}

func TestPodman_Environment_Cleanup_Running(t *testing.T) {
	runMock, env := podmanPrepare(t, "image12")

	runMock.OnStart("image12", "c31b2ec8b325").
		Return(nil)

	runMock.OnCall([]string{"exec", "c31b2ec8b325", "which", "echo"}).
		Return(nil)

	executable, err := env.FindExecutable(t.Context(), "echo")
	require.NotNil(t, executable)
	require.NoError(t, err)

	runMock.OnCall([]string{"stop", "c31b2ec8b325"}).
		Return(nil)

	err = env.Cleanup(t.Context())
	require.NoError(t, err)

	runMock.AssertExpectations(t)
}

func TestPodman_Environment_Cleanup_Running_Failed(t *testing.T) {
	runMock, env := podmanPrepare(t, "image13")

	runMock.OnStart("image13", "e01c76bb1351").
		Return(nil)

	runMock.OnCall([]string{"exec", "e01c76bb1351", "which", "echo"}).
		Return(nil)

	executable, err := env.FindExecutable(t.Context(), "echo")
	require.NotNil(t, executable)
	require.NoError(t, err)

	runMock.OnCall([]string{"stop", "e01c76bb1351"}).
		Return(errors.New("failed"))

	err = env.Cleanup(t.Context())
	require.EqualError(t, err, "podman run failed: failed")

	runMock.AssertExpectations(t)
}

func TestPodman_Environment_FindExecutable(t *testing.T) {
	runMock, env := podmanPrepare(t, "image14")

	runMock.OnStart("image14", "15dc587d6d92").
		Return(nil)

	err := env.Start(t.Context())
	require.NoError(t, err)

	// Is running
	runMock.OnInspectResult("15dc587d6d92", true).
		Return(nil).
		Once()

	runMock.OnCallOutput([]string{"exec", "15dc587d6d92", "which", "tool1"}, "/usr/bin/tool1").
		Return(nil)

	executable, err := env.FindExecutable(t.Context(), "tool1")
	require.NotNil(t, executable)
	require.NoError(t, err)
	assert.Equal(t, "/usr/bin/tool1", executable.Path)

	runMock.AssertExpectations(t)
}

func TestPodman_Environment_FindExecutable_WorkingDirectory(t *testing.T) {
	// FIXME: path conversion based on podman mode
	if runtime.GOOS == "windows" { // nolint: goconst
		t.Skip("Problem with path conversion between Windows and podman")
	}

	dir := t.TempDir()
	binDir := filepath.Join(dir, "bin")
	require.NoError(t, os.MkdirAll(binDir, os.ModePerm))
	require.NoError(t, os.WriteFile(filepath.Join(binDir, "echo"), []byte(""), os.ModePerm))

	t.Chdir(dir)

	runMock := podmanRunMock{}
	runMock.Test(t)

	podman := tool.NewPodman(runMock.LazyExecutable("podman"))
	assert.NotNil(t, podman)

	env, err := podman.CreateEnvironment(dir, "image14")
	require.NoError(t, err)
	assert.NotNil(t, env)

	runMock.OnCallOutput([]string{
		"run", "--rm", "-d",
		"-v", fmt.Sprintf("%s:/work", dir),
		"-w", "/work",
		"image14",
		"/bin/bash", "-c", "trap : TERM INT; sleep infinity & wait",
	}, "15dc587d6d92").Return(nil)

	require.NoError(t, env.Start(t.Context()))

	// Is running
	runMock.OnInspectResult("15dc587d6d92", true).
		Return(nil).
		Once()

	runMock.OnCallOutput([]string{"exec", "15dc587d6d92", "which", "bin/echo"}, "/usr/bin/tool1").
		Return(nil)

	executable, err := env.FindExecutable(t.Context(), filepath.Join(dir, "bin", "echo"))
	require.NotNil(t, executable)
	require.NoError(t, err)
	assert.Equal(t, "/usr/bin/tool1", executable.Path)

	runMock.AssertExpectations(t)
}

func TestPodman_Environment_FindExecutable_NotFound(t *testing.T) {
	runMock, env := podmanPrepare(t, "image15")

	runMock.OnStart("image15", "15dc587d6d92").
		Return(nil)

	err := env.Start(t.Context())
	require.NoError(t, err)

	// Is running
	runMock.OnInspectResult("15dc587d6d92", true).
		Return(nil).
		Once()

	runMock.OnCallOutput([]string{"exec", "15dc587d6d92", "which", "tool1"}, "/usr/bin/tool1").
		Return(errors.New("failed"))

	executable, err := env.FindExecutable(t.Context(), "tool1")
	assert.Nil(t, executable)
	require.NoError(t, err)

	runMock.AssertExpectations(t)
}

func TestPodman_Environment_FindExecutable_AutoStart(t *testing.T) {
	runMock, env := podmanPrepare(t, "image16")

	// Start
	runMock.OnStart("image16", "e1ff62268161").
		Return(nil)

	runMock.OnCallOutput([]string{"exec", "e1ff62268161", "which", "tool1"}, "/usr/bin/tool1").
		Return(nil)

	executable, err := env.FindExecutable(t.Context(), "tool1")
	require.NotNil(t, executable)
	require.NoError(t, err)
	assert.Equal(t, "/usr/bin/tool1", executable.Path)

	runMock.AssertExpectations(t)
}

func TestPodman_Environment_FindExecutable_AutoStart_Failed(t *testing.T) {
	runMock, env := podmanPrepare(t, "image17")

	// Start
	runMock.OnStart("image17", "db0ac83ce405").
		Return(errors.New("failed"))

	executable, err := env.FindExecutable(t.Context(), "tool1")
	assert.Nil(t, executable)
	require.EqualError(t, err, "podman run failed: failed")

	runMock.AssertExpectations(t)
}

func TestPodman_Environment_RunExecutable(t *testing.T) {
	runMock, env := podmanPrepare(t, "image18")

	runMock.OnStart("image18", "33910ef9319e").
		Return(nil)

	err := env.Start(t.Context())
	require.NoError(t, err)

	// Is running
	runMock.OnInspectResult("33910ef9319e", true).
		Return(nil).
		Once()

	runMock.OnCall([]string{"exec", "33910ef9319e", "tool1", "arg1", "arg2"}).
		Return(nil).
		Once()

	err = env.RunExecutable(t.Context(), internal.RunOptions{}, "tool1", []string{"arg1", "arg2"})
	require.NoError(t, err)

	runMock.AssertExpectations(t)
}

func TestPodman_Environment_RunExecutable_Failed(t *testing.T) {
	runMock, env := podmanPrepare(t, "image19")

	runMock.OnStart("image19", "db0ac83ce405").
		Return(nil)

	err := env.Start(t.Context())
	require.NoError(t, err)

	// Is running
	runMock.OnInspectResult("db0ac83ce405", true).
		Return(nil).
		Once()

	runMock.OnCall([]string{"exec", "db0ac83ce405", "tool1", "arg1", "arg2"}).
		Return(errors.New("failed")).
		Once()

	err = env.RunExecutable(t.Context(), internal.RunOptions{}, "tool1", []string{"arg1", "arg2"})
	require.EqualError(t, err, "failed")

	runMock.AssertExpectations(t)
}

func TestPodman_Environment_RunExecutable_Env(t *testing.T) {
	runMock, env := podmanPrepare(t, "image18")

	runMock.OnStart("image18", "33910ef9319e").
		Return(nil)

	err := env.Start(t.Context())
	require.NoError(t, err)

	// Is running
	runMock.OnInspectResult("33910ef9319e", true).
		Return(nil).
		Once()

	runMock.OnCall([]string{"exec", "-e", "VAR1=value1", "-e", "VAR2=value2", "33910ef9319e", "tool1", "arg1", "arg2"}).
		Return(nil).
		Once()

	err = env.RunExecutable(t.Context(), internal.RunOptions{
		Env: []string{"VAR1=value1", "VAR2=value2"},
	}, "tool1", []string{"arg1", "arg2"})
	require.NoError(t, err)

	runMock.AssertExpectations(t)
}

func TestPodman_Environment_RunExecutable_AutoStart(t *testing.T) {
	runMock, env := podmanPrepare(t, "image20")

	// Start
	runMock.OnStart("image20", "3961248ad455").
		Return(nil).
		Once()

	runMock.OnCall([]string{"exec", "3961248ad455", "tool1", "arg1", "arg2"}).
		Return(nil).
		Once()

	err := env.RunExecutable(t.Context(), internal.RunOptions{}, "tool1", []string{"arg1", "arg2"})
	require.NoError(t, err)

	runMock.AssertExpectations(t)
}

func TestPodman_Environment_RunExecutable_AutoStart_Failed(t *testing.T) {
	runMock, env := podmanPrepare(t, "image21")

	// Start
	runMock.OnStart("image21", "3961248ad455").
		Return(errors.New("failed")).
		Once()

	err := env.RunExecutable(
		t.Context(),
		internal.RunOptions{},
		"/usr/bin/tool1",
		[]string{"arg1", "arg2"},
	)
	require.EqualError(t, err, "podman run failed: failed")

	runMock.AssertExpectations(t)
}

func podmanDetect(ctx context.Context, options internal.DetectOptions) internal.Tool {
	return tool.DetectPodman(ctx, options)
}
