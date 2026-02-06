package tool_test

import (
	"testing"

	"cdt/internal"
	"cdt/internal/test"
	"cdt/internal/tool"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBlack_DetectBlack(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("black").
		Return(env.NewExecutable("/bin/black"), nil)

	black := tool.DetectBlack(t.Context(), tool.DetectOptions{Environment: env})
	assert.NotNil(t, black)
	assert.Equal(t, "black", black.Id())
	assert.True(t, black.IsAvailable())

	if executable := black.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/black", executable.Path)
	}

	env.AssertExpectations(t)
}

func TestBlack_DetectBlack_NotFound(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("black").
		Return(nil, nil)

	black := tool.DetectBlack(t.Context(), tool.DetectOptions{Environment: env})
	assert.NotNil(t, black)
	assert.Equal(t, "black", black.Id())
	assert.False(t, black.IsAvailable())
	assert.Nil(t, black.Executable())

	env.AssertExpectations(t)
}

func TestBlack_DetectBandit_Config(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("black-1").
		Return(env.NewExecutable("/bin/black"), nil)

	black := tool.DetectBlack(t.Context(), tool.DetectOptions{
		Environment: env,
		ToolsPaths:  map[string]string{"black": "black-1"},
	})
	assert.NotNil(t, black)
	assert.Equal(t, "black", black.Id())
	assert.True(t, black.IsAvailable())

	if executable := black.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/black", executable.Path)
	}

	env.AssertExpectations(t)
}

func TestBlack_Black_FormatFiles_All(t *testing.T) {
	exec := test.NewExecutable(t)

	black := tool.NewBlack(exec.LazyExecutable("format"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("format", []string{"."}).
		Return(nil)

	err := black.FormatFiles(t.Context(), internal.ProjectFormatterOptions{ProjectInfo: info})
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestBlack_Black_FormatFiles(t *testing.T) {
	exec := test.NewExecutable(t)

	black := tool.NewBlack(exec.LazyExecutable("format"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("format", []string{"tests/*"}).
		Return(nil)

	err := black.FormatFiles(
		t.Context(),
		internal.ProjectFormatterOptions{ProjectInfo: info, Filenames: &[]string{"tests/*"}},
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestBlack_Black_FormatFiles_CheckAll(t *testing.T) {
	exec := test.NewExecutable(t)

	black := tool.NewBlack(exec.LazyExecutable("format"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("format", []string{"--check", "."}).
		Return(nil)

	err := black.FormatFiles(
		t.Context(),
		internal.ProjectFormatterOptions{ProjectInfo: info, CheckOnly: true},
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestBlack_Black_FormatFiles_Check(t *testing.T) {
	exec := test.NewExecutable(t)

	black := tool.NewBlack(exec.LazyExecutable("format"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("format", []string{"--check", "tests/*", "/path/to/file.py"}).
		Return(nil)

	err := black.FormatFiles(
		t.Context(),
		internal.ProjectFormatterOptions{
			ProjectInfo: info,
			CheckOnly:   true,
			Filenames:   &[]string{"tests/*", "/path/to/file.py"},
		},
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}
