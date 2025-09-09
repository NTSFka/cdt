package command

import (
	"cdt/internal"
	"cdt/internal/test"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func runProject(structureProvider internal.ProjectStructureProvider, args ...string) error {
	return test.RunCommand(NewProjectCommand(), internal.Context{
		Project: internal.Project{Directory: "", StructureProvider: structureProvider},
	}, args...)
}

func TestProjectTargets(t *testing.T) {
	structure := test.NewStructureProvider(t)
	structure.On("Structure", mock.Anything).Return(&internal.ProjectStructure{}, nil)

	err := runProject(structure, "targets")

	assert.NoError(t, err)
	structure.AssertExpectations(t)
}

func TestProjectFiles(t *testing.T) {
	structure := test.NewStructureProvider(t)
	structure.On("Structure", mock.Anything).Return(&internal.ProjectStructure{}, nil)

	err := runProject(structure, "files")

	assert.NoError(t, err)
	structure.AssertExpectations(t)
}

// TODO: other variants
