package tool_test

import (
	"errors"
	"testing"

	"cdt/internal"
	"cdt/internal/test"
	"cdt/internal/tool"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGolangCILint_DetectGolangCILint(t *testing.T) {
	env := test.NewEnvironment(t)
	env.OnFindExecutable("golangci-lint").
		Return(env.NewExecutable("/bin/golangci-lint"), nil)

	golangCILint := tool.DetectGolangCILint(t.Context(), tool.DetectOptions{Environment: env})
	assert.NotNil(t, golangCILint)
	assert.Equal(t, "golangci-lint", golangCILint.Id())
	assert.True(t, golangCILint.IsAvailable())

	if executable := golangCILint.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/golangci-lint", executable.Path)
	}

	env.AssertExpectations(t)
}

func TestGolangCILint_DetectGolangCILint_NotFound(t *testing.T) {
	env := test.NewEnvironment(t)
	env.OnFindExecutable("golangci-lint").
		Return(nil, nil)

	golangCILint := tool.DetectGolangCILint(t.Context(), tool.DetectOptions{Environment: env})
	assert.NotNil(t, golangCILint)
	assert.Equal(t, "golangci-lint", golangCILint.Id())
	assert.False(t, golangCILint.IsAvailable())
	assert.Nil(t, golangCILint.Executable())

	env.AssertExpectations(t)
}

func TestGolangCILint_DetectGolangCILint_Config(t *testing.T) {
	env := test.NewEnvironment(t)
	env.OnFindExecutable("golangci-lint-1").
		Return(env.NewExecutable("/bin/golangci-lint"), nil)

	golangCILint := tool.DetectGolangCILint(t.Context(), tool.DetectOptions{
		Environment: env,
		ToolsPaths:  map[string]string{"golangci-lint": "golangci-lint-1"},
	})
	assert.NotNil(t, golangCILint)
	assert.Equal(t, "golangci-lint", golangCILint.Id())
	assert.True(t, golangCILint.IsAvailable())

	if executable := golangCILint.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/golangci-lint", executable.Path)
	}

	env.AssertExpectations(t)
}

func TestGolangCILint_GolangCILint_LintFiles_All(t *testing.T) {
	exec := test.NewExecutable(t)

	golangCILint := tool.NewGolangCILint(exec.LazyExecutable("lint"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("lint", []string{"run"}).
		Return(nil)

	err := golangCILint.LintFiles(t.Context(), internal.ProjectLinterOptions{ProjectInfo: info})
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestGolangCILint_GolangCILint_LintFiles(t *testing.T) {
	exec := test.NewExecutable(t)

	golangCILint := tool.NewGolangCILint(exec.LazyExecutable("lint"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("lint", []string{"run", "mod1"}).
		Return(nil)

	err := golangCILint.LintFiles(
		t.Context(),
		internal.ProjectLinterOptions{ProjectInfo: info, Filenames: &[]string{"mod1"}},
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestGolangCILint_GolangCILint_LintFiles_OutputFormat_Raw(t *testing.T) {
	exec := test.NewExecutable(t)

	golangCILint := tool.NewGolangCILint(exec.LazyExecutable("lint"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("lint", []string{"run"}).
		Return(nil)

	err := golangCILint.LintFiles(
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

func TestGolangCILint_GolangCILint_LintFiles_OutputFormat_Unsupported(t *testing.T) {
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

			golangCILint := tool.NewGolangCILint(exec.LazyExecutable("lint"))

			info := internal.ProjectInfo{Directory: "."}

			err := golangCILint.LintFiles(
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

func TestGolangCILint_FormatFiles_All(t *testing.T) {
	exec := test.NewExecutable(t)

	golangCILint := tool.NewGolangCILint(exec.LazyExecutable("lint"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("lint", []string{"fmt"}).
		Return(nil)

	err := golangCILint.FormatFiles(
		t.Context(),
		internal.ProjectFormatterOptions{ProjectInfo: info},
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestGolangCILint_FormatFiles_All_Failed(t *testing.T) {
	exec := test.NewExecutable(t)

	golangCILint := tool.NewGolangCILint(exec.LazyExecutable("lint"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("lint", []string{"fmt"}).
		Return(errors.New("failed"))

	err := golangCILint.FormatFiles(
		t.Context(),
		internal.ProjectFormatterOptions{ProjectInfo: info},
	)
	require.EqualError(t, err, "failed")

	exec.AssertExpectations(t)
}

func TestGolangCILint_FormatFiles(t *testing.T) {
	exec := test.NewExecutable(t)

	golangCILint := tool.NewGolangCILint(exec.LazyExecutable("lint"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("lint", []string{"fmt", "file1"}).
		Return(nil)

	err := golangCILint.FormatFiles(
		t.Context(),
		internal.ProjectFormatterOptions{ProjectInfo: info, Filenames: &[]string{"file1"}},
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestGolangCILint_FormatFiles_Failed(t *testing.T) {
	exec := test.NewExecutable(t)

	goTool := tool.NewGolangCILint(exec.LazyExecutable("lint"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("lint", []string{"fmt", "file1"}).
		Return(errors.New("failed"))

	err := goTool.FormatFiles(
		t.Context(),
		internal.ProjectFormatterOptions{ProjectInfo: info, Filenames: &[]string{"file1"}},
	)
	require.EqualError(t, err, "failed")

	exec.AssertExpectations(t)
}

func TestGolangCILint_FormatFiles_CheckAll(t *testing.T) {
	exec := test.NewExecutable(t)

	golangCILint := tool.NewGolangCILint(exec.LazyExecutable("lint"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("lint", []string{"fmt", "--diff"}).
		Return(nil)

	err := golangCILint.FormatFiles(
		t.Context(),
		internal.ProjectFormatterOptions{ProjectInfo: info, CheckOnly: true},
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestGolangCILint_FormatFiles_CheckAll_Failed(t *testing.T) {
	exec := test.NewExecutable(t)

	golangCILint := tool.NewGolangCILint(exec.LazyExecutable("lint"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("lint", []string{"fmt", "--diff"}).
		Return(errors.New("failed"))

	err := golangCILint.FormatFiles(
		t.Context(),
		internal.ProjectFormatterOptions{ProjectInfo: info, CheckOnly: true},
	)
	require.EqualError(t, err, "failed")

	exec.AssertExpectations(t)
}

func TestGolangCILint_FormatFiles_Check(t *testing.T) {
	exec := test.NewExecutable(t)

	golangCILint := tool.NewGolangCILint(exec.LazyExecutable("lint"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("lint", []string{"fmt", "--diff", "file1"}).
		Return(nil)

	err := golangCILint.FormatFiles(
		t.Context(),
		internal.ProjectFormatterOptions{
			ProjectInfo: info,
			CheckOnly:   true,
			Filenames:   &[]string{"file1"},
		},
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestGolangCILint_FormatFiles_Check_Failed(t *testing.T) {
	exec := test.NewExecutable(t)

	goTool := tool.NewGolangCILint(exec.LazyExecutable("lint"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("lint", []string{"fmt", "--diff", "file1"}).
		Return(errors.New("failed"))

	err := goTool.FormatFiles(
		t.Context(),
		internal.ProjectFormatterOptions{
			ProjectInfo: info,
			CheckOnly:   true,
			Filenames:   &[]string{"file1"},
		},
	)
	require.EqualError(t, err, "failed")

	exec.AssertExpectations(t)
}
