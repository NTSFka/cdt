package pkg

import (
	"cdt/internal"
	"cdt/internal/test"
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/urfave/cli/v3"
)

func TestApp_Run_ContextBuild_Failed(t *testing.T) {
	app := NewApp(func(cfg internal.Config) (*internal.Context, error) {
		return nil, errors.New("failed")
	})

	// A fake command to do nothing
	app.Commands = append(app.Commands, &cli.Command{Name: "__test__", Action: func(ctx context.Context, cmd *cli.Command) error {
		return nil
	}})

	err := app.Run(context.Background(), []string{"cdt", "__test__"})

	assert.EqualError(t, err, "failed")
}

func TestApp_Run_Environment_Cleanup(t *testing.T) {
	env := test.NewEnvironment(t)

	app := NewApp(func(cfg internal.Config) (*internal.Context, error) {
		return &internal.Context{
			Environment: env,
		}, nil
	})

	// A fake command to do nothing
	app.Commands = append(app.Commands, &cli.Command{Name: "__test__", Action: func(ctx context.Context, cmd *cli.Command) error {
		return nil
	}})

	env.On("Cleanup", mock.Anything).Return(nil)

	err := app.Run(context.Background(), []string{"cdt", "__test__"})

	assert.NoError(t, err)

	env.AssertExpectations(t)
}

func TestApp_Run_Debug(t *testing.T) {
	app := NewApp(func(cfg internal.Config) (*internal.Context, error) {
		return &internal.Context{}, nil
	})

	// A fake command to do nothing
	app.Commands = append(app.Commands, &cli.Command{Name: "__test__", Action: func(ctx context.Context, cmd *cli.Command) error {
		return nil
	}})

	err := app.Run(context.Background(), []string{"cdt", "--debug", "__test__"})

	assert.NoError(t, err)
}

func TestApp_Run_ConfigDefault(t *testing.T) {
	var config *internal.Config

	app := NewApp(func(cfg internal.Config) (*internal.Context, error) {
		config = &cfg

		return &internal.Context{}, nil
	})

	// A fake command to do nothing
	app.Commands = append(app.Commands, &cli.Command{Name: "__test__", Action: func(ctx context.Context, cmd *cli.Command) error {
		return nil
	}})

	err := app.Run(context.Background(), []string{"cdt", "__test__"})

	if assert.NoError(t, err) {
		if assert.NotNil(t, config) {
			assert.Equal(t, internal.Config{
				RootDirectory:  ".",
				WorkDirectory:  nil,
				BuildDirectory: nil,
				Environment:    nil,
				Workflow:       nil,
			}, *config)
		}
	}
}

func TestApp_Run_ConfigFull(t *testing.T) {
	var config *internal.Config

	app := NewApp(func(cfg internal.Config) (*internal.Context, error) {
		config = &cfg

		return &internal.Context{}, nil
	})

	// A fake command to do nothing
	app.Commands = append(app.Commands, &cli.Command{Name: "__test__", Action: func(ctx context.Context, cmd *cli.Command) error {
		return nil
	}})

	err := app.Run(context.Background(), []string{"cdt",
		"--root", "/path/to/project",
		"--build", "/path/to/build",
		"--environment", "env:arg",
		"--type", "go",
		"__test__",
	})

	if assert.NoError(t, err) {
		if assert.NotNil(t, config) {
			assert.Equal(t, internal.Config{
				RootDirectory:  "/path/to/project",
				WorkDirectory:  nil,
				BuildDirectory: internal.StrPtr("/path/to/build"),
				Environment:    internal.StrPtr("env:arg"),
				Workflow:       "go",
			}, *config)
		}
	}
}

func TestApp_Run_ConfigFullAlias(t *testing.T) {
	var config *internal.Config

	app := NewApp(func(cfg internal.Config) (*internal.Context, error) {
		config = &cfg

		return &internal.Context{}, nil
	})

	// A fake command to do nothing
	app.Commands = append(app.Commands, &cli.Command{Name: "__test__", Action: func(ctx context.Context, cmd *cli.Command) error {
		return nil
	}})

	err := app.Run(context.Background(), []string{"cdt",
		"-r", "/path/to/project",
		"-b", "/path/to/build",
		"-e", "env:arg",
		"-t", "go",
		"__test__",
	})

	if assert.NoError(t, err) {
		if assert.NotNil(t, config) {
			assert.Equal(t, internal.Config{
				RootDirectory:  "/path/to/project",
				WorkDirectory:  nil,
				BuildDirectory: internal.StrPtr("/path/to/build"),
				Environment:    internal.StrPtr("env:arg"),
				Workflow:       "go",
			}, *config)
		}
	}
}

func TestApp_Run_ConfigFile_DefaultPath(t *testing.T) {
	var config *internal.Config

	app := NewApp(func(cfg internal.Config) (*internal.Context, error) {
		config = &cfg

		return &internal.Context{}, nil
	})

	// A fake command to do nothing
	app.Commands = append(app.Commands, &cli.Command{Name: "__test__", Action: func(ctx context.Context, cmd *cli.Command) error {
		return nil
	}})

	tempDir := t.TempDir()

	assert.NoError(t, os.WriteFile(filepath.Join(tempDir, ConfigFileName), []byte(`
project:
    work-directory: /path/to/project
    build-directory: /path/to/build
    environment: env:arg
`), 0666))

	err := app.Run(context.Background(), []string{"cdt",
		"-r", tempDir,
		"__test__",
	})

	if assert.NoError(t, err) {
		if assert.NotNil(t, config) {
			assert.Equal(t, internal.Config{
				RootDirectory:  tempDir,
				WorkDirectory:  internal.StrPtr("/path/to/project"),
				BuildDirectory: internal.StrPtr("/path/to/build"),
				Environment:    internal.StrPtr("env:arg"),
			}, *config)
		}
	}
}

func TestApp_Run_ConfigFile_CustomPath(t *testing.T) {
	var config *internal.Config

	app := NewApp(func(cfg internal.Config) (*internal.Context, error) {
		config = &cfg

		return &internal.Context{}, nil
	})

	// A fake command to do nothing
	app.Commands = append(app.Commands, &cli.Command{Name: "__test__", Action: func(ctx context.Context, cmd *cli.Command) error {
		return nil
	}})

	configFilePath := filepath.Join(t.TempDir(), "my-cdt.yml")

	assert.NoError(t, os.WriteFile(configFilePath, []byte(`
project:
    work-directory: /path/to/work
    build-directory: /path/to/build
    environment: env:arg
`), 0666))

	err := app.Run(context.Background(), []string{"cdt",
		"--config", configFilePath,
		"__test__",
	})

	if assert.NoError(t, err) {
		if assert.NotNil(t, config) {
			assert.Equal(t, internal.Config{
				RootDirectory:  ".",
				WorkDirectory:  internal.StrPtr("/path/to/work"),
				BuildDirectory: internal.StrPtr("/path/to/build"),
				Environment:    internal.StrPtr("env:arg"),
			}, *config)
		}
	}
}

func TestApp_Run_ConfigFile_CustomPath_Alias(t *testing.T) {
	var config *internal.Config

	app := NewApp(func(cfg internal.Config) (*internal.Context, error) {
		config = &cfg

		return &internal.Context{}, nil
	})

	// A fake command to do nothing
	app.Commands = append(app.Commands, &cli.Command{Name: "__test__", Action: func(ctx context.Context, cmd *cli.Command) error {
		return nil
	}})

	configFilePath := filepath.Join(t.TempDir(), "my-cdt.yml")

	assert.NoError(t, os.WriteFile(configFilePath, []byte(`
project:
    work-directory: /path/to/project
    build-directory: /path/to/build
    environment: env:arg
`), 0666))

	err := app.Run(context.Background(), []string{"cdt",
		"-c", configFilePath,
		"__test__",
	})

	if assert.NoError(t, err) {
		if assert.NotNil(t, config) {
			assert.Equal(t, internal.Config{
				RootDirectory:  ".",
				WorkDirectory:  internal.StrPtr("/path/to/project"),
				BuildDirectory: internal.StrPtr("/path/to/build"),
				Environment:    internal.StrPtr("env:arg"),
			}, *config)
		}
	}
}

func TestApp_Run_ConfigFile_CustomPath_UnreadableFile(t *testing.T) {
	app := NewApp(func(cfg internal.Config) (*internal.Context, error) {
		return &internal.Context{}, nil
	})

	// A fake command to do nothing
	app.Commands = append(app.Commands, &cli.Command{Name: "__test__", Action: func(ctx context.Context, cmd *cli.Command) error {
		return nil
	}})

	configFilePath := filepath.Join(t.TempDir(), "my-cdt.yml")

	assert.NoError(t, os.WriteFile(configFilePath, []byte(`
Hello world!
`), 0000))

	err := app.Run(context.Background(), []string{"cdt",
		"--config", configFilePath,
		"__test__",
	})

	assert.ErrorContains(t, err, "failed to open configuration file: ")
}

func TestApp_Run_ConfigFile_CustomPath_InvalidContent(t *testing.T) {
	app := NewApp(func(cfg internal.Config) (*internal.Context, error) {
		return &internal.Context{}, nil
	})

	// A fake command to do nothing
	app.Commands = append(app.Commands, &cli.Command{Name: "__test__", Action: func(ctx context.Context, cmd *cli.Command) error {
		return nil
	}})

	configFilePath := filepath.Join(t.TempDir(), "my-cdt.yml")

	assert.NoError(t, os.WriteFile(configFilePath, []byte(`
Hello world!
`), 0666))

	err := app.Run(context.Background(), []string{"cdt",
		"--config", configFilePath,
		"__test__",
	})

	assert.ErrorContains(t, err, "failed to load configuration file: config load failed: yaml")
}
