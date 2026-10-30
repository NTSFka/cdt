package tool_test

import (
	"bytes"
	"cdt/internal/reporter"
	"cdt/internal/test"
	"context"
	"errors"
	"io"
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

	test.InvokeAnalyzeFiles(
		t,
		phpMDAnalyzer,
		internal.ProjectAnalyzerOptions{
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
	test.InvokeAnalyzeFiles(
		t,
		phpMDAnalyzer,
		internal.ProjectAnalyzerOptions{
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

func TestPHPMD_PHPMD_2_AnalyzeFiles_All_Success(t *testing.T) {
	dataSet := []string{
		"phpmd.xml",
		"phpmd.ruleset.xml",
	}

	for _, fileName := range dataSet {
		t.Run(fileName, func(t *testing.T) {
			directory := t.TempDir()

			test.CreateFile(t, directory, fileName)

			test.InvokeAnalyzeFiles(
				t,
				phpMDAnalyzer,
				internal.ProjectAnalyzerOptions{
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

func TestPHPMD_PHPMD_2_AnalyzeFiles_All_Directories(t *testing.T) {
	directory := t.TempDir()

	test.CreateFile(t, directory, "phpmd.xml")
	test.CreateDirectory(t, directory, "src")
	test.CreateDirectory(t, directory, "tests")

	test.InvokeAnalyzeFiles(
		t,
		phpMDAnalyzer,
		internal.ProjectAnalyzerOptions{
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

func TestPHPMD_PHPMD_2_AnalyzeFiles_All_Failed(t *testing.T) {
	directory := t.TempDir()

	test.CreateFile(t, directory, "phpmd.xml")

	expectedErr := errors.New("failed")

	test.InvokeAnalyzeFiles(
		t,
		phpMDAnalyzer,
		internal.ProjectAnalyzerOptions{
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

func TestPHPMD_PHPMD_2_AnalyzeFiles_All_MissingRulesFile(t *testing.T) {
	test.InvokeAnalyzeFiles(
		t,
		phpMDAnalyzer,
		internal.ProjectAnalyzerOptions{},
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

func TestPHPMD_PHPMD_2_AnalyzeFiles_All_CustomRulesFile(t *testing.T) {
	directory := t.TempDir()

	test.InvokeAnalyzeFiles(
		t,
		phpMDAnalyzer,
		internal.ProjectAnalyzerOptions{
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

func TestPHPMD_PHPMD_2_AnalyzeFiles(t *testing.T) {
	directory := t.TempDir()

	test.CreateFile(t, directory, "phpmd.xml")

	fileNames := []string{"file.php", "/path/to/file2.php"}

	test.InvokeAnalyzeFiles(
		t,
		phpMDAnalyzer,
		internal.ProjectAnalyzerOptions{
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

func TestPHPMD_PHPMD_2_AnalyzeFiles_Failed(t *testing.T) {
	directory := t.TempDir()

	test.CreateFile(t, directory, "phpmd.xml")

	fileNames := []string{"file.php", "/path/to/file2.php"}

	expectedErr := errors.New("failed")

	test.InvokeAnalyzeFiles(
		t,
		phpMDAnalyzer,
		internal.ProjectAnalyzerOptions{
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

func TestPHPMD_2_AnalyzeFiles_Reporter_NoErrors(t *testing.T) {
	directory := t.TempDir()

	result := bytes.Buffer{}

	analyzeReporter := reporter.NewSimpleSarifReporter(&result)

	test.InvokeAnalyzeFiles(
		t,
		phpMDAnalyzer,
		internal.ProjectAnalyzerOptions{
			ProjectInfo: internal.ProjectInfo{
				Directory: directory,
			},
			Output: internal.AnalyzeOutputOptions{
				Reporter: analyzeReporter,
				Output:   io.Discard,
			},
			ExtraArgs: []string{"phpmd.rules.xml"},
		},
		func(exec *test.Executable, executablePath string) {
			exec.OnRunOutput(executablePath, []string{"--version"}, phpMDVersion2).
				Return(nil)

			exec.OnRunOutput(executablePath, []string{directory, "sarif", "phpmd.rules.xml"}, "{}").
				Return(nil)
		},
		func(err error) {
			require.NoError(t, err)
		},
	)

	assert.Empty(t, result.String())
}

// nolint:funlen
func TestPHPMD_2_AnalyzeFiles_Reporter_SomeErrors(t *testing.T) {
	directory := t.TempDir()

	result := bytes.Buffer{}

	analyzeReporter := reporter.NewSimpleSarifReporter(&result)

	test.InvokeAnalyzeFiles(
		t,
		phpMDAnalyzer,
		internal.ProjectAnalyzerOptions{
			ProjectInfo: internal.ProjectInfo{
				Directory: directory,
			},
			Output: internal.AnalyzeOutputOptions{
				Reporter: analyzeReporter,
				Output:   io.Discard,
			},
			ExtraArgs: []string{"phpmd.rules.xml"},
		},
		func(executable *test.Executable, executablePath string) {
			executable.OnRunOutput(executablePath, []string{"--version"}, phpMDVersion2).
				Return(nil)

			executable.OnRunOutput(executablePath, []string{directory, "sarif", "phpmd.rules.xml"}, `{
    "version": "2.1.0",
    "$schema": "https://raw.githubusercontent.com/oasis-tcs/sarif-spec/master/Schemata/sarif-schema-2.1.0.json",
    "runs": [
      {
        "tool": {
          "driver": {
            "informationUri": "https://github.com/carthage-software/phpmd",
            "name": "phpmd",
            "version": "1.57.0"
          }
        },
        "artifacts": [
          {
            "location": {
              "uri": "file:///test.php"
            }
          }
        ],
        "results": [
          {
            "level": "error",
            "message": {
              "text": "'x' is assigned a value but never used."
            },
            "locations": [
              {
                "physicalLocation": {
                  "artifactLocation": {
                    "uri": "file:///test.php",
                    "index": 0
                  },
                  "region": {
                    "startLine": 1,
                    "startColumn": 5
                  }
                }
              }
            ],
            "ruleId": "no-unused-vars",
            "ruleIndex": 0
          }
        ]
      }
    ]
  }`).
				Return(nil)
		},
		func(err error) {
			require.NoError(t, err)
		},
	)

	require.NoError(t, analyzeReporter.Finish())

	assert.NotEmpty(t, result.String())
}

func TestPHPMD_2_AnalyzeFiles_Reporter_NoOutput(t *testing.T) {
	directory := t.TempDir()

	result := bytes.Buffer{}

	analyzeReporter := reporter.NewSimpleSarifReporter(&result)

	test.InvokeAnalyzeFiles(
		t,
		phpMDAnalyzer,
		internal.ProjectAnalyzerOptions{
			ProjectInfo: internal.ProjectInfo{
				Directory: directory,
			},
			Output: internal.AnalyzeOutputOptions{
				Reporter: analyzeReporter,
				Output:   io.Discard,
			},
			ExtraArgs: []string{"phpmd.rules.xml"},
		},
		func(executable *test.Executable, executablePath string) {
			executable.OnRunOutput(executablePath, []string{"--version"}, phpMDVersion2).
				Return(nil)

			executable.OnRunOutput(executablePath, []string{directory, "sarif", "phpmd.rules.xml"}, "").
				Return(nil)
		},
		func(err error) {
			require.NoError(t, err)
		},
	)

	assert.Empty(t, result.String())
}

func TestPHPMD_2_AnalyzeFiles_Reporter_InvalidJson(t *testing.T) {
	directory := t.TempDir()

	result := bytes.Buffer{}

	analyzeReporter := reporter.NewSimpleSarifReporter(&result)

	test.InvokeAnalyzeFiles(
		t,
		phpMDAnalyzer,
		internal.ProjectAnalyzerOptions{
			ProjectInfo: internal.ProjectInfo{
				Directory: directory,
			},
			Output: internal.AnalyzeOutputOptions{
				Reporter: analyzeReporter,
				Output:   io.Discard,
			},
			ExtraArgs: []string{"phpmd.rules.xml"},
		},
		func(executable *test.Executable, executablePath string) {
			executable.OnRunOutput(executablePath, []string{"--version"}, phpMDVersion2).
				Return(nil)

			executable.OnRunOutput(executablePath, []string{directory, "sarif", "phpmd.rules.xml"}, ".").
				Return(nil)
		},
		func(err error) {
			require.ErrorContains(t, err, "failed to parse SARIF output:")
		},
	)

	assert.Empty(t, result.String())
}

func TestPHPMD_PHPMD_2_AnalyzeFiles_OutputFile(t *testing.T) {
	directory := t.TempDir()

	test.CreateFile(t, directory, "phpmd.xml")

	const content = "ok"

	outputFilename := t.TempDir() + "/report.txt"

	file, err := os.OpenFile(outputFilename, os.O_CREATE|os.O_WRONLY, 0644)
	require.NoError(t, err)

	defer func() {
		_ = file.Close()
	}()

	test.InvokeAnalyzeFiles(
		t,
		phpMDAnalyzer,
		internal.ProjectAnalyzerOptions{
			ProjectInfo: internal.ProjectInfo{
				Directory: directory,
			},
			Output: internal.AnalyzeOutputOptions{
				Output: file,
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

func TestPHPMD_PHPMD_3_AnalyzeFiles_All_Success(t *testing.T) {
	directory := t.TempDir()

	test.InvokeAnalyzeFiles(
		t,
		phpMDAnalyzer,
		internal.ProjectAnalyzerOptions{
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

func TestPHPMD_PHPMD_3_AnalyzeFiles_All_Failed(t *testing.T) {
	expectedErr := errors.New("failed")

	test.InvokeAnalyzeFiles(
		t,
		phpMDAnalyzer,
		internal.ProjectAnalyzerOptions{
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

func TestPHPMD_PHPMD_3_AnalyzeFiles(t *testing.T) {
	fileNames := []string{"file.php", "/path/to/file2.php"}

	test.InvokeAnalyzeFiles(
		t,
		phpMDAnalyzer,
		internal.ProjectAnalyzerOptions{
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

func TestPHPMD_PHPMD_3_AnalyzeFiles_Failed(t *testing.T) {
	fileNames := []string{"file.php", "/path/to/file2.php"}

	expectedErr := errors.New("failed")

	test.InvokeAnalyzeFiles(
		t,
		phpMDAnalyzer,
		internal.ProjectAnalyzerOptions{
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

func TestPHPMD_3_AnalyzeFiles_Reporter_NoErrors(t *testing.T) {
	result := bytes.Buffer{}

	analyzeReporter := reporter.NewSimpleSarifReporter(&result)

	test.InvokeAnalyzeFiles(
		t,
		phpMDAnalyzer,
		internal.ProjectAnalyzerOptions{
			ProjectInfo: internal.ProjectInfo{
				Directory: ".",
			},
			Output: internal.AnalyzeOutputOptions{
				Reporter: analyzeReporter,
				Output:   io.Discard,
			},
		},
		func(exec *test.Executable, executablePath string) {
			exec.OnRunOutput(executablePath, []string{"--version"}, phpMDVersion3).
				Return(nil)

			exec.OnRunOutput(executablePath, []string{"analyze", "--format=sarif"}, "{}").
				Return(nil)
		},
		func(err error) {
			require.NoError(t, err)
		},
	)

	assert.Empty(t, result.String())
}

// nolint:funlen
func TestPHPMD_3_AnalyzeFiles_Reporter_SomeErrors(t *testing.T) {
	result := bytes.Buffer{}

	analyzeReporter := reporter.NewSimpleSarifReporter(&result)

	test.InvokeAnalyzeFiles(
		t,
		phpMDAnalyzer,
		internal.ProjectAnalyzerOptions{
			ProjectInfo: internal.ProjectInfo{
				Directory: ".",
			},
			Output: internal.AnalyzeOutputOptions{
				Reporter: analyzeReporter,
				Output:   io.Discard,
			},
		},
		func(executable *test.Executable, executablePath string) {
			executable.OnRunOutput(executablePath, []string{"--version"}, phpMDVersion3).
				Return(nil)

			executable.OnRunOutput(executablePath, []string{"analyze", "--format=sarif"}, `{
    "version": "2.1.0",
    "$schema": "https://raw.githubusercontent.com/oasis-tcs/sarif-spec/master/Schemata/sarif-schema-2.1.0.json",
    "runs": [
      {
        "tool": {
          "driver": {
            "informationUri": "https://github.com/carthage-software/phpmd",
            "name": "phpmd",
            "version": "1.57.0"
          }
        },
        "artifacts": [
          {
            "location": {
              "uri": "file:///test.php"
            }
          }
        ],
        "results": [
          {
            "level": "error",
            "message": {
              "text": "'x' is assigned a value but never used."
            },
            "locations": [
              {
                "physicalLocation": {
                  "artifactLocation": {
                    "uri": "file:///test.php",
                    "index": 0
                  },
                  "region": {
                    "startLine": 1,
                    "startColumn": 5
                  }
                }
              }
            ],
            "ruleId": "no-unused-vars",
            "ruleIndex": 0
          }
        ]
      }
    ]
  }`).
				Return(nil)
		},
		func(err error) {
			require.NoError(t, err)
		},
	)

	require.NoError(t, analyzeReporter.Finish())

	assert.NotEmpty(t, result.String())
}

func TestPHPMD_3_AnalyzeFiles_Reporter_NoOutput(t *testing.T) {
	result := bytes.Buffer{}

	analyzeReporter := reporter.NewSimpleSarifReporter(&result)

	test.InvokeAnalyzeFiles(
		t,
		phpMDAnalyzer,
		internal.ProjectAnalyzerOptions{
			ProjectInfo: internal.ProjectInfo{
				Directory: ".",
			},
			Output: internal.AnalyzeOutputOptions{
				Reporter: analyzeReporter,
				Output:   io.Discard,
			},
		},
		func(executable *test.Executable, executablePath string) {
			executable.OnRunOutput(executablePath, []string{"--version"}, phpMDVersion3).
				Return(nil)

			executable.OnRunOutput(executablePath, []string{"analyze", "--format=sarif"}, "").
				Return(nil)
		},
		func(err error) {
			require.NoError(t, err)
		},
	)

	assert.Empty(t, result.String())
}

func TestPHPMD_3_AnalyzeFiles_Reporter_InvalidJson(t *testing.T) {
	result := bytes.Buffer{}

	analyzeReporter := reporter.NewSimpleSarifReporter(&result)

	test.InvokeAnalyzeFiles(
		t,
		phpMDAnalyzer,
		internal.ProjectAnalyzerOptions{
			ProjectInfo: internal.ProjectInfo{
				Directory: ".",
			},
			Output: internal.AnalyzeOutputOptions{
				Reporter: analyzeReporter,
				Output:   io.Discard,
			},
		},
		func(executable *test.Executable, executablePath string) {
			executable.OnRunOutput(executablePath, []string{"--version"}, phpMDVersion3).
				Return(nil)

			executable.OnRunOutput(executablePath, []string{"analyze", "--format=sarif"}, ".").
				Return(nil)
		},
		func(err error) {
			require.ErrorContains(t, err, "failed to parse SARIF output:")
		},
	)

	assert.Empty(t, result.String())
}

func TestPHPMD_PHPMD_3_AnalyzeFiles_OutputFile(t *testing.T) {
	directory := t.TempDir()

	const content = "ok"

	outputFilename := directory + "/report.txt"

	file, err := os.OpenFile(outputFilename, os.O_CREATE|os.O_WRONLY, 0644)
	require.NoError(t, err)

	defer func() {
		_ = file.Close()
	}()

	test.InvokeAnalyzeFiles(
		t,
		phpMDAnalyzer,
		internal.ProjectAnalyzerOptions{
			ProjectInfo: internal.ProjectInfo{
				Directory: directory,
			},
			Output: internal.AnalyzeOutputOptions{
				Output: file,
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

func phpMDAnalyzer(executable func() (*internal.Executable, error)) internal.ProjectAnalyzer {
	return tool.NewPHPMD(executable)
}
