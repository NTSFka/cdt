package cli

import (
	"cdt/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func runProject(structureProvider internal.ProjectStructureProvider, args ...string) error {
	var runArgs []string
	runArgs = append(runArgs, "project")
	runArgs = append(runArgs, args...)

	return runMainWithProject(internal.MakeProject(
		".",
		"",
		structureProvider,
		internal.Workflow{},
	), runArgs...)
}

func TestProjectTargets(t *testing.T) {
	structure := testStructureProvider{}
	structure.On("Structure", mock.Anything).Return(&internal.ProjectStructure{}, nil)

	err := runProject(&structure, "targets")

	assert.NoError(t, err)
	structure.AssertExpectations(t)
}

func TestProjectFiles(t *testing.T) {
	structure := testStructureProvider{}
	structure.On("Structure", mock.Anything).Return(&internal.ProjectStructure{}, nil)

	err := runProject(&structure, "files")

	assert.NoError(t, err)
	structure.AssertExpectations(t)
}

// TODO: other variants
