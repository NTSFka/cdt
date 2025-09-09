package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProject_ProjectStructure_GetFiles(t *testing.T) {
	structure := ProjectStructure{
		Targets: map[string]ProjectTarget{
			"target1": {
				Files: []string{"file1.go"},
			},
			"target2": {
				Files: []string{"file2.go", "file3.go"},
			},
		},
	}

	assert.Equal(t, []string{"file1.go", "file2.go", "file3.go"}, structure.GetFiles())
}

func TestProject_EmptyStructureProvider(t *testing.T) {
	provider := EmptyProjectStructureProvider{}
	structure, err := provider.Structure(Project{})

	assert.NoError(t, err)
	assert.NotNil(t, structure)
	assert.Empty(t, structure.Targets)
}

func TestProject_FixedStructureProvider(t *testing.T) {
	provider := FixedProjectStructureProvider{
		ProjectStructure: ProjectStructure{
			Targets: map[string]ProjectTarget{
				"target1": {
					Files: []string{"file1.go"},
				},
			},
		},
	}
	structure, err := provider.Structure(Project{})

	assert.NoError(t, err)
	if assert.NotNil(t, structure) {
		assert.Equal(t, ProjectStructure{
			Targets: map[string]ProjectTarget{
				"target1": {
					Files: []string{"file1.go"},
				},
			},
		}, *structure)
	}
}
