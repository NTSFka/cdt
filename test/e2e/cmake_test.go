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

	buildDir := t.TempDir()

	var err error

	// Configure
	err = runCdt(t.Context(), "-w", "cmake", "-r", "data/cmake", "-b", buildDir,
		"configure",
	)
	require.NoError(t, err)

	// Build all
	err = runCdt(t.Context(), "-w", "cmake", "-r", "data/cmake", "-b", buildDir,
		"build",
	)
	require.NoError(t, err)

	// Run
	buffer := bytes.Buffer{}
	err = runCdtOutput(t.Context(), io.MultiWriter(os.Stdout, &buffer), os.Stderr,
		"-w", "cmake", "-r", "data/cmake", "-b", buildDir,
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

	buildDir := t.TempDir()

	buffer := bytes.Buffer{}
	err := runCdtOutput(t.Context(), io.MultiWriter(os.Stdout, &buffer), os.Stderr,
		"-w", "cmake", "-r", "data/cmake", "-b", buildDir,
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

	buildDir := t.TempDir()

	buffer := bytes.Buffer{}
	err := runCdtOutput(t.Context(), io.MultiWriter(os.Stdout, &buffer), os.Stderr,
		"-w", "cmake", "-r", "data/cmake", "-b", buildDir,
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

	buildDir := t.TempDir()

	buffer := bytes.Buffer{}
	err := runCdtOutput(t.Context(), io.MultiWriter(os.Stdout, &buffer), os.Stderr,
		"-w", "cmake", "-r", "data/cmake", "-b", buildDir,
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

	buildDir := t.TempDir()

	buffer := bytes.Buffer{}
	err := runCdtOutput(t.Context(), io.MultiWriter(os.Stdout, &buffer), os.Stderr,
		"-w", "cmake", "-r", "data/cmake", "-b", buildDir,
		"lint",
	)
	require.NoError(t, err)
}
