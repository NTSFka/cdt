package command

import (
	"cdt/internal"
	"cdt/internal/test"
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func runProject(ctx context.Context, structureProvider internal.ProjectStructureProvider, args ...string) error {
	return test.RunCommand(ctx, NewProjectCommand(), internal.Context{
		Project: internal.Project{
			Info: internal.ProjectInfo{Directory: "", StructureProvider: structureProvider},
		},
	}, args...)
}

func TestProjectTargets(t *testing.T) {
	structure := test.NewStructureProvider(t)
	structure.On("Structure", mock.Anything, mock.Anything).Return(&internal.ProjectStructure{}, nil)

	err := runProject(context.Background(), structure, "targets")

	require.NoError(t, err)
	structure.AssertExpectations(t)
}

func TestProjectFiles(t *testing.T) {
	structure := test.NewStructureProvider(t)
	structure.On("Structure", mock.Anything, mock.Anything).Return(&internal.ProjectStructure{}, nil)

	err := runProject(context.Background(), structure, "files")

	require.NoError(t, err)
	structure.AssertExpectations(t)
}

// TODO: other variants
