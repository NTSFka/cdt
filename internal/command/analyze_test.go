package command_test

import (
	"context"
	"errors"
	"testing"

	"cdt/internal"
	"cdt/internal/command"
	"cdt/internal/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func runAnalyze(ctx context.Context, analyzer internal.ProjectAnalyzer, args ...string) error {
	return test.RunCommand(ctx, command.NewAnalyzeCommand(), internal.Context{
		Project: internal.Project{
			Workflow: internal.Workflow{
				Analyzer: analyzer,
			},
		},
	}, args...)
}

func runAnalyzeTool(ctx context.Context, analyzer internal.Tool, args ...string) error {
	return test.RunCommand(ctx, command.NewAnalyzeCommand(), internal.Context{
		Tools: []internal.Tool{
			analyzer,
		},
	}, args...)
}

func TestAnalyze_Analyze_CannotBeAnalyzed(t *testing.T) {
	err := runAnalyze(t.Context(), nil)

	require.Error(t, err)
	assert.Equal(t, "project doesn't support analyzing", err.Error())
}

func TestAnalyze_AnalyzeFiles_All_Success(t *testing.T) {
	analyzer := test.NewProjectAnalyzer(t)
	analyzer.On("AnalyzeFiles", mock.Anything, mock.MatchedBy(func(options internal.ProjectAnalyzerOptions) bool {
		return assert.Nil(t, options.Filenames)
	})).
		Return(nil)

	err := runAnalyze(t.Context(), analyzer)

	require.NoError(t, err)
	analyzer.AssertExpectations(t)
}

func TestAnalyze_AnalyzeFiles_All_Failure(t *testing.T) {
	analyzer := test.NewProjectAnalyzer(t)
	analyzer.On("AnalyzeFiles", mock.Anything, mock.Anything).
		Return(errors.New("failed"))

	err := runAnalyze(t.Context(), analyzer)

	require.Error(t, err)
	assert.Equal(t, "failed", err.Error())

	analyzer.AssertExpectations(t)
}

type testAnalyzerTool struct {
	internal.ExecutableTool
	test.ProjectAnalyzer
}

func newTestAnalyzerTool(t *testing.T) *testAnalyzerTool {
	linter := &testAnalyzerTool{
		internal.MakeExecutableTool("tool1", "", "", internal.Tags{}, nil),
		test.ProjectAnalyzer{},
	}
	linter.Test(t)

	return linter
}

func TestAnalyze_Tool_Success(t *testing.T) {
	linter := newTestAnalyzerTool(t)
	linter.On("AnalyzeFiles", mock.Anything, mock.Anything).
		Return(nil)

	err := runAnalyzeTool(t.Context(), linter, "--tool", "tool1")

	require.NoError(t, err)
	linter.AssertExpectations(t)
}

func TestAnalyze_Tool_Failed(t *testing.T) {
	linter := newTestAnalyzerTool(t)
	linter.On("AnalyzeFiles", mock.Anything, mock.Anything).
		Return(errors.New("failed"))

	err := runAnalyzeTool(t.Context(), linter, "--tool", "tool1")

	require.Error(t, err)
	assert.Equal(t, "failed", err.Error())

	linter.AssertExpectations(t)
}

func TestAnalyze_Tool_NotFound(t *testing.T) {
	linter := newTestAnalyzerTool(t)

	err := runAnalyzeTool(t.Context(), linter, "--tool", "tool2")

	require.Error(t, err)
	assert.Equal(t, "tool 'tool2' not found", err.Error())

	linter.AssertExpectations(t)
}

func TestAnalyze_Tool_NotSupported(t *testing.T) {
	linter := struct {
		internal.ExecutableTool
	}{
		internal.MakeExecutableTool("tool1", "", "", internal.Tags{}, nil),
	}

	err := runAnalyzeTool(t.Context(), &linter, "--tool", "tool1")

	require.Error(t, err)
	assert.Equal(t, "tool 'tool1' doesn't support analyzing", err.Error())
}

func TestAnalyze_AnalyzeFiles_Success(t *testing.T) {
	linter := test.NewProjectAnalyzer(t)
	linter.On(
		"AnalyzeFiles",
		mock.Anything,
		mock.MatchedBy(func(opts internal.ProjectAnalyzerOptions) bool {
			return assert.NotNil(t, opts.Filenames) &&
				assert.ElementsMatch(t, []string{"file1", "file2"}, *opts.Filenames)
		}),
	).
		Return(nil)

	err := runAnalyze(t.Context(), linter, "file1", "file2")

	require.NoError(t, err)
	linter.AssertExpectations(t)
}

func TestAnalyze_AnalyzeFiles_Failure(t *testing.T) {
	linter := test.NewProjectAnalyzer(t)
	linter.On(
		"AnalyzeFiles",
		mock.Anything,
		mock.MatchedBy(func(opts internal.ProjectAnalyzerOptions) bool {
			return assert.NotNil(t, opts.Filenames) &&
				assert.ElementsMatch(t, []string{"file1", "file2"}, *opts.Filenames)
		}),
	).
		Return(errors.New("failed"))

	err := runAnalyze(t.Context(), linter, "file1", "file2")

	require.Error(t, err)
	assert.Equal(t, "failed", err.Error())

	linter.AssertExpectations(t)
}

func TestAnalyze_AnalyzeFiles_CustomOutput(t *testing.T) {
	linter := test.NewProjectAnalyzer(t)
	linter.On("AnalyzeFiles", mock.Anything, mock.MatchedBy(func(options internal.ProjectAnalyzerOptions) bool {
		return assert.Equal(t, internal.AnalyzeReportFormatSarif, options.Output.Format) &&
			assert.Nil(t, options.Output.Filename)
	})).
		Return(nil)

	err := runAnalyze(t.Context(), linter, "--output", "sarif")

	require.NoError(t, err)
	linter.AssertExpectations(t)
}

func TestAnalyze_AnalyzeFiles_CustomOutputWithFile(t *testing.T) {
	linter := test.NewProjectAnalyzer(t)
	linter.On("AnalyzeFiles", mock.Anything, mock.MatchedBy(func(options internal.ProjectAnalyzerOptions) bool {
		return assert.Equal(t, internal.AnalyzeReportFormatSarif, options.Output.Format) &&
			assert.NotNil(t, options.Output.Filename) &&
			assert.Equal(t, "report.sarif", *options.Output.Filename)
	})).
		Return(nil)

	err := runAnalyze(t.Context(), linter, "-o", "sarif:report.sarif")

	require.NoError(t, err)
	linter.AssertExpectations(t)
}
