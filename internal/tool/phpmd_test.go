package tool_test

import (
	"cdt/internal/test"
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"cdt/internal"
	"cdt/internal/tool"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const phpMDVersion2 = "PHPMD 2.0.0"
const phpMDVersion3 = "PHPMD 3.0.0"

func TestPHPMD_DetectPHPMD_Phar(t *testing.T) {
	phpMD := test.RunDetectToolLastFound(
		t,
		phpMDDetect,
		[]string{
			"phpmd.phar",
		},
	)

	assert.Equal(t, tool.IdPHPMD, phpMD.Id())
}

func TestPHPMD_DetectPHPMD_Composer(t *testing.T) {
	phpMD := test.RunDetectToolLastFound(
		t,
		phpMDDetect,
		[]string{
			"phpmd.phar",
			filepath.Join("vendor", "bin", "phpmd"),
		},
	)

	assert.Equal(t, tool.IdPHPMD, phpMD.Id())
}

func TestPHPMD_DetectPHPMD_System(t *testing.T) {
	phpMD := test.RunDetectToolLastFound(
		t,
		phpMDDetect,
		[]string{
			"phpmd.phar",
			filepath.Join("vendor", "bin", "phpmd"),
			"phpmd",
		},
	)

	assert.Equal(t, tool.IdPHPMD, phpMD.Id())
}

func TestPHPMD_DetectPHPMD_NotFound(t *testing.T) {
	phpMD := test.RunDetectToolNotFound(
		t,
		phpMDDetect,
		[]string{
			"phpmd.phar",
			filepath.Join("vendor", "bin", "phpmd"),
			"phpmd",
		},
	)

	assert.Equal(t, tool.IdPHPMD, phpMD.Id())
}

func TestPHPMD_DetectPHPMD_Config(t *testing.T) {
	phpMD := test.RunDetectToolConfig(t, phpMDDetect, "phpmd")

	assert.Equal(t, tool.IdPHPMD, phpMD.Id())
}

func TestPHPMD_PHPMD_Version_Failed(t *testing.T) {
	expectedErr := errors.New("version error")

	test.InvokeLintFiles(
		t,
		phpMDLinter,
		internal.ProjectLinterOptions{
			ProjectInfo: internal.ProjectInfo{
				Directory: ".",
			},
		},
		func(exec *test.Executable, executablePath string) {
			exec.OnRunOutput(executablePath, []string{"--version"}, phpMDVersion2).
				Return(expectedErr)
		},
		func(err error) {
			require.ErrorIs(t, err, expectedErr)
		},
	)
}

func TestPHPMD_PHPMD_Version_Unknown(t *testing.T) {
	test.InvokeLintFiles(
		t,
		phpMDLinter,
		internal.ProjectLinterOptions{
			ProjectInfo: internal.ProjectInfo{
				Directory: ".",
			},
		},
		func(exec *test.Executable, executablePath string) {
			exec.OnRunOutput(executablePath, []string{"--version"}, "PHPMD 1.0.0").
				Return(nil)
		},
		func(err error) {
			require.EqualError(t, err, "unsupported PHPMD version: PHPMD 1.0.0")
		},
	)
}

func TestPHPMD_PHPMD_2_LintFiles_All_Success(t *testing.T) {
	dataSet := []string{
		"phpmd.xml",
		"phpmd.ruleset.xml",
	}

	for _, fileName := range dataSet {
		t.Run(fileName, func(t *testing.T) {
			directory := t.TempDir()

			test.CreateFile(t, directory, fileName)

			test.InvokeLintFiles(
				t,
				phpMDLinter,
				internal.ProjectLinterOptions{
					ProjectInfo: internal.ProjectInfo{
						Directory: directory,
					},
				},
				func(exec *test.Executable, executablePath string) {
					exec.OnRunOutput(executablePath, []string{"--version"}, phpMDVersion2).
						Return(nil)

					exec.OnRun(executablePath, []string{directory, "ansi", fileName}).
						Return(nil)
				},
				func(err error) {
					require.NoError(t, err)
				},
			)
		})
	}
}

func TestPHPMD_PHPMD_2_LintFiles_All_Directories(t *testing.T) {
	directory := t.TempDir()

	test.CreateFile(t, directory, "phpmd.xml")
	test.CreateDirectory(t, directory, "src")
	test.CreateDirectory(t, directory, "tests")

	test.InvokeLintFiles(
		t,
		phpMDLinter,
		internal.ProjectLinterOptions{
			ProjectInfo: internal.ProjectInfo{
				Directory: directory,
			},
		},
		func(exec *test.Executable, executablePath string) {
			exec.OnRunOutput(executablePath, []string{"--version"}, phpMDVersion2).
				Return(nil)

			exec.OnRun(executablePath, []string{"src,tests", "ansi", "phpmd.xml"}).
				Return(nil)
		},
		func(err error) {
			require.NoError(t, err)
		},
	)
}

func TestPHPMD_PHPMD_2_LintFiles_All_Failed(t *testing.T) {
	directory := t.TempDir()

	test.CreateFile(t, directory, "phpmd.xml")

	expectedErr := errors.New("failed")

	test.InvokeLintFiles(
		t,
		phpMDLinter,
		internal.ProjectLinterOptions{
			ProjectInfo: internal.ProjectInfo{
				Directory: directory,
			},
		},
		func(exec *test.Executable, executablePath string) {
			exec.OnRunOutput(executablePath, []string{"--version"}, phpMDVersion2).
				Return(nil)

			exec.OnRun(executablePath, []string{directory, "ansi", "phpmd.xml"}).
				Return(expectedErr)
		},
		func(err error) {
			require.ErrorIs(t, err, expectedErr)
		},
	)
}

func TestPHPMD_PHPMD_2_LintFiles_All_MissingRulesFile(t *testing.T) {
	test.InvokeLintFiles(
		t,
		phpMDLinter,
		internal.ProjectLinterOptions{},
		func(exec *test.Executable, executablePath string) {
			exec.OnRunOutput(executablePath, []string{"--version"}, phpMDVersion2).
				Return(nil)
		},
		func(err error) {
			require.EqualError(
				t,
				err,
				"ruleset required, create a phpmd.xml or pass additional argument for ruleset",
			)
		},
	)
}

func TestPHPMD_PHPMD_2_LintFiles_All_CustomRulesFile(t *testing.T) {
	directory := t.TempDir()

	test.InvokeLintFiles(
		t,
		phpMDLinter,
		internal.ProjectLinterOptions{
			ProjectInfo: internal.ProjectInfo{
				Directory: directory,
			},
			ExtraArgs: []string{"phpmd.rules.xml"},
		},
		func(exec *test.Executable, executablePath string) {
			exec.OnRunOutput(executablePath, []string{"--version"}, phpMDVersion2).
				Return(nil)

			exec.OnRun(executablePath, []string{directory, "ansi", "phpmd.rules.xml"}).
				Return(nil)
		},
		func(err error) {
			require.NoError(t, err)
		},
	)
}

func TestPHPMD_PHPMD_2_LintFiles(t *testing.T) {
	directory := t.TempDir()

	test.CreateFile(t, directory, "phpmd.xml")

	fileNames := []string{"file.php", "/path/to/file2.php"}

	test.InvokeLintFiles(
		t,
		phpMDLinter,
		internal.ProjectLinterOptions{
			ProjectInfo: internal.ProjectInfo{
				Directory: directory,
			},
			Filenames: &fileNames,
		},
		func(exec *test.Executable, executablePath string) {
			exec.OnRunOutput(executablePath, []string{"--version"}, phpMDVersion2).
				Return(nil)

			exec.OnRun(executablePath, []string{strings.Join(fileNames, ","), "ansi", "phpmd.xml"}).
				Return(nil)
		},
		func(err error) {
			require.NoError(t, err)
		},
	)
}

func TestPHPMD_PHPMD_2_LintFiles_Failed(t *testing.T) {
	directory := t.TempDir()

	test.CreateFile(t, directory, "phpmd.xml")

	fileNames := []string{"file.php", "/path/to/file2.php"}

	expectedErr := errors.New("failed")

	test.InvokeLintFiles(
		t,
		phpMDLinter,
		internal.ProjectLinterOptions{
			ProjectInfo: internal.ProjectInfo{
				Directory: directory,
			},
			Filenames: &fileNames,
		},
		func(exec *test.Executable, executablePath string) {
			exec.OnRunOutput(executablePath, []string{"--version"}, phpMDVersion2).
				Return(nil)

			exec.OnRun(executablePath, []string{strings.Join(fileNames, ","), "ansi", "phpmd.xml"}).
				Return(expectedErr)
		},
		func(err error) {
			require.ErrorIs(t, err, expectedErr)
		},
	)
}

func TestPHPMD_PHPMD_2_LintFiles_OutputFormat(t *testing.T) {
	directory := t.TempDir()

	test.CreateFile(t, directory, "phpmd.xml")

	dataSet := []struct {
		format internal.LintReportFormat
		arg    string
	}{
		{internal.LintReportFormatRaw, "ansi"},
		{internal.LintReportFormatJson, "json"},
		{internal.LintReportFormatGitHub, "github"},
		{internal.LintReportFormatGitLab, "gitlab"},
		{internal.LintReportFormatXml, "xml"},
		{internal.LintReportFormatHtml, "html"},
		{internal.LintReportFormatSarif, "sarif"},
	}

	for _, data := range dataSet {
		t.Run(string(data.format), func(t *testing.T) {
			test.InvokeLintFiles(
				t,
				phpMDLinter,
				internal.ProjectLinterOptions{
					ProjectInfo: internal.ProjectInfo{
						Directory: directory,
					},
					Output: internal.OutputOptions[internal.LintReportFormat]{
						Format: data.format,
					},
				},
				func(exec *test.Executable, executablePath string) {
					exec.OnRunOutput(executablePath, []string{"--version"}, phpMDVersion2).
						Return(nil)

					exec.OnRun(executablePath, []string{directory, data.arg, "phpmd.xml"}).
						Return(nil)
				},
				func(err error) {
					require.NoError(t, err)
				},
			)
		})
	}
}

func TestPHPMD_PHPMD_2_LintFiles_OutputFormat_Unsupported(t *testing.T) {
	directory := t.TempDir()

	formats := []internal.LintReportFormat{
		"test-unsupported",
		internal.LintReportFormatJUnit,
		internal.LintReportFormatTeamCity,
	}

	for _, format := range formats {
		t.Run(string(format), func(t *testing.T) {
			test.InvokeLintFiles(
				t,
				phpMDLinter,
				internal.ProjectLinterOptions{
					ProjectInfo: internal.ProjectInfo{
						Directory: directory,
					},
					Output: internal.OutputOptions[internal.LintReportFormat]{
						Format: format,
					},
				},
				func(exec *test.Executable, executablePath string) {
					exec.OnRunOutput(executablePath, []string{"--version"}, phpMDVersion2).
						Return(nil)
				},
				func(err error) {
					require.Error(t, err)
					require.EqualError(t, err, "unsupported report format: "+string(format))
				},
			)
		})
	}
}

func TestPHPMD_PHPMD_2_LintFiles_OutputFile(t *testing.T) {
	directory := t.TempDir()

	test.CreateFile(t, directory, "phpmd.xml")

	const content = "ok"

	outputFilename := t.TempDir() + "/report.txt"

	test.InvokeLintFiles(
		t,
		phpMDLinter,
		internal.ProjectLinterOptions{
			ProjectInfo: internal.ProjectInfo{
				Directory: directory,
			},
			Output: internal.OutputOptions[internal.LintReportFormat]{
				Filename: &outputFilename,
			},
		},
		func(exec *test.Executable, executablePath string) {
			exec.OnRunOutput(executablePath, []string{"--version"}, phpMDVersion2).
				Return(nil)

			exec.OnRunOutput(executablePath, []string{directory, "ansi", "phpmd.xml"}, content).
				Return(nil)
		},
		func(err error) {
			require.NoError(t, err)
			require.FileExists(t, outputFilename)

			outputContent, err := os.ReadFile(outputFilename)
			require.NoError(t, err)
			assert.Equal(t, content, string(outputContent))
		},
	)
}

func TestPHPMD_PHPMD_3_LintFiles_All_Success(t *testing.T) {
	directory := t.TempDir()

	test.InvokeLintFiles(
		t,
		phpMDLinter,
		internal.ProjectLinterOptions{
			ProjectInfo: internal.ProjectInfo{
				Directory: directory,
			},
		},
		func(exec *test.Executable, executablePath string) {
			exec.OnRunOutput(executablePath, []string{"--version"}, phpMDVersion3).
				Return(nil)

			exec.OnRun(executablePath, []string{"analyze"}).
				Return(nil)
		},
		func(err error) {
			require.NoError(t, err)
		},
	)
}

func TestPHPMD_PHPMD_3_LintFiles_All_Failed(t *testing.T) {
	expectedErr := errors.New("failed")

	test.InvokeLintFiles(
		t,
		phpMDLinter,
		internal.ProjectLinterOptions{
			ProjectInfo: internal.ProjectInfo{
				Directory: ".",
			},
		},
		func(exec *test.Executable, executablePath string) {
			exec.OnRunOutput(executablePath, []string{"--version"}, phpMDVersion3).
				Return(nil)

			exec.OnRun(executablePath, []string{"analyze"}).
				Return(expectedErr)
		},
		func(err error) {
			require.ErrorIs(t, err, expectedErr)
		},
	)
}

func TestPHPMD_PHPMD_3_LintFiles(t *testing.T) {
	fileNames := []string{"file.php", "/path/to/file2.php"}

	test.InvokeLintFiles(
		t,
		phpMDLinter,
		internal.ProjectLinterOptions{
			ProjectInfo: internal.ProjectInfo{
				Directory: ".",
			},
			Filenames: &fileNames,
		},
		func(exec *test.Executable, executablePath string) {
			exec.OnRunOutput(executablePath, []string{"--version"}, phpMDVersion3).
				Return(nil)

			exec.OnRun(executablePath, append([]string{"analyze"}, fileNames...)).
				Return(nil)
		},
		func(err error) {
			require.NoError(t, err)
		},
	)
}

func TestPHPMD_PHPMD_3_LintFiles_Failed(t *testing.T) {
	fileNames := []string{"file.php", "/path/to/file2.php"}

	expectedErr := errors.New("failed")

	test.InvokeLintFiles(
		t,
		phpMDLinter,
		internal.ProjectLinterOptions{
			ProjectInfo: internal.ProjectInfo{
				Directory: ".",
			},
			Filenames: &fileNames,
		},
		func(exec *test.Executable, executablePath string) {
			exec.OnRunOutput(executablePath, []string{"--version"}, phpMDVersion3).
				Return(nil)

			exec.OnRun(executablePath, append([]string{"analyze"}, fileNames...)).
				Return(expectedErr)
		},
		func(err error) {
			require.ErrorIs(t, err, expectedErr)
		},
	)
}

func TestPHPMD_PHPMD_3_LintFiles_OutputFormat(t *testing.T) {
	dataSet := []struct {
		format internal.LintReportFormat
		args   []string
	}{
		{internal.LintReportFormatRaw, []string{}},
		{internal.LintReportFormatJson, []string{"--format=json"}},
		{internal.LintReportFormatGitHub, []string{"--format=github"}},
		{internal.LintReportFormatGitLab, []string{"--format=gitlab"}},
		{internal.LintReportFormatXml, []string{"--format=xml"}},
		{internal.LintReportFormatHtml, []string{"--format=html"}},
		{internal.LintReportFormatSarif, []string{"--format=sarif"}},
	}

	for _, data := range dataSet {
		t.Run(string(data.format), func(t *testing.T) {
			test.InvokeLintFiles(
				t,
				phpMDLinter,
				internal.ProjectLinterOptions{
					ProjectInfo: internal.ProjectInfo{
						Directory: ".",
					},
					Output: internal.OutputOptions[internal.LintReportFormat]{
						Format: data.format,
					},
				},
				func(exec *test.Executable, executablePath string) {
					exec.OnRunOutput(executablePath, []string{"--version"}, phpMDVersion3).
						Return(nil)

					exec.OnRun(executablePath, append([]string{"analyze"}, data.args...)).
						Return(nil)
				},
				func(err error) {
					require.NoError(t, err)
				},
			)
		})
	}
}

func TestPHPMD_PHPMD_3_LintFiles_OutputFormat_Unsupported(t *testing.T) {
	formats := []internal.LintReportFormat{
		"test-unsupported",
		internal.LintReportFormatJUnit,
		internal.LintReportFormatTeamCity,
	}

	for _, format := range formats {
		t.Run(string(format), func(t *testing.T) {
			test.InvokeLintFiles(
				t,
				phpMDLinter,
				internal.ProjectLinterOptions{
					ProjectInfo: internal.ProjectInfo{
						Directory: ".",
					},
					Output: internal.OutputOptions[internal.LintReportFormat]{
						Format: format,
					},
				},
				func(exec *test.Executable, executablePath string) {
					exec.OnRunOutput(executablePath, []string{"--version"}, phpMDVersion3).
						Return(nil)
				},
				func(err error) {
					require.Error(t, err)
					require.EqualError(t, err, "unsupported report format: "+string(format))
				},
			)
		})
	}
}

func TestPHPMD_PHPMD_3_LintFiles_OutputFile(t *testing.T) {
	directory := t.TempDir()

	const content = "ok"

	outputFilename := directory + "/report.txt"

	test.InvokeLintFiles(
		t,
		phpMDLinter,
		internal.ProjectLinterOptions{
			ProjectInfo: internal.ProjectInfo{
				Directory: directory,
			},
			Output: internal.OutputOptions[internal.LintReportFormat]{
				Filename: &outputFilename,
			},
		},
		func(exec *test.Executable, executablePath string) {
			exec.OnRunOutput(executablePath, []string{"--version"}, phpMDVersion3).
				Return(nil)

			exec.OnRunOutput(executablePath, []string{"analyze"}, content).
				Return(nil)
		},
		func(err error) {
			require.NoError(t, err)
			require.FileExists(t, outputFilename)

			outputContent, err := os.ReadFile(outputFilename)
			require.NoError(t, err)
			assert.Equal(t, content, string(outputContent))
		},
	)
}

func phpMDDetect(ctx context.Context, options internal.DetectOptions) internal.Tool {
	return tool.DetectPHPMD(ctx, options)
}

func phpMDLinter(executable func() (*internal.Executable, error)) internal.ProjectLinter {
	return tool.NewPHPMD(executable)
}
