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
	environment := internal.SystemEnvironment

	checkTool(t, t.Context(), environment, "php")

	outputDir := t.TempDir()

	// Install dependencies
	var err error
	err = runCdt(t.Context(), "-w", "php", "-r", "data/php", "-o", outputDir, "dep", "install")
	require.NoError(t, err)

	buffer := bytes.Buffer{}
	err = runCdtOutput(t.Context(), os.Stdout, io.MultiWriter(os.Stderr, &buffer),
		"-w", "php", "-r", "data/php", "-o", outputDir,
		"test",
	)
	require.NoError(t, err)
	require.Contains(t, buffer.String(), "OK")
}

func TestPhpFormat(t *testing.T) {
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
		"format", "--check",
	)
	require.NoError(t, err)
}

func TestPhpLint(t *testing.T) {
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

func TestPhpRun_Docker(t *testing.T) {
	environment := internal.SystemEnvironment

	checkTool(t, t.Context(), environment, "docker")

	outputDir := t.TempDir()

	var err error

	buffer := bytes.Buffer{}
	err = runCdtOutput(t.Context(), io.MultiWriter(os.Stdout, &buffer), os.Stderr,
		"-e", "d:php:8.4", "-w", "php", "-r", "data/php", "-o", outputDir,
		"run", "bin/hello.php",
	)
	require.NoError(t, err)
	require.Contains(t, buffer.String(), "Hello World!\n")
}

func TestPhpTest_Docker(t *testing.T) {
	environment := internal.SystemEnvironment

	checkTool(t, t.Context(), environment, "docker")

	outputDir := t.TempDir()

	// Install dependencies
	var err error
	err = runCdt(t.Context(), "-w", "php", "-r", "data/php", "-o", outputDir, "dep", "install")
	require.NoError(t, err)

	buffer := bytes.Buffer{}
	err = runCdtOutput(t.Context(), os.Stdout, io.MultiWriter(os.Stderr, &buffer),
		"-e", "d:php:8.4", "-w", "php", "-r", "data/php", "-o", outputDir,
		"test",
	)
	require.NoError(t, err)
	require.Contains(t, buffer.String(), "OK")
}

func TestPhpFormat_Docker(t *testing.T) {
	environment := internal.SystemEnvironment

	checkTool(t, t.Context(), environment, "docker")

	outputDir := t.TempDir()

	// Install dependencies
	var err error
	err = runCdt(t.Context(), "-w", "php", "-r", "data/php", "-o", outputDir, "dep", "install")
	require.NoError(t, err)

	buffer := bytes.Buffer{}
	err = runCdtOutput(t.Context(), io.MultiWriter(os.Stdout, &buffer), os.Stderr,
		"-e", "d:php:8.4", "-w", "php", "-r", "data/php", "-o", outputDir,
		"format", "--check",
	)
	require.NoError(t, err)
}

func TestPhpLint_Docker(t *testing.T) {
	environment := internal.SystemEnvironment

	checkTool(t, t.Context(), environment, "docker")

	outputDir := t.TempDir()

	// Install dependencies
	var err error
	err = runCdt(t.Context(), "-w", "php", "-r", "data/php", "-o", outputDir, "dep", "install")
	require.NoError(t, err)

	buffer := bytes.Buffer{}
	err = runCdtOutput(t.Context(), io.MultiWriter(os.Stdout, &buffer), os.Stderr,
		"-e", "d:php:8.4", "-w", "php", "-r", "data/php", "-o", outputDir,
		"lint",
	)
	require.NoError(t, err)
}
