package e2e_test

import (
	"bytes"
	"context"
	"io"
	"os"
	"os/exec"
	"testing"

	"cdt/internal"

	"github.com/stretchr/testify/require"
)

func initPythonEnvironment(ctx context.Context, directory string) error {
	return exec.CommandContext(ctx, "python3", "-m", "venv", directory).Run()
}

func TestPythonRun(t *testing.T) {
	environment := internal.SystemEnvironment

	checkTool(t, t.Context(), environment, "python3")

	envDir := t.TempDir()
	require.NoError(t, initPythonEnvironment(t.Context(), envDir))

	buildDir := t.TempDir()

	var err error

	buffer := bytes.Buffer{}
	err = runCdtOutput(t.Context(), io.MultiWriter(os.Stdout, &buffer), os.Stderr,
		"-e", "pyenv:"+envDir,
		"-w", "python", "-r", "data/python", "-b", buildDir,
		"run", "hello.py",
	)
	require.NoError(t, err)
	require.Contains(t, buffer.String(), "Hello World!\n")
}

func TestPythonTest(t *testing.T) {
	environment := internal.SystemEnvironment

	checkTool(t, t.Context(), environment, "python3")

	envDir := t.TempDir()
	require.NoError(t, initPythonEnvironment(t.Context(), envDir))

	buildDir := t.TempDir()

	// Install dependencies
	var err error
	err = runCdt(t.Context(),
		"-e", "pyenv:"+envDir,
		"-w", "python", "-r", "data/python", "-b", buildDir,
		"dep", "add", "pytest",
	)
	require.NoError(t, err)

	buffer := bytes.Buffer{}
	err = runCdtOutput(t.Context(), io.MultiWriter(os.Stdout, &buffer), os.Stderr,
		"-e", "pyenv:"+envDir,
		"-w", "python", "-r", "data/python", "-b", buildDir,
		"test",
	)
	require.NoError(t, err)
	require.Contains(t, buffer.String(), "passed")
}

func TestPythonFormat(t *testing.T) {
	environment := internal.SystemEnvironment

	checkTool(t, t.Context(), environment, "python3")

	envDir := t.TempDir()
	require.NoError(t, initPythonEnvironment(t.Context(), envDir))

	buildDir := t.TempDir()

	// Install dependencies
	var err error
	err = runCdt(t.Context(),
		"-e", "pyenv:"+envDir,
		"-w", "python", "-r", "data/python", "-b", buildDir,
		"dep", "add", "black",
	)
	require.NoError(t, err)

	buffer := bytes.Buffer{}
	err = runCdtOutput(t.Context(), io.MultiWriter(os.Stdout, &buffer), os.Stderr,
		"-e", "pyenv:"+envDir,
		"-w", "python", "-r", "data/python", "-b", buildDir,
		"format",
	)
	require.NoError(t, err)
}

func TestPythonLint(t *testing.T) {
	environment := internal.SystemEnvironment

	checkTool(t, t.Context(), environment, "python3")

	envDir := t.TempDir()
	require.NoError(t, initPythonEnvironment(t.Context(), envDir))

	buildDir := t.TempDir()

	// Install dependencies
	var err error
	err = runCdt(t.Context(),
		"-e", "pyenv:"+envDir,
		"-w", "python", "-r", "data/python", "-b", buildDir,
		"dep", "add", "pylint",
	)
	require.NoError(t, err)

	buffer := bytes.Buffer{}
	err = runCdtOutput(t.Context(), io.MultiWriter(os.Stdout, &buffer), os.Stderr,
		"-e", "pyenv:"+envDir,
		"-w", "python", "-r", "data/python", "-b", buildDir,
		"lint",
	)
	require.NoError(t, err)
}
