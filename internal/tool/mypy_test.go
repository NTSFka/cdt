package tool_test

import (
	"testing"

	"cdt/internal"
	"cdt/internal/test"
	"cdt/internal/tool"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMyPy_DetectMyPy(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("mypy").
		Return(env.NewExecutable("/bin/mypy"), nil)

	mypy := tool.DetectMyPy(t.Context(), tool.DetectOptions{Environment: env})
	assert.NotNil(t, mypy)
	assert.Equal(t, "mypy", mypy.Id())
	assert.True(t, mypy.IsAvailable())

	if executable := mypy.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/mypy", executable.Path)
	}

	env.AssertExpectations(t)
}

func TestMyPy_DetectMyPy_NotFound(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("mypy").
		Return(nil, nil)

	mypy := tool.DetectMyPy(t.Context(), tool.DetectOptions{Environment: env})
	assert.NotNil(t, mypy)
	assert.Equal(t, "mypy", mypy.Id())
	assert.False(t, mypy.IsAvailable())
	assert.Nil(t, mypy.Executable())

	env.AssertExpectations(t)
}

func TestMyPy_DetectMyPy_Config(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("mypy-2").
		Return(env.NewExecutable("/bin/mypy"), nil)

	mypy := tool.DetectMyPy(t.Context(), tool.DetectOptions{
		Environment: env,
		ToolsPaths:  map[string]string{"mypy": "mypy-2"},
	})
	assert.NotNil(t, mypy)
	assert.Equal(t, "mypy", mypy.Id())
	assert.True(t, mypy.IsAvailable())

	if executable := mypy.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/mypy", executable.Path)
	}

	env.AssertExpectations(t)
}

func TestMyPy_MyPy_LintFiles_All(t *testing.T) {
	exec := test.NewExecutable(t)

	mypy := tool.NewMyPy(exec.LazyExecutable("lint"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("lint", []string{"*.py"}).
		Return(nil)

	err := mypy.LintFiles(t.Context(), internal.ProjectLinterOptions{ProjectInfo: info})
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestMyPy_MyPy_LintFiles(t *testing.T) {
	exec := test.NewExecutable(t)

	mypy := tool.NewMyPy(exec.LazyExecutable("lint"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("lint", []string{"file.py", "/path/to/file2.py"}).
		Return(nil)

	err := mypy.LintFiles(
		t.Context(),
		internal.ProjectLinterOptions{
			ProjectInfo: info,
			Filenames:   &[]string{"file.py", "/path/to/file2.py"},
		},
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestMyPy_MyPy_LintFiles_OutputFormat_Raw(t *testing.T) {
	exec := test.NewExecutable(t)

	linter := tool.NewMyPy(exec.LazyExecutable("lint"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("lint", []string{"*.py"}).
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

func TestMyPy_MyPy_LintFiles_OutputFormat_Unsupported(t *testing.T) {
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

			linter := tool.NewMyPy(exec.LazyExecutable("lint"))

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

func TestMyPy_MyPy_LintFiles_OutputFile(t *testing.T) {
	runTestLintFilesOutputFile(
		t,
		func(executable func() (*internal.Executable, error)) internal.ProjectLinter {
			return tool.NewMyPy(executable)
		},
		[]string{"*.py"},
		nil,
	)
}
