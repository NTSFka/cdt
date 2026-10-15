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
	err := runCdt(t.Context(), "-w", "php", "-r", "data/php", "-o", outputDir, "dep", "install")
	require.NoError(t, err)

	buffer := bytes.Buffer{}
	err = runCdtOutput(t.Context(), os.Stdout, io.MultiWriter(os.Stderr, &buffer),
		"-w", "php", "-r", "data/php", "-o", outputDir,
		"test",
	)
	require.NoError(t, err)
	require.Contains(t, buffer.String(), "OK")
}

func TestPhpCoverage(t *testing.T) {
	environment := internal.SystemEnvironment

	checkTool(t, t.Context(), environment, "php")

	outputDir := t.TempDir()

	// Install dependencies
	err := runCdt(t.Context(), "-w", "php", "-r", "data/php", "-o", outputDir, "dep", "install")
	require.NoError(t, err)

	buffer := bytes.Buffer{}
	err = runCdtOutput(t.Context(), os.Stdout, io.MultiWriter(os.Stderr, &buffer),
		"-w", "php", "-r", "data/php", "-o", outputDir,
		"coverage",
	)
	require.NoError(t, err)
	require.Contains(t, buffer.String(), "OK")
}

func TestPhpFormat(t *testing.T) {
	environment := internal.SystemEnvironment

	checkTool(t, t.Context(), environment, "php")

	outputDir := t.TempDir()

	// Install dependencies
	err := runCdt(t.Context(), "-w", "php", "-r", "data/php", "-o", outputDir, "dep", "install")
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
	err := runCdt(t.Context(), "-w", "php", "-r", "data/php", "-o", outputDir, "dep", "install")
	require.NoError(t, err)

	buffer := bytes.Buffer{}
	err = runCdtOutput(t.Context(), io.MultiWriter(os.Stdout, &buffer), os.Stderr,
		"-w", "php", "-r", "data/php", "-o", outputDir,
		"lint",
	)
	require.NoError(t, err)
}

func TestPhpRun_Environment(t *testing.T) {
	if runtime.GOOS == "windows" { // nolint: goconst
		t.Skip("docker issues on Windows")
	}

	dataSet := []struct {
		toolName string
		envName  string
	}{
		{"docker", "docker"},
		{"podman", "podman"},
	}

	environment := internal.SystemEnvironment

	for _, data := range dataSet {
		t.Run(data.toolName, func(t *testing.T) {
			checkTool(t, t.Context(), environment, data.toolName)

			outputDir := t.TempDir()

			var err error

			buffer := bytes.Buffer{}
			err = runCdtOutput(
				t.Context(),
				io.MultiWriter(os.Stdout, &buffer),
				os.Stderr,
				"-e",
				data.envName+":"+phpEnvironmentImage(data.envName),
				"-w",
				"php",
				"-r",
				"data/php",
				"-o",
				outputDir,
				"run",
				"bin/hello.php",
			)
			require.NoError(t, err)
			require.Contains(t, buffer.String(), "Hello World!\n")
		})
	}
}

func TestPhpTest_Environment(t *testing.T) {
	if runtime.GOOS == "windows" { // nolint: goconst
		t.Skip("docker issues on Windows")
	}

	dataSet := []struct {
		toolName string
		envName  string
	}{
		{"docker", "docker"},
		{"podman", "podman"},
	}

	environment := internal.SystemEnvironment

	for _, data := range dataSet {
		t.Run(data.toolName, func(t *testing.T) {
			checkTool(t, t.Context(), environment, data.toolName)

			outputDir := t.TempDir()

			// Install dependencies
			err := runCdt(
				t.Context(),
				"-w",
				"php",
				"-r",
				"data/php",
				"-o",
				outputDir,
				"dep",
				"install",
			)
			require.NoError(t, err)

			buffer := bytes.Buffer{}
			err = runCdtOutput(
				t.Context(),
				os.Stdout,
				io.MultiWriter(os.Stderr, &buffer),
				"-e",
				data.envName+":"+phpEnvironmentImage(data.envName),
				"-w",
				"php",
				"-r",
				"data/php",
				"-o",
				outputDir,
				"test",
			)
			require.NoError(t, err)
			require.Contains(t, buffer.String(), "OK")
		})
	}
}

func TestPhpFormat_Environment(t *testing.T) {
	if runtime.GOOS == "windows" { // nolint: goconst
		t.Skip("docker issues on Windows")
	}

	dataSet := []struct {
		toolName string
		envName  string
	}{
		{"docker", "docker"},
		{"podman", "podman"},
	}

	environment := internal.SystemEnvironment

	for _, data := range dataSet {
		t.Run(data.toolName, func(t *testing.T) {
			checkTool(t, t.Context(), environment, data.toolName)

			outputDir := t.TempDir()

			// Install dependencies
			err := runCdt(
				t.Context(),
				"-w",
				"php",
				"-r",
				"data/php",
				"-o",
				outputDir,
				"dep",
				"install",
			)
			require.NoError(t, err)

			buffer := bytes.Buffer{}
			err = runCdtOutput(t.Context(), io.MultiWriter(os.Stdout, &buffer), os.Stderr,
				"-e",
				data.envName+":"+phpEnvironmentImage(data.envName),
				"-w", "php", "-r", "data/php", "-o", outputDir,
				"format", "--check",
			)
			require.NoError(t, err)
		})
	}
}

func TestPhpLint_Environment(t *testing.T) {
	if runtime.GOOS == "windows" { // nolint: goconst
		t.Skip("docker issues on Windows")
	}

	dataSet := []struct {
		toolName string
		envName  string
	}{
		{"docker", "docker"},
		{"podman", "podman"},
	}

	environment := internal.SystemEnvironment

	for _, data := range dataSet {
		t.Run(data.toolName, func(t *testing.T) {
			checkTool(t, t.Context(), environment, data.toolName)

			outputDir := t.TempDir()

			// Install dependencies
			err := runCdt(
				t.Context(),
				"-w",
				"php",
				"-r",
				"data/php",
				"-o",
				outputDir,
				"dep",
				"install",
			)
			require.NoError(t, err)

			buffer := bytes.Buffer{}
			err = runCdtOutput(
				t.Context(),
				io.MultiWriter(os.Stdout, &buffer),
				os.Stderr,
				"-e",
				data.envName+":"+phpEnvironmentImage(data.envName),
				"-w",
				"php",
				"-r",
				"data/php",
				"-o",
				outputDir,
				"lint",
			)
			require.NoError(t, err)
		})
	}
}

func phpEnvironmentImage(environment string) string {
	if environment == "podman" {
		return "docker.io/php:8.4"
	}

	return "php:8.4"
}
