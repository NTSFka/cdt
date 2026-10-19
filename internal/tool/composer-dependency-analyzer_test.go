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

func TestComposerDependencyAnalyzer_DetectComposerDependencyAnalyzer_Composer(t *testing.T) {
	composerDependencyAnalyzer := test.RunDetectToolLastFound(
		t,
		composerDependencyAnalyzerDetect,
		[]string{
			filepath.Join("vendor", "bin", "composer-dependency-analyzer"),
		},
	)

	assert.Equal(t, tool.IdComposerDependencyAnalyzer, composerDependencyAnalyzer.Id())
}

func TestComposerDependencyAnalyzer_DetectComposerDependencyAnalyzer_System(t *testing.T) {
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

func TestComposerDependencyAnalyzer_DetectComposerDependencyAnalyzer_NotFound(t *testing.T) {
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

func TestComposerDependencyAnalyzer_DetectComposerDependencyAnalyzer_Config(t *testing.T) {
	composerDependencyAnalyzer := test.RunDetectToolConfig(
		t,
		composerDependencyAnalyzerDetect,
		"composer-dependency-analyzer",
	)

	assert.Equal(t, tool.IdComposerDependencyAnalyzer, composerDependencyAnalyzer.Id())
}

func TestComposerDependencyAnalyzer_ComposerDependencyAnalyzer_LintFiles_All(t *testing.T) {
	test.InvokeLintFilesSuccess(
		t,
		composerDependencyAnalyzerLinter,
		internal.ProjectLinterOptions{},
		[]string{},
	)
}

func TestComposerDependencyAnalyzer_ComposerDependencyAnalyzer_LintFiles_All_Failed(t *testing.T) {
	test.InvokeLintFilesFailed(
		t,
		composerDependencyAnalyzerLinter,
		internal.ProjectLinterOptions{},
		[]string{},
	)
}

func TestComposerDependencyAnalyzer_ComposerDependencyAnalyzer_LintFiles(t *testing.T) {
	fileNames := []string{"file.php", "/path/to/file2.php"}

	test.InvokeLintFilesSuccess(
		t,
		composerDependencyAnalyzerLinter,
		internal.ProjectLinterOptions{Filenames: &fileNames},
		fileNames,
	)
}

func TestComposerDependencyAnalyzer_ComposerDependencyAnalyzer_LintFiles_Failed(t *testing.T) {
	fileNames := []string{"file.php", "/path/to/file2.php"}

	test.InvokeLintFilesFailed(
		t,
		composerDependencyAnalyzerLinter,
		internal.ProjectLinterOptions{Filenames: &fileNames},
		fileNames,
	)
}

func composerDependencyAnalyzerDetect(
	ctx context.Context,
	options internal.DetectOptions,
) internal.Tool {
	return tool.DetectComposerDependencyAnalyzer(ctx, options)
}

func composerDependencyAnalyzerLinter(
	executable func() (*internal.Executable, error),
) internal.ProjectLinter {
	return tool.NewComposerDependencyAnalyzer(executable)
}
