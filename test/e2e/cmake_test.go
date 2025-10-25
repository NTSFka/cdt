package e2e_test

import (
	"bytes"
	"io"
	"os"
	"runtime"
	"testing"

	"cdt/internal"

	"github.com/stretchr/testify/require"
)

func TestCMakeConfigureAndBuildAndRun(t *testing.T) {
	if runtime.GOOS == "windows" { // nolint: goconst
		t.Skip("project structure is not detected properly on Windows using MSVC Generator")
	}

	environment := internal.SystemEnvironment

	checkTool(t, t.Context(), environment, "cmake")

	outputDir := t.TempDir()

	var err error

	// Configure
	err = runCdt(t.Context(), "-w", "cmake", "-r", "data/cmake", "-o", outputDir,
		"configure",
	)
	require.NoError(t, err)

	// Build all
	err = runCdt(t.Context(), "-w", "cmake", "-r", "data/cmake", "-o", outputDir,
		"build",
	)
	require.NoError(t, err)

	// Run
	buffer := bytes.Buffer{}
	err = runCdtOutput(t.Context(), io.MultiWriter(os.Stdout, &buffer), os.Stderr,
		"-w", "cmake", "-r", "data/cmake", "-o", outputDir,
		"run", "main",
	)
	require.NoError(t, err)
	require.Contains(t, buffer.String(), "Hello World!\n")
}

func TestCMakeBuild(t *testing.T) {
	if runtime.GOOS == "windows" { // nolint: goconst
		t.Skip("project structure is not detected properly on Windows using MSVC Generator")
	}

	environment := internal.SystemEnvironment

	checkTool(t, t.Context(), environment, "cmake")

	outputDir := t.TempDir()

	buffer := bytes.Buffer{}
	err := runCdtOutput(t.Context(), io.MultiWriter(os.Stdout, &buffer), os.Stderr,
		"-w", "cmake", "-r", "data/cmake", "-o", outputDir,
		"build",
	)
	require.NoError(t, err)
}

func TestCMakeTest(t *testing.T) {
	if runtime.GOOS == "windows" { // nolint: goconst
		t.Skip("project structure is not detected properly on Windows using MSVC Generator")
	}

	environment := internal.SystemEnvironment

	checkTool(t, t.Context(), environment, "cmake")

	outputDir := t.TempDir()

	buffer := bytes.Buffer{}
	err := runCdtOutput(t.Context(), io.MultiWriter(os.Stdout, &buffer), os.Stderr,
		"-w", "cmake", "-r", "data/cmake", "-o", outputDir,
		"test",
	)
	require.NoError(t, err)
	require.Contains(t, buffer.String(), "Passed")
}

func TestCMakeFormat(t *testing.T) {
	if runtime.GOOS == "windows" { // nolint: goconst
		t.Skip("project structure is not detected properly on Windows using MSVC Generator")
	}

	environment := internal.SystemEnvironment

	checkTool(t, t.Context(), environment, "cmake")
	checkTool(t, t.Context(), environment, "clang-format")

	outputDir := t.TempDir()

	buffer := bytes.Buffer{}
	err := runCdtOutput(t.Context(), io.MultiWriter(os.Stdout, &buffer), os.Stderr,
		"-w", "cmake", "-r", "data/cmake", "-o", outputDir,
		"format",
	)
	require.NoError(t, err)
}

func TestCMakeLint(t *testing.T) {
	if runtime.GOOS == "windows" { // nolint: goconst
		t.Skip("project structure is not detected properly on Windows using MSVC Generator")
	}

	environment := internal.SystemEnvironment

	checkTool(t, t.Context(), environment, "cmake")
	checkTool(t, t.Context(), environment, "clang-tidy")

	outputDir := t.TempDir()

	buffer := bytes.Buffer{}
	err := runCdtOutput(t.Context(), io.MultiWriter(os.Stdout, &buffer), os.Stderr,
		"-w", "cmake", "-r", "data/cmake", "-o", outputDir,
		"lint",
	)
	require.NoError(t, err)
}
