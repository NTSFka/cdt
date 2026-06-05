package tool_test

import (
	"testing"

	"cdt/internal"
	"cdt/internal/test"
	"cdt/internal/tool"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRuff_DetectRuff(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("ruff").
		Return(env.NewExecutable("/bin/ruff"), nil)

	ruff := tool.DetectRuff(t.Context(), tool.DetectOptions{Environment: env})
	assert.NotNil(t, ruff)
	assert.Equal(t, "ruff", ruff.Id())
	assert.True(t, ruff.IsAvailable())

	if executable := ruff.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/ruff", executable.Path)
	}

	env.AssertExpectations(t)
}

func TestRuff_DetectRuff_NotFound(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("ruff").
		Return(nil, nil)

	ruff := tool.DetectRuff(t.Context(), tool.DetectOptions{Environment: env})
	assert.NotNil(t, ruff)
	assert.Equal(t, "ruff", ruff.Id())
	assert.False(t, ruff.IsAvailable())
	assert.Nil(t, ruff.Executable())

	env.AssertExpectations(t)
}

func TestRuff_DetectRuff_Config(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("ruff2").
		Return(env.NewExecutable("/bin/ruff"), nil)

	ruff := tool.DetectRuff(t.Context(), tool.DetectOptions{
		Environment: env,
		ToolsPaths:  map[string]string{"ruff": "ruff2"},
	})
	assert.NotNil(t, ruff)
	assert.Equal(t, "ruff", ruff.Id())
	assert.True(t, ruff.IsAvailable())

	if executable := ruff.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/ruff", executable.Path)
	}

	env.AssertExpectations(t)
}

func TestRuff_Ruff_LintFiles_All(t *testing.T) {
	exec := test.NewExecutable(t)

	ruff := tool.NewRuff(exec.LazyExecutable("lint"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("lint", []string{"check"}).
		Return(nil)

	err := ruff.LintFiles(t.Context(), internal.ProjectLinterOptions{ProjectInfo: info})
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestRuff_Ruff_LintFiles(t *testing.T) {
	exec := test.NewExecutable(t)

	ruff := tool.NewRuff(exec.LazyExecutable("lint"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("lint", []string{"check", "file.py", "/path/to/file2.py"}).
		Return(nil)

	err := ruff.LintFiles(
		t.Context(),
		internal.ProjectLinterOptions{
			ProjectInfo: info,
			Filenames:   &[]string{"file.py", "/path/to/file2.py"},
		},
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestRuff_Ruff_LintFiles_OutputFormat_Raw(t *testing.T) {
	exec := test.NewExecutable(t)

	linter := tool.NewRuff(exec.LazyExecutable("lint"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("lint", []string{"check"}).
		Return(nil)

	err := linter.LintFiles(
		t.Context(),
		internal.ProjectLinterOptions{
			ProjectInfo: info,
			Output: internal.OutputOptions[internal.LintReportFormat]{
				Format: internal.LintReportFormatRaw,
			},
		},
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestRuff_Ruff_LintFiles_OutputFormat_Unsupported(t *testing.T) {
	dataSet := []struct {
		Format internal.LintReportFormat
	}{
		{internal.LintReportFormatJson},
		{internal.LintReportFormatJUnit},
		{internal.LintReportFormatGitHub},
		{internal.LintReportFormatGitLab},
		{internal.LintReportFormatTeamCity},
		{"test-unsupported"},
	}

	for _, data := range dataSet {
		t.Run(string(data.Format), func(t *testing.T) {
			exec := test.NewExecutable(t)

			linter := tool.NewRuff(exec.LazyExecutable("lint"))

			info := internal.ProjectInfo{Directory: "."}

			err := linter.LintFiles(
				t.Context(),
				internal.ProjectLinterOptions{
					ProjectInfo: info,
					Output: internal.OutputOptions[internal.LintReportFormat]{
						Format: data.Format,
					},
				},
			)
			require.EqualError(t, err, "unsupported report format: "+string(data.Format))

			exec.AssertExpectations(t)
		})
	}
}
func TestRuff_Ruff_FormatFiles_All(t *testing.T) {
	exec := test.NewExecutable(t)

	ruff := tool.NewRuff(exec.LazyExecutable("format"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("format", []string{"format"}).
		Return(nil)

	err := ruff.FormatFiles(t.Context(), internal.ProjectFormatterOptions{ProjectInfo: info})
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestRuff_Ruff_FormatFiles(t *testing.T) {
	exec := test.NewExecutable(t)

	ruff := tool.NewRuff(exec.LazyExecutable("format"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("format", []string{"format", "tests/*"}).
		Return(nil)

	err := ruff.FormatFiles(
		t.Context(),
		internal.ProjectFormatterOptions{ProjectInfo: info, Filenames: &[]string{"tests/*"}},
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestRuff_Ruff_FormatFiles_CheckAll(t *testing.T) {
	exec := test.NewExecutable(t)

	ruff := tool.NewRuff(exec.LazyExecutable("format"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("format", []string{"format", "--check"}).
		Return(nil)

	err := ruff.FormatFiles(
		t.Context(),
		internal.ProjectFormatterOptions{ProjectInfo: info, CheckOnly: true},
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestRuff_Ruff_FormatFiles_Check(t *testing.T) {
	exec := test.NewExecutable(t)

	ruff := tool.NewRuff(exec.LazyExecutable("format"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("format", []string{"format", "--check", "tests/*", "/path/to/file.py"}).
		Return(nil)

	err := ruff.FormatFiles(
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
