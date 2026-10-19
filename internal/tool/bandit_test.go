package tool_test

import (
	"context"
	"testing"

	"cdt/internal"
	"cdt/internal/test"
	"cdt/internal/tool"

	"github.com/stretchr/testify/assert"
)

func TestBandit_DetectBandit(t *testing.T) {
	bandit := test.RunDetectToolLastFound(t, banditDetect, []string{"bandit"})

	assert.Equal(t, tool.IdBandit, bandit.Id())
}

func TestBandit_DetectBandit_NotFound(t *testing.T) {
	bandit := test.RunDetectToolNotFound(t, banditDetect, []string{"bandit"})

	assert.Equal(t, tool.IdBandit, bandit.Id())
}

func TestBandit_DetectBandit_Config(t *testing.T) {
	bandit := test.RunDetectToolConfig(t, banditDetect, "bandit")

	assert.Equal(t, tool.IdBandit, bandit.Id())
}

func TestBandit_Bandit_LintFiles_All(t *testing.T) {
	test.InvokeLintFilesSuccess(
		t,
		banditLinter,
		internal.ProjectLinterOptions{},
		[]string{"*"},
	)
}

func TestBandit_Bandit_LintFiles_All_Failed(t *testing.T) {
	test.InvokeLintFilesFailed(
		t,
		banditLinter,
		internal.ProjectLinterOptions{},
		[]string{"*"},
	)
}

func TestBandit_Bandit_LintFiles(t *testing.T) {
	fileNames := []string{"file.py", "/path/to/file2.py"}

	test.InvokeLintFilesSuccess(
		t,
		banditLinter,
		internal.ProjectLinterOptions{Filenames: &fileNames},
		fileNames,
	)
}

func TestBandit_Bandit_LintFiles_Failed(t *testing.T) {
	fileNames := []string{"file.py", "/path/to/file2.py"}

	test.InvokeLintFilesFailed(
		t,
		banditLinter,
		internal.ProjectLinterOptions{Filenames: &fileNames},
		fileNames,
	)
}

func TestBandit_Bandit_LintFiles_OutputFormat_Raw(t *testing.T) {
	test.InvokeLintFilesOutputFormatCheck(
		t,
		banditLinter,
		".",
		internal.LintReportFormatRaw,
		[]string{"*"},
		nil,
	)
}

func TestBandit_Bandit_LintFiles_OutputFormat_Unsupported(t *testing.T) {
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
			test.InvokeLintFilesOutputFormatUnsupported(t, banditLinter, ".", format, nil)
		})
	}
}

func TestBandit_Bandit_LintFiles_OutputFile(t *testing.T) {
	test.InvokeLintFilesOutputFile(t, banditLinter, ".", []string{"*"}, nil)
}

func banditDetect(ctx context.Context, options internal.DetectOptions) internal.Tool {
	return tool.DetectBandit(ctx, options)
}

func banditLinter(executable func() (*internal.Executable, error)) internal.ProjectLinter {
	return tool.NewBandit(executable)
}
