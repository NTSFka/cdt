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

func TestGoBuild(t *testing.T) {
	environment := internal.SystemEnvironment

	checkTool(t, t.Context(), environment, "go")

	outputDir := t.TempDir()

	buffer := bytes.Buffer{}
	err := runCdtOutput(t.Context(), io.MultiWriter(os.Stdout, &buffer), os.Stderr,
		"-w", "go", "-r", "data/go", "-o", outputDir,
		"build",
	)
	require.NoError(t, err)
}

func TestGoRun(t *testing.T) {
	environment := internal.SystemEnvironment

	checkTool(t, t.Context(), environment, "go")

	outputDir := t.TempDir()

	var err error

	// Run
	buffer := bytes.Buffer{}
	err = runCdtOutput(t.Context(), io.MultiWriter(os.Stdout, &buffer), os.Stderr,
		"-w", "go", "-r", "data/go", "-o", outputDir,
		"run", "hello",
	)
	require.NoError(t, err)
	require.Contains(t, buffer.String(), "Hello World!\n")
}

func TestGoTest(t *testing.T) {
	environment := internal.SystemEnvironment

	checkTool(t, t.Context(), environment, "go")

	outputDir := t.TempDir()

	buffer := bytes.Buffer{}
	err := runCdtOutput(t.Context(), io.MultiWriter(os.Stdout, &buffer), os.Stderr,
		"-w", "go", "-r", "data/go", "-o", outputDir,
		"test",
	)
	require.NoError(t, err)
	require.Contains(t, buffer.String(), "ok")
	require.NotContains(t, buffer.String(), "FAIL")
}

func TestGoFormat(t *testing.T) {
	environment := internal.SystemEnvironment

	checkTool(t, t.Context(), environment, "go")

	outputDir := t.TempDir()

	buffer := bytes.Buffer{}
	err := runCdtOutput(t.Context(), io.MultiWriter(os.Stdout, &buffer), os.Stderr,
		"-w", "go", "-r", "data/go", "-o", outputDir,
		"format",
	)
	require.NoError(t, err)
}

func TestGoLint(t *testing.T) {
	environment := internal.SystemEnvironment

	checkTool(t, t.Context(), environment, "go")

	outputDir := t.TempDir()

	buffer := bytes.Buffer{}
	err := runCdtOutput(t.Context(), io.MultiWriter(os.Stdout, &buffer), os.Stderr,
		"-w", "go", "-r", "data/go", "-o", outputDir,
		"lint",
	)
	require.NoError(t, err)
}

func TestGoBuild_Docker(t *testing.T) {
	if runtime.GOOS == "windows" { // nolint: goconst
		t.Skip("docker issues on Windows")
	}

	environment := internal.SystemEnvironment

	checkTool(t, t.Context(), environment, "docker")

	outputDir := t.TempDir()

	buffer := bytes.Buffer{}
	err := runCdtOutput(t.Context(), io.MultiWriter(os.Stdout, &buffer), os.Stderr,
		"-e", "d:golang:1.25", "-w", "go", "-r", "data/go", "-o", outputDir,
		"build",
	)
	require.NoError(t, err)
}

func TestGoRun_Docker(t *testing.T) {
	if runtime.GOOS == "windows" { // nolint: goconst
		t.Skip("docker issues on Windows")
	}

	environment := internal.SystemEnvironment

	checkTool(t, t.Context(), environment, "docker")

	outputDir := t.TempDir()

	var err error

	// Run
	buffer := bytes.Buffer{}
	err = runCdtOutput(t.Context(), io.MultiWriter(os.Stdout, &buffer), os.Stderr,
		"-e", "d:golang:1.25", "-w", "go", "-r", "data/go", "-o", outputDir,
		"run", "hello",
	)
	require.NoError(t, err)
	require.Contains(t, buffer.String(), "Hello World!\n")
}

func TestGoTest_Docker(t *testing.T) {
	if runtime.GOOS == "windows" { // nolint: goconst
		t.Skip("docker issues on Windows")
	}

	environment := internal.SystemEnvironment

	checkTool(t, t.Context(), environment, "docker")

	outputDir := t.TempDir()

	buffer := bytes.Buffer{}
	err := runCdtOutput(t.Context(), io.MultiWriter(os.Stdout, &buffer), os.Stderr,
		"-e", "d:golang:1.25", "-w", "go", "-r", "data/go", "-o", outputDir,
		"test",
	)
	require.NoError(t, err)
	require.Contains(t, buffer.String(), "ok")
	require.NotContains(t, buffer.String(), "FAIL")
}

func TestGoFormat_Docker(t *testing.T) {
	if runtime.GOOS == "windows" { // nolint: goconst
		t.Skip("docker issues on Windows")
	}

	environment := internal.SystemEnvironment

	checkTool(t, t.Context(), environment, "docker")

	outputDir := t.TempDir()

	buffer := bytes.Buffer{}
	err := runCdtOutput(t.Context(), io.MultiWriter(os.Stdout, &buffer), os.Stderr,
		"-e", "d:golang:1.25", "-w", "go", "-r", "data/go", "-o", outputDir,
		"format",
	)
	require.NoError(t, err)
}

func TestGoLint_Docker(t *testing.T) {
	if runtime.GOOS == "windows" { // nolint: goconst
		t.Skip("docker issues on Windows")
	}

	environment := internal.SystemEnvironment

	checkTool(t, t.Context(), environment, "docker")

	outputDir := t.TempDir()

	buffer := bytes.Buffer{}
	err := runCdtOutput(t.Context(), io.MultiWriter(os.Stdout, &buffer), os.Stderr,
		"-e", "d:golang:1.25", "-w", "go", "-r", "data/go", "-o", outputDir,
		"lint",
	)
	require.NoError(t, err)
}
