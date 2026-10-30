package tool_test

import (
	"bytes"
	"cdt/internal/reporter"
	"cdt/internal/test"
	"context"
	"path/filepath"
	"testing"

	"cdt/internal"
	"cdt/internal/tool"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestComposerDependencyAnalyzer_Detect_Composer(t *testing.T) {
	composerDependencyAnalyzer := test.RunDetectToolLastFound(
		t,
		composerDependencyAnalyzerDetect,
		[]string{
			filepath.Join("vendor", "bin", "composer-dependency-analyzer"),
		},
	)

	assert.Equal(t, tool.IdComposerDependencyAnalyzer, composerDependencyAnalyzer.Id())
}

func TestComposerDependencyAnalyzer_Detect_System(t *testing.T) {
	composerDependencyAnalyzer := test.RunDetectToolLastFound(
		t,
		composerDependencyAnalyzerDetect,
		[]string{
			filepath.Join("vendor", "bin", "composer-dependency-analyzer"),
			"composer-dependency-analyzer",
		},
	)

	assert.Equal(t, tool.IdComposerDependencyAnalyzer, composerDependencyAnalyzer.Id())
}

func TestComposerDependencyAnalyzer_Detect_NotFound(t *testing.T) {
	composerDependencyAnalyzer := test.RunDetectToolNotFound(
		t,
		composerDependencyAnalyzerDetect,
		[]string{
			filepath.Join("vendor", "bin", "composer-dependency-analyzer"),
			"composer-dependency-analyzer",
		},
	)

	assert.Equal(t, tool.IdComposerDependencyAnalyzer, composerDependencyAnalyzer.Id())
}

func TestComposerDependencyAnalyzer_Detect_Config(t *testing.T) {
	composerDependencyAnalyzer := test.RunDetectToolConfig(
		t,
		composerDependencyAnalyzerDetect,
		"composer-dependency-analyzer",
	)

	assert.Equal(t, tool.IdComposerDependencyAnalyzer, composerDependencyAnalyzer.Id())
}

func TestComposerDependencyAnalyzer_AnalyzeFiles_All(t *testing.T) {
	test.InvokeAnalyzeFilesSuccess(
		t,
		composerDependencyAnalyzerAnalyzer,
		internal.ProjectAnalyzerOptions{},
		[]string{},
	)
}

func TestComposerDependencyAnalyzer_AnalyzeFiles_All_Failed(t *testing.T) {
	test.InvokeAnalyzeFilesFailed(
		t,
		composerDependencyAnalyzerAnalyzer,
		internal.ProjectAnalyzerOptions{},
		[]string{},
	)
}

func TestComposerDependencyAnalyzer_AnalyzeFiles(t *testing.T) {
	fileNames := []string{"file.php", "/path/to/file2.php"}

	test.InvokeAnalyzeFilesSuccess(
		t,
		composerDependencyAnalyzerAnalyzer,
		internal.ProjectAnalyzerOptions{Filenames: &fileNames},
		fileNames,
	)
}

func TestComposerDependencyAnalyzer_AnalyzeFiles_Failed(t *testing.T) {
	fileNames := []string{"file.php", "/path/to/file2.php"}

	test.InvokeAnalyzeFilesFailed(
		t,
		composerDependencyAnalyzerAnalyzer,
		internal.ProjectAnalyzerOptions{Filenames: &fileNames},
		fileNames,
	)
}

func TestComposerDependencyAnalyzer_AnalyzeFiles_Reporter_NoErrors(t *testing.T) {
	result := bytes.Buffer{}

	test.InvokeAnalyzeFilesReporter(
		t,
		composerDependencyAnalyzerAnalyzer,
		".",
		reporter.NewSimpleSarifReporter(&result),
		[]string{"--format", "junit"},
		"<?xml version=\"1.0\" encoding=\"UTF-8\"?>",
		func(err error) {
			require.NoError(t, err)
		},
	)

	assert.Equal(t, "\nTotal issues: 0\n\n", result.String())
}

func TestComposerDependencyAnalyzer_AnalyzeFiles_Reporter_SomeErrors(t *testing.T) {
	result := bytes.Buffer{}

	test.InvokeAnalyzeFilesReporter(
		t,
		composerDependencyAnalyzerAnalyzer,
		".",
		reporter.NewSimpleSarifReporter(&result),
		[]string{"--format", "junit"},
		`<?xml version="1.0" encoding="UTF-8"?>
		<testsuites>
		  <testsuite name="shadow dependencies" failures="4">
			<testcase name="package1">
			  <failure message="Class1">www/index.php:10</failure>
			</testcase>
		  </testsuite>
		  <!-- showing only first example failure usage -->
		</testsuites>`,
		func(err error) {
			require.NoError(t, err)
		},
	)

	assert.NotEmpty(t, result.String())
}

func TestComposerDependencyAnalyzer_AnalyzeFiles_Reporter_NoOutput(t *testing.T) {
	result := bytes.Buffer{}

	test.InvokeAnalyzeFilesReporter(
		t,
		composerDependencyAnalyzerAnalyzer,
		".",
		reporter.NewSimpleSarifReporter(&result),
		[]string{"--format", "junit"},
		"",
		func(err error) {
			require.NoError(t, err)
		},
	)

	assert.Equal(t, "\nTotal issues: 0\n\n", result.String())
}

func TestComposerDependencyAnalyzer_AnalyzeFiles_Reporter_InvalidXml(t *testing.T) {
	result := bytes.Buffer{}

	test.InvokeAnalyzeFilesReporter(
		t,
		composerDependencyAnalyzerAnalyzer,
		".",
		reporter.NewSimpleSarifReporter(&result),
		[]string{"--format", "junit"},
		".",
		func(err error) {
			// The tool returns invalid XML when no issue is found, so it's not possible to distinguish between
			// no issue and invalid XML
			require.NoError(t, err)
		},
	)

	assert.Equal(t, "\nTotal issues: 0\n\n", result.String())
}

func TestComposerDependencyAnalyzer_AnalyzeFiles_OutputFile(t *testing.T) {
	test.InvokeAnalyzeFilesOutputFile(
		t,
		composerDependencyAnalyzerAnalyzer,
		".",
		[]string{},
		nil,
	)
}

func composerDependencyAnalyzerDetect(
	ctx context.Context,
	options internal.DetectOptions,
) internal.Tool {
	return tool.DetectComposerDependencyAnalyzer(ctx, options)
}

func composerDependencyAnalyzerAnalyzer(
	executable func() (*internal.Executable, error),
) internal.ProjectAnalyzer {
	return tool.NewComposerDependencyAnalyzer(executable)
}
