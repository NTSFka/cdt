package tool_test

import (
	"cdt/internal"
	"cdt/internal/test"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Run a test that verifies linting files with an output report a file.
func runTestLintFilesOutputFile(
	t *testing.T,
	build func(executable func() (*internal.Executable, error)) internal.ProjectLinter,
	toolArgs []string,
	filenames *[]string,
) {
	exec := test.NewExecutable(t)

	linter := build(exec.LazyExecutable("lint"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRunOutput("lint", toolArgs, "ok").
		Return(nil)

	outputFilename := t.TempDir() + "/report.txt"

	err := linter.LintFiles(
		t.Context(),
		internal.ProjectLinterOptions{
			ProjectInfo: info,
			Filenames:   filenames,
			Output: internal.OutputOptions[internal.LintReportFormat]{
				Filename: &outputFilename,
			},
		},
	)
	require.NoError(t, err)
	require.FileExists(t, outputFilename)

	outputContent, err := os.ReadFile(outputFilename)
	require.NoError(t, err)
	assert.Equal(t, "ok", string(outputContent))

	exec.AssertExpectations(t)
}
