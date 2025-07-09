package tool

import (
	"cdt/internal"
	"cdt/internal/test"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestClangFormat_DetectClangFormat(t *testing.T) {
	environment := test.Environment{}
	environment.On("FindExecutable", "clang-format").Return(environment.MakeExecutable("/bin/clang-format"))

	tool := DetectClangFormat(&environment, nil)
	assert.NotNil(t, tool)
	assert.Equal(t, "clang-format", tool.Id())
	assert.True(t, tool.IsAvailable())

	if path := tool.ExecutablePath(); assert.NotNil(t, path) {
		assert.Equal(t, "/bin/clang-format", *path)
	}

	environment.AssertExpectations(t)
}

func TestClangFormat_DetectClangFormat_Version(t *testing.T) {
	environment := test.Environment{}
	environment.On("FindExecutable", "clang-format").Return(nil)
	environment.On("FindExecutable", fmt.Sprintf("clang-format-%v", clangTidyMaxVersion)).
		Return(environment.MakeExecutable("/bin/clang-format"))

	tool := DetectClangFormat(&environment, nil)
	assert.NotNil(t, tool)
	assert.Equal(t, "clang-format", tool.Id())
	assert.True(t, tool.IsAvailable())

	if path := tool.ExecutablePath(); assert.NotNil(t, path) {
		assert.Equal(t, "/bin/clang-format", *path)
	}

	environment.AssertExpectations(t)
}

func TestClangFormat_DetectClangFormat_Preferred(t *testing.T) {
	environment := test.Environment{}
	environment.On("FindExecutable", "clang-format-20").Return(environment.MakeExecutable("/bin/clang-format-20"))

	version := 20
	tool := DetectClangFormat(&environment, &version)
	assert.NotNil(t, tool)
	assert.Equal(t, "clang-format", tool.Id())
	assert.True(t, tool.IsAvailable())

	if path := tool.ExecutablePath(); assert.NotNil(t, path) {
		assert.Equal(t, "/bin/clang-format-20", *path)
	}

	environment.AssertExpectations(t)
}

func TestClangFormat_DetectClangFormat_NotFound(t *testing.T) {
	environment := test.Environment{}
	environment.On("FindExecutable", "clang-format").Return(nil)

	for version := clangTidyMaxVersion; version >= clangTidyMinVersion; version-- {
		environment.On("FindExecutable", fmt.Sprintf("clang-format-%v", version)).Return(nil)
	}

	tool := DetectClangFormat(&environment, nil)
	assert.NotNil(t, tool)
	assert.Equal(t, "clang-format", tool.Id())
	assert.False(t, tool.IsAvailable())
	assert.Nil(t, tool.ExecutablePath())

	environment.AssertExpectations(t)
}

func TestClangFormat_FormatAll(t *testing.T) {
	environment := test.Environment{}

	tool := NewClangFormat(environment.DetectExecutable("clang-format"))

	p := internal.MakeProject("project", "build", &internal.FixedProjectStructureProvider{
		ProjectStructure: internal.ProjectStructure{
			Targets: map[string]internal.ProjectTarget{
				"target1": {
					Files: []string{"file1.go", "file2.go"},
				},
			},
		},
	}, internal.Workflow{})

	environment.OnRunSuccess(p, "clang-format", []string{
		"--Werror",
		"-i",
		filepath.Join(p.RootDirectory(), "file1.go"),
		filepath.Join(p.RootDirectory(), "file2.go"),
	})

	err := tool.FormatAll(p, []string{})
	assert.NoError(t, err)

	environment.AssertExpectations(t)
}

func TestClangFormat_FormatAll_Failed(t *testing.T) {
	environment := test.Environment{}

	tool := NewClangFormat(environment.DetectExecutable("clang-format"))

	p := internal.MakeProject("project", "build", &internal.FixedProjectStructureProvider{
		ProjectStructure: internal.ProjectStructure{
			Targets: map[string]internal.ProjectTarget{
				"target1": {
					Files: []string{"file1.go", "file2.go"},
				},
			},
		},
	}, internal.Workflow{})

	environment.OnRunError(p, "clang-format", []string{
		"--Werror",
		"-i",
		filepath.Join(p.RootDirectory(), "file1.go"),
		filepath.Join(p.RootDirectory(), "file2.go"),
	}, errors.New("failed"))

	err := tool.FormatAll(p, []string{})
	assert.EqualError(t, err, "failed")

	environment.AssertExpectations(t)
}

func TestClangFormat_FormatAll_CustomConfig(t *testing.T) {
	environment := test.Environment{}

	tool := NewClangFormat(environment.DetectExecutable("clang-format"))

	p := internal.MakeProject(t.TempDir(), "build", &internal.FixedProjectStructureProvider{
		ProjectStructure: internal.ProjectStructure{
			Targets: map[string]internal.ProjectTarget{
				"target1": {
					Files: []string{"file1.go", "file2.go"},
				},
			},
		},
	}, internal.Workflow{})

	_, err := os.Create(filepath.Join(p.RootDirectory(), ".clang-format"))
	assert.NoError(t, err)

	environment.OnRunSuccess(p, "clang-format", []string{
		fmt.Sprintf("--style=file:%v", filepath.Join(p.RootDirectory(), ".clang-format")),
		"--Werror",
		"-i",
		filepath.Join(p.RootDirectory(), "file1.go"),
		filepath.Join(p.RootDirectory(), "file2.go"),
	})

	err = tool.FormatAll(p, []string{})
	assert.NoError(t, err)

	environment.AssertExpectations(t)
}

func TestClangFormat_FormatFiles(t *testing.T) {
	environment := test.Environment{}

	tool := NewClangFormat(environment.DetectExecutable("clang-format"))

	p := internal.MakeProject(t.TempDir(), "build", nil, internal.Workflow{})

	environment.OnRunSuccess(p, "clang-format", []string{
		"--Werror",
		"-i",
		filepath.Join(p.RootDirectory(), "file1.go"),
		filepath.Join(p.RootDirectory(), "file3.go"),
	})

	err := tool.FormatFiles(p, []string{"file1.go", filepath.Join(p.RootDirectory(), "file3.go")}, []string{})
	assert.NoError(t, err)

	environment.AssertExpectations(t)
}

func TestClangFormat_FormatFiles_Failed(t *testing.T) {
	environment := test.Environment{}

	tool := NewClangFormat(environment.DetectExecutable("clang-format"))

	p := internal.MakeProject("project", "build", nil, internal.Workflow{})

	environment.OnRunError(p, "clang-format", []string{
		"--Werror",
		"-i",
		filepath.Join(p.RootDirectory(), "file1.go"),
	}, errors.New("failed"))

	err := tool.FormatFiles(p, []string{"file1.go"}, []string{})
	assert.EqualError(t, err, "failed")

	environment.AssertExpectations(t)
}

func TestClangFormat_FormatFiles_CustomConfig(t *testing.T) {
	environment := test.Environment{}

	tool := NewClangFormat(environment.DetectExecutable("clang-format"))

	p := internal.MakeProject(t.TempDir(), "build", nil, internal.Workflow{})

	_, err := os.Create(filepath.Join(p.RootDirectory(), ".clang-format"))
	assert.NoError(t, err)

	environment.OnRunSuccess(p, "clang-format", []string{
		fmt.Sprintf("--style=file:%v", filepath.Join(p.RootDirectory(), ".clang-format")),
		"--Werror",
		"-i",
		filepath.Join(p.RootDirectory(), "file1.go"),
	})

	err = tool.FormatFiles(p, []string{"file1.go"}, []string{})
	assert.NoError(t, err)

	environment.AssertExpectations(t)
}

func TestClangFormat_FormatCheckAll(t *testing.T) {
	environment := test.Environment{}

	tool := NewClangFormat(environment.DetectExecutable("clang-format"))

	p := internal.MakeProject("project", "build", &internal.FixedProjectStructureProvider{
		ProjectStructure: internal.ProjectStructure{
			Targets: map[string]internal.ProjectTarget{
				"target1": {
					Files: []string{"file1.go", "file2.go"},
				},
			},
		},
	}, internal.Workflow{})

	environment.OnRunSuccess(p, "clang-format", []string{
		"--Werror",
		"--dry-run",
		filepath.Join(p.RootDirectory(), "file1.go"),
		filepath.Join(p.RootDirectory(), "file2.go"),
	})

	err := tool.FormatCheckAll(p, []string{})
	assert.NoError(t, err)

	environment.AssertExpectations(t)
}

func TestClangFormat_FormatCheckAll_Failed(t *testing.T) {
	environment := test.Environment{}

	tool := NewClangFormat(environment.DetectExecutable("clang-format"))

	p := internal.MakeProject("project", "build", &internal.FixedProjectStructureProvider{
		ProjectStructure: internal.ProjectStructure{
			Targets: map[string]internal.ProjectTarget{
				"target1": {
					Files: []string{"file1.go", "file2.go"},
				},
			},
		},
	}, internal.Workflow{})

	environment.OnRunError(p, "clang-format", []string{
		"--Werror",
		"--dry-run",
		filepath.Join(p.RootDirectory(), "file1.go"),
		filepath.Join(p.RootDirectory(), "file2.go"),
	}, errors.New("failed"))

	err := tool.FormatCheckAll(p, []string{})
	assert.EqualError(t, err, "failed")

	environment.AssertExpectations(t)
}

func TestClangFormat_FormatCheckAll_CustomConfig(t *testing.T) {
	environment := test.Environment{}

	tool := NewClangFormat(environment.DetectExecutable("clang-format"))

	p := internal.MakeProject(t.TempDir(), "build", &internal.FixedProjectStructureProvider{
		ProjectStructure: internal.ProjectStructure{
			Targets: map[string]internal.ProjectTarget{
				"target1": {
					Files: []string{"file1.go", "file2.go"},
				},
			},
		},
	}, internal.Workflow{})

	_, err := os.Create(filepath.Join(p.RootDirectory(), ".clang-format"))
	assert.NoError(t, err)

	environment.OnRunSuccess(p, "clang-format", []string{
		fmt.Sprintf("--style=file:%v", filepath.Join(p.RootDirectory(), ".clang-format")),
		"--Werror",
		"--dry-run",
		filepath.Join(p.RootDirectory(), "file1.go"),
		filepath.Join(p.RootDirectory(), "file2.go"),
	})

	err = tool.FormatCheckAll(p, []string{})
	assert.NoError(t, err)

	environment.AssertExpectations(t)
}

func TestClangFormat_FormatCheckFiles(t *testing.T) {
	environment := test.Environment{}

	tool := NewClangFormat(environment.DetectExecutable("clang-format"))

	p := internal.MakeProject(t.TempDir(), "build", nil, internal.Workflow{})

	environment.OnRunSuccess(p, "clang-format", []string{
		"--Werror",
		"--dry-run",
		filepath.Join(p.RootDirectory(), "file1.go"),
		filepath.Join(p.RootDirectory(), "file3.go"),
	})

	err := tool.FormatCheckFiles(p, []string{"file1.go", filepath.Join(p.RootDirectory(), "file3.go")}, []string{})
	assert.NoError(t, err)

	environment.AssertExpectations(t)
}

func TestClangFormat_FormatCheckFiles_Failed(t *testing.T) {
	environment := test.Environment{}

	tool := NewClangFormat(environment.DetectExecutable("clang-format"))

	p := internal.MakeProject("project", "build", nil, internal.Workflow{})

	environment.OnRunError(p, "clang-format", []string{
		"--Werror",
		"--dry-run",
		filepath.Join(p.RootDirectory(), "file1.go"),
	}, errors.New("failed"))

	err := tool.FormatCheckFiles(p, []string{"file1.go"}, []string{})
	assert.EqualError(t, err, "failed")

	environment.AssertExpectations(t)
}

func TestClangFormat_FormatCheckFiles_CustomConfig(t *testing.T) {
	environment := test.Environment{}

	tool := NewClangFormat(environment.DetectExecutable("clang-format"))

	p := internal.MakeProject(t.TempDir(), "build", nil, internal.Workflow{})

	_, err := os.Create(filepath.Join(p.RootDirectory(), ".clang-format"))
	assert.NoError(t, err)

	environment.OnRunSuccess(p, "clang-format", []string{
		fmt.Sprintf("--style=file:%v", filepath.Join(p.RootDirectory(), ".clang-format")),
		"--Werror",
		"--dry-run",
		filepath.Join(p.RootDirectory(), "file1.go"),
	})

	err = tool.FormatCheckFiles(p, []string{"file1.go"}, []string{})
	assert.NoError(t, err)

	environment.AssertExpectations(t)
}

func TestClangFormat_Run(t *testing.T) {
	environment := test.Environment{}

	tool := NewClangFormat(environment.DetectExecutable("clang-format"))

	p := internal.MakeProject("project", "build", nil, internal.Workflow{})

	environment.OnRunSuccess(p, "clang-format", []string{})

	err := tool.Run(p, []string{})
	assert.NoError(t, err)

	environment.AssertExpectations(t)
}

func TestClangFormat_Run_Failed(t *testing.T) {
	environment := test.Environment{}

	tool := NewClangFormat(environment.DetectExecutable("clang-format"))

	p := internal.MakeProject("project", "build", nil, internal.Workflow{})

	environment.OnRunError(p, "clang-format", []string{}, errors.New("failed"))

	err := tool.Run(p, []string{})
	assert.EqualError(t, err, "failed")

	environment.AssertExpectations(t)
}
