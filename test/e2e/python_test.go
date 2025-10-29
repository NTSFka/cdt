package e2e_test

import (
	"bytes"
	"context"
	"io"
	"os"
	"os/exec"
	"runtime"
	"testing"

	"cdt/internal"

	"github.com/stretchr/testify/require"
)

func initPythonEnvironment(ctx context.Context, directory string) error {
	return exec.CommandContext(ctx, "python3", "-m", "venv", directory).Run()
}

func TestPythonRun(t *testing.T) {
	// FIXME: Python detection on Windows in Github Actions
	if runtime.GOOS == "windows" { // nolint: goconst
		t.Skip("python3 detection is not working on Windows")
	}

	environment := internal.SystemEnvironment

	checkTool(t, t.Context(), environment, "python3")

	envDir := t.TempDir()
	require.NoError(t, initPythonEnvironment(t.Context(), envDir))

	outputDir := t.TempDir()

	var err error

	buffer := bytes.Buffer{}
	err = runCdtOutput(t.Context(), io.MultiWriter(os.Stdout, &buffer), os.Stderr,
		"-e", "pyenv:"+envDir,
		"-w", "python", "-r", "data/python", "-o", outputDir,
		"run", "hello.py",
	)
	require.NoError(t, err)
	require.Contains(t, buffer.String(), "Hello World!\n")
}

func TestPythonTest(t *testing.T) {
	// FIXME: Python detection on Windows in Github Actions
	if runtime.GOOS == "windows" { // nolint: goconst
		t.Skip("python3 detection is not working on Windows")
	}

	environment := internal.SystemEnvironment

	checkTool(t, t.Context(), environment, "python3")

	envDir := t.TempDir()
	require.NoError(t, initPythonEnvironment(t.Context(), envDir))

	outputDir := t.TempDir()

	// Install dependencies
	var err error
	err = runCdt(t.Context(),
		"-e", "pyenv:"+envDir,
		"-w", "python", "-r", "data/python", "-o", outputDir,
		"dep", "add", "pytest",
	)
	require.NoError(t, err)

	buffer := bytes.Buffer{}
	err = runCdtOutput(t.Context(), io.MultiWriter(os.Stdout, &buffer), os.Stderr,
		"-e", "pyenv:"+envDir,
		"-w", "python", "-r", "data/python", "-o", outputDir,
		"test",
	)
	require.NoError(t, err)
	require.Contains(t, buffer.String(), "passed")
}

func TestPythonFormat(t *testing.T) {
	// FIXME: Python detection on Windows in Github Actions
	if runtime.GOOS == "windows" { // nolint: goconst
		t.Skip("python3 detection is not working on Windows")
	}

	environment := internal.SystemEnvironment

	checkTool(t, t.Context(), environment, "python3")

	envDir := t.TempDir()
	require.NoError(t, initPythonEnvironment(t.Context(), envDir))

	outputDir := t.TempDir()

	// Install dependencies
	var err error
	err = runCdt(t.Context(),
		"-e", "pyenv:"+envDir,
		"-w", "python", "-r", "data/python", "-o", outputDir,
		"dep", "add", "black",
	)
	require.NoError(t, err)

	buffer := bytes.Buffer{}
	err = runCdtOutput(t.Context(), io.MultiWriter(os.Stdout, &buffer), os.Stderr,
		"-e", "pyenv:"+envDir,
		"-w", "python", "-r", "data/python", "-o", outputDir,
		"format",
	)
	require.NoError(t, err)
}

func TestPythonLint(t *testing.T) {
	// FIXME: Python detection on Windows in Github Actions
	if runtime.GOOS == "windows" { // nolint: goconst
		t.Skip("python3 detection is not working on Windows")
	}

	environment := internal.SystemEnvironment

	checkTool(t, t.Context(), environment, "python3")

	envDir := t.TempDir()
	require.NoError(t, initPythonEnvironment(t.Context(), envDir))

	outputDir := t.TempDir()

	// Install dependencies
	var err error
	err = runCdt(t.Context(),
		"-e", "pyenv:"+envDir,
		"-w", "python", "-r", "data/python", "-o", outputDir,
		"dep", "add", "pylint",
	)
	require.NoError(t, err)

	buffer := bytes.Buffer{}
	err = runCdtOutput(t.Context(), io.MultiWriter(os.Stdout, &buffer), os.Stderr,
		"-e", "pyenv:"+envDir,
		"-w", "python", "-r", "data/python", "-o", outputDir,
		"lint",
	)
	require.NoError(t, err)
}

func TestPythonRun_Docker(t *testing.T) {
	if runtime.GOOS == "windows" { // nolint: goconst
		t.Skip("docker issues on Windows")
	}

	environment := internal.SystemEnvironment

	checkTool(t, t.Context(), environment, "docker")

	outputDir := t.TempDir()

	var err error

	buffer := bytes.Buffer{}
	err = runCdtOutput(t.Context(), io.MultiWriter(os.Stdout, &buffer), os.Stderr,
		"-e", "d:python:3.14",
		"-w", "python", "-r", "data/python", "-o", outputDir,
		"run", "hello.py",
	)
	require.NoError(t, err)
	require.Contains(t, buffer.String(), "Hello World!\n")
}

func TestPythonTest_Docker(t *testing.T) {
	t.Skip("cannot share container between invocations")

	environment := internal.SystemEnvironment

	checkTool(t, t.Context(), environment, "docker")

	outputDir := t.TempDir()

	// Install dependencies
	var err error
	err = runCdt(t.Context(),
		"-e", "d:python:3.14",
		"-w", "python", "-r", "data/python", "-o", outputDir,
		"dep", "add", "pytest",
	)
	require.NoError(t, err)

	buffer := bytes.Buffer{}
	err = runCdtOutput(t.Context(), io.MultiWriter(os.Stdout, &buffer), os.Stderr,
		"-e", "d:python:3.14",
		"-w", "python", "-r", "data/python", "-o", outputDir,
		"test",
	)
	require.NoError(t, err)
	require.Contains(t, buffer.String(), "passed")
}

func TestPythonFormat_Docker(t *testing.T) {
	t.Skip("cannot share container between invocations")

	environment := internal.SystemEnvironment

	checkTool(t, t.Context(), environment, "docker")

	envDir := t.TempDir()
	require.NoError(t, initPythonEnvironment(t.Context(), envDir))

	outputDir := t.TempDir()

	// Install dependencies
	var err error
	err = runCdt(t.Context(),
		"-e", "pyenv:"+envDir,
		"-w", "python", "-r", "data/python", "-o", outputDir,
		"dep", "add", "black",
	)
	require.NoError(t, err)

	buffer := bytes.Buffer{}
	err = runCdtOutput(t.Context(), io.MultiWriter(os.Stdout, &buffer), os.Stderr,
		"-e", "pyenv:"+envDir,
		"-w", "python", "-r", "data/python", "-o", outputDir,
		"format",
	)
	require.NoError(t, err)
}

func TestPythonLint_Docker(t *testing.T) {
	t.Skip("cannot share container between invocations")

	environment := internal.SystemEnvironment

	checkTool(t, t.Context(), environment, "docker")

	envDir := t.TempDir()
	require.NoError(t, initPythonEnvironment(t.Context(), envDir))

	outputDir := t.TempDir()

	// Install dependencies
	var err error
	err = runCdt(t.Context(),
		"-e", "pyenv:"+envDir,
		"-w", "python", "-r", "data/python", "-o", outputDir,
		"dep", "add", "pylint",
	)
	require.NoError(t, err)

	buffer := bytes.Buffer{}
	err = runCdtOutput(t.Context(), io.MultiWriter(os.Stdout, &buffer), os.Stderr,
		"-e", "pyenv:"+envDir,
		"-w", "python", "-r", "data/python", "-o", outputDir,
		"lint",
	)
	require.NoError(t, err)
}
