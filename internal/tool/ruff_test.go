package tool_test

import (
	"context"
	"testing"

	"cdt/internal"
	"cdt/internal/test"
	"cdt/internal/tool"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRuff_DetectRuff(t *testing.T) {
	ruff := test.RunDetectToolLastFound(t, ruffDetect, []string{"ruff"})

	assert.Equal(t, tool.IdRuff, ruff.Id())
}

func TestRuff_DetectRuff_NotFound(t *testing.T) {
	ruff := test.RunDetectToolNotFound(t, ruffDetect, []string{"ruff"})

	assert.Equal(t, tool.IdRuff, ruff.Id())
}

func TestRuff_DetectRuff_Config(t *testing.T) {
	ruff := test.RunDetectToolConfig(t, ruffDetect, "ruff")

	assert.Equal(t, tool.IdRuff, ruff.Id())
}

func TestRuff_Ruff_LintFiles_All(t *testing.T) {
	test.RunLintFilesSuccess(
		t,
		ruffBuildLinter,
		internal.ProjectLinterOptions{},
		[]string{"check"},
	)
}

func TestRuff_Ruff_LintFiles(t *testing.T) {
	test.RunLintFilesSuccess(
		t,
		ruffBuildLinter,
		internal.ProjectLinterOptions{
			Filenames: &[]string{"file.py", "/path/to/file2.py"},
		},
		[]string{"check", "file.py", "/path/to/file2.py"},
	)
}

func TestRuff_Ruff_LintFiles_OutputFormat_Raw(t *testing.T) {
	test.RunLintFilesOutputFormatCheck(
		t,
		ruffBuildLinter,
		internal.LintReportFormatRaw,
		[]string{"check"},
		nil,
	)
}

func TestRuff_Ruff_LintFiles_OutputFormat_Unsupported(t *testing.T) {
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
			test.RunLintFilesOutputFormatUnsupported(t, ruffBuildLinter, format, nil)
		})
	}
}

func TestRuff_Ruff_LintFiles_OutputFile(t *testing.T) {
	test.RunLintFilesOutputFile(
		t,
		func(executable func() (*internal.Executable, error)) internal.ProjectLinter {
			return tool.NewRuff(executable)
		},
		[]string{"check"},
		nil,
	)
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

func ruffDetect(ctx context.Context, options tool.DetectOptions) internal.Tool {
	return tool.DetectRuff(ctx, options)
}

func ruffBuildLinter(executable func() (*internal.Executable, error)) internal.ProjectLinter {
	return tool.NewRuff(executable)
}
