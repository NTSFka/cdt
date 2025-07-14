package internal

import (
	"github.com/stretchr/testify/assert"
	"io"
	"strings"
	"testing"
)

func TestLoadConfigFile_Empty(t *testing.T) {
	reader := strings.NewReader("")

	config, err := LoadConfigFile(reader)
	assert.ErrorIs(t, err, io.EOF)
	assert.Nil(t, config)
}

func TestLoadConfigFile_EmptyProject(t *testing.T) {
	reader := strings.NewReader("project:\n")

	config, err := LoadConfigFile(reader)
	assert.NoError(t, err)

	if assert.NotNil(t, config) {
		if assert.NotNil(t, config.Project) {
			assert.Nil(t, config.Project.Directory)
			assert.Nil(t, config.Project.BuildDirectory)
			assert.Nil(t, config.Project.Environment)
		}
	}
}

func TestLoadConfigFile_Whole(t *testing.T) {
	reader := strings.NewReader(`
project:
    directory: /path/to/project
    build-directory: /path/to/build
    environment: docker:golang
`,
	)

	config, err := LoadConfigFile(reader)
	assert.NoError(t, err)

	if assert.NotNil(t, config) {
		if assert.NotNil(t, config.Project) {
			if assert.NotNil(t, config.Project.Directory) {
				assert.Equal(t, "/path/to/project", *config.Project.Directory)
			}

			if assert.NotNil(t, config.Project.BuildDirectory) {
				assert.Equal(t, "/path/to/build", *config.Project.BuildDirectory)
			}

			if assert.NotNil(t, config.Project.Environment) {
				assert.Equal(t, "docker:golang", *config.Project.Environment)
			}
		}
	}
}
