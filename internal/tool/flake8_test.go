package tool_test

import (
	"testing"

	"cdt/internal"
	"cdt/internal/test"
	"cdt/internal/tool"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFlake8_DetectFlake8(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("flake8").
		Return(env.NewExecutable("/bin/flake8"), nil)

	flake8 := tool.DetectFlake8(t.Context(), tool.DetectOptions{Environment: env})
	assert.NotNil(t, flake8)
	assert.Equal(t, "flake8", flake8.Id())
	assert.True(t, flake8.IsAvailable())

	if executable := flake8.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/flake8", executable.Path)
	}

	env.AssertExpectations(t)
}

func TestFlake8_DetectFlake8_NotFound(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("flake8").
		Return(nil, nil)

	flake8 := tool.DetectFlake8(t.Context(), tool.DetectOptions{Environment: env})
	assert.NotNil(t, flake8)
	assert.Equal(t, "flake8", flake8.Id())
	assert.False(t, flake8.IsAvailable())
	assert.Nil(t, flake8.Executable())

	env.AssertExpectations(t)
}

func TestFlake8_DetectFlake8_Config(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("flake8-2").
		Return(env.NewExecutable("/bin/flake8"), nil)

	flake8 := tool.DetectFlake8(t.Context(), tool.DetectOptions{
		Environment: env,
		ToolsPaths:  map[string]string{"flake8": "flake8-2"},
	})
	assert.NotNil(t, flake8)
	assert.Equal(t, "flake8", flake8.Id())
	assert.True(t, flake8.IsAvailable())

	if executable := flake8.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/flake8", executable.Path)
	}

	env.AssertExpectations(t)
}

func TestFlake8_Flake8_LintFiles_All(t *testing.T) {
	exec := test.NewExecutable(t)

	flake8 := tool.NewFlake8(exec.LazyExecutable("lint"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("lint", []string{}).
		Return(nil)

	err := flake8.LintFiles(t.Context(), internal.ProjectLinterOptions{ProjectInfo: info})
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestFlake8_Flake8_LintFiles(t *testing.T) {
	exec := test.NewExecutable(t)

	flake8 := tool.NewFlake8(exec.LazyExecutable("lint"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("lint", []string{"file.py", "/path/to/file2.py"}).
		Return(nil)

	err := flake8.LintFiles(
		t.Context(),
		internal.ProjectLinterOptions{
			ProjectInfo: info,
			Filenames:   &[]string{"file.py", "/path/to/file2.py"},
		},
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestFlake8_Flake8_LintFiles_OutputFormat_Raw(t *testing.T) {
	exec := test.NewExecutable(t)

	linter := tool.NewFlake8(exec.LazyExecutable("lint"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("lint", []string{}).
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

func TestFlake8_Flake8_LintFiles_OutputFormat_Unsupported(t *testing.T) {
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

			linter := tool.NewFlake8(exec.LazyExecutable("lint"))

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

func TestFlake8_Flake8_LintFiles_OutputFile(t *testing.T) {
	runTestLintFilesOutputFile(
		t,
		func(executable func() (*internal.Executable, error)) internal.ProjectLinter {
			return tool.NewFlake8(executable)
		},
		[]string{},
		nil,
	)
}
