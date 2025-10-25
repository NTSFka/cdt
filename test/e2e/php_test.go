package e2e_test

import (
	"bytes"
	"io"
	"os"
	"testing"

	"cdt/internal"

	"github.com/stretchr/testify/require"
)

func TestPhpRun(t *testing.T) {
	environment := internal.SystemEnvironment

	checkTool(t, t.Context(), environment, "php")

	outputDir := t.TempDir()

	var err error

	buffer := bytes.Buffer{}
	err = runCdtOutput(t.Context(), io.MultiWriter(os.Stdout, &buffer), os.Stderr,
		"-w", "php", "-r", "data/php", "-o", outputDir,
		"run", "bin/hello.php",
	)
	require.NoError(t, err)
	require.Contains(t, buffer.String(), "Hello World!\n")
}

func TestPhpTest(t *testing.T) {
	// FIXME: working directory problems
	t.Skip("working directory problems")

	environment := internal.SystemEnvironment

	checkTool(t, t.Context(), environment, "php")

	outputDir := t.TempDir()

	// Install dependencies
	var err error
	err = runCdt(t.Context(), "-w", "php", "-r", "data/php", "-o", outputDir, "dep", "install")
	require.NoError(t, err)

	buffer := bytes.Buffer{}
	err = runCdtOutput(t.Context(), io.MultiWriter(os.Stdout, &buffer), os.Stderr,
		"-w", "php", "-r", "data/php", "-o", outputDir,
		"test",
	)
	require.NoError(t, err)
	require.Contains(t, buffer.String(), "ok")
	require.NotContains(t, buffer.String(), "FAIL")
}

func TestPhpFormat(t *testing.T) {
	// FIXME: working directory problems
	t.Skip("working directory problems")

	environment := internal.SystemEnvironment

	checkTool(t, t.Context(), environment, "php")

	outputDir := t.TempDir()

	// Install dependencies
	var err error
	err = runCdt(t.Context(), "-w", "php", "-r", "data/php", "-o", outputDir, "dep", "install")
	require.NoError(t, err)

	buffer := bytes.Buffer{}
	err = runCdtOutput(t.Context(), io.MultiWriter(os.Stdout, &buffer), os.Stderr,
		"-w", "php", "-r", "data/php", "-o", outputDir,
		"format",
	)
	require.NoError(t, err)
}

func TestPhpLint(t *testing.T) {
	// FIXME: working directory problems
	t.Skip("working directory problems")

	environment := internal.SystemEnvironment

	checkTool(t, t.Context(), environment, "php")

	outputDir := t.TempDir()

	// Install dependencies
	var err error
	err = runCdt(t.Context(), "-w", "php", "-r", "data/php", "-o", outputDir, "dep", "install")
	require.NoError(t, err)

	buffer := bytes.Buffer{}
	err = runCdtOutput(t.Context(), io.MultiWriter(os.Stdout, &buffer), os.Stderr,
		"-w", "php", "-r", "data/php", "-o", outputDir,
		"lint",
	)
	require.NoError(t, err)
}
