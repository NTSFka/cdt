package tool_test

import (
	"cdt/internal/test"
	"context"
	"path/filepath"
	"testing"

	"cdt/internal"
	"cdt/internal/tool"

	"github.com/stretchr/testify/assert"
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
