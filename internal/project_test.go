package internal_test

import (
	"testing"

	"cdt/internal"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProject_ProjectStructure_GetFiles(t *testing.T) {
	structure := internal.ProjectStructure{
		Targets: map[string]internal.ProjectTarget{
			"target1": {
				Files: []string{"file1.go"},
			},
			"target2": {
				Files: []string{"file2.go", "file3.go"},
			},
		},
	}

	assert.ElementsMatch(t, []string{"file1.go", "file2.go", "file3.go"}, structure.GetFiles())
}

func TestProject_EmptyStructureProvider(t *testing.T) {
	provider := internal.EmptyProjectStructureProvider{}
	structure, err := provider.Structure(t.Context(), internal.ProjectInfo{})

	require.EqualError(t, err, "no structure provided")
	assert.Nil(t, structure)
}

func TestProject_FixedStructureProvider(t *testing.T) {
	provider := internal.FixedProjectStructureProvider{
		ProjectStructure: internal.ProjectStructure{
			Targets: map[string]internal.ProjectTarget{
				"target1": {
					Files: []string{"file1.go"},
				},
			},
		},
	}
	structure, err := provider.Structure(t.Context(), internal.ProjectInfo{})

	require.NoError(t, err)
	require.NotNil(t, structure)
	assert.Equal(t, internal.ProjectStructure{
		Targets: map[string]internal.ProjectTarget{
			"target1": {
				Files: []string{"file1.go"},
			},
		},
	}, *structure)
}
