package tool_test

import (
	"context"
	"errors"
	"testing"

	"cdt/internal"
	"cdt/internal/test"
	"cdt/internal/tool"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGolangCILint_DetectGolangCILint(t *testing.T) {
	golangCILint := test.RunDetectToolLastFound(t, golangCILintDetect, []string{"golangci-lint"})

	assert.Equal(t, tool.IdGolangCILint, golangCILint.Id())
}

func TestGolangCILint_DetectGolangCILint_NotFound(t *testing.T) {
	golangCILint := test.RunDetectToolNotFound(t, golangCILintDetect, []string{"golangci-lint"})

	assert.Equal(t, tool.IdGolangCILint, golangCILint.Id())
}

func TestGolangCILint_DetectGolangCILint_Config(t *testing.T) {
	golangCILint := test.RunDetectToolConfig(t, golangCILintDetect, "golangci-lint")

	assert.Equal(t, tool.IdGolangCILint, golangCILint.Id())
}

func TestGolangCILint_GolangCILint_LintFiles_All(t *testing.T) {
	test.RunLintFilesSuccess(
		t,
		golangCILintBuildLinter,
		internal.ProjectLinterOptions{},
		[]string{"run"},
	)
}

func TestGolangCILint_GolangCILint_LintFiles(t *testing.T) {
	test.RunLintFilesSuccess(
		t,
		golangCILintBuildLinter,
		internal.ProjectLinterOptions{Filenames: &[]string{"mod1"}},
		[]string{"run", "mod1"},
	)
}

func TestGolangCILint_GolangCILint_LintFiles_OutputFormat_Raw(t *testing.T) {
	dataSet := []struct {
		format internal.LintReportFormat
		args   []string
	}{
		{internal.LintReportFormatRaw, []string{}},
	}

	for _, data := range dataSet {
		t.Run(string(data.format), func(t *testing.T) {
			test.RunLintFilesOutputFormatCheck(
				t,
				golangCILintBuildLinter,
				data.format,
				append([]string{"run"}, data.args...),
				nil,
			)
		})
	}
}

func TestGolangCILint_GolangCILint_LintFiles_OutputFormat_Unsupported(t *testing.T) {
	dataSet := []internal.LintReportFormat{
		internal.LintReportFormatJson,
		internal.LintReportFormatJUnit,
		internal.LintReportFormatGitHub,
		internal.LintReportFormatGitLab,
		internal.LintReportFormatTeamCity,
		"test-unsupported",
	}

	for _, format := range dataSet {
		t.Run(string(format), func(t *testing.T) {
			test.RunLintFilesOutputFormatUnsupported(t, golangCILintBuildLinter, format, nil)
		})
	}
}

func TestGolangCILint_GolangCILint_LintFiles_OutputFile(t *testing.T) {
	test.RunLintFilesOutputFile(
		t,
		func(executable func() (*internal.Executable, error)) internal.ProjectLinter {
			return tool.NewGolangCILint(executable)
		},
		[]string{"run"},
		nil,
	)
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

func golangCILintDetect(ctx context.Context, options tool.DetectOptions) internal.Tool {
	return tool.DetectGolangCILint(ctx, options)
}

func golangCILintBuildLinter(
	executable func() (*internal.Executable, error),
) internal.ProjectLinter {
	return tool.NewGolangCILint(executable)
}
