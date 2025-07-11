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
	environment.OnFindExecutable("clang-format").
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

func TestClangFormat_DetectClangFormat_Version(t *testing.T) {
	environment := test.Environment{}
	environment.OnFindExecutable("clang-format").
		Return(nil)
	environment.OnFindExecutable(fmt.Sprintf("clang-format-%v", clangTidyMaxVersion)).
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
	environment.OnFindExecutable("clang-format-20").
		Return(environment.MakeExecutable("/bin/clang-format-20"))

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
	environment.OnFindExecutable("clang-format").
		Return(nil)

	for version := clangTidyMaxVersion; version >= clangTidyMinVersion; version-- {
		environment.OnFindExecutable(fmt.Sprintf("clang-format-%v", version)).
			Return(nil)
	}

	tool := DetectClangFormat(&environment, nil)
	assert.NotNil(t, tool)
	assert.Equal(t, "clang-format", tool.Id())
	assert.False(t, tool.IsAvailable())
	assert.Nil(t, tool.ExecutablePath())

	environment.AssertExpectations(t)
}

func TestClangFormat_FormatAll(t *testing.T) {
	runMock := test.Executable{}

	tool := NewClangFormat(runMock.LazyExecutable("clang-format"))

	p := internal.MakeProject("project", "build", &internal.FixedProjectStructureProvider{
		ProjectStructure: internal.ProjectStructure{
			Targets: map[string]internal.ProjectTarget{
				"target1": {
					Files: []string{"file1.go", "file2.go"},
				},
			},
		},
	}, internal.Workflow{})

	runMock.OnRun("clang-format", []string{
		"--Werror",
		"-i",
		filepath.Join(p.RootDirectory(), "file1.go"),
		filepath.Join(p.RootDirectory(), "file2.go"),
	}).
		Return(nil)

	err := tool.FormatAll(p, []string{})
	assert.NoError(t, err)

	runMock.AssertExpectations(t)
}

func TestClangFormat_FormatAll_Failed(t *testing.T) {
	runMock := test.Executable{}

	tool := NewClangFormat(runMock.LazyExecutable("clang-format"))

	p := internal.MakeProject("project", "build", &internal.FixedProjectStructureProvider{
		ProjectStructure: internal.ProjectStructure{
			Targets: map[string]internal.ProjectTarget{
				"target1": {
					Files: []string{"file1.go", "file2.go"},
				},
			},
		},
	}, internal.Workflow{})

	runMock.OnRun("clang-format", []string{
		"--Werror",
		"-i",
		filepath.Join(p.RootDirectory(), "file1.go"),
		filepath.Join(p.RootDirectory(), "file2.go"),
	}).
		Return(errors.New("failed"))

	err := tool.FormatAll(p, []string{})
	assert.EqualError(t, err, "failed")

	runMock.AssertExpectations(t)
}

func TestClangFormat_FormatAll_CustomConfig(t *testing.T) {
	runMock := test.Executable{}

	tool := NewClangFormat(runMock.LazyExecutable("clang-format"))

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

	runMock.OnRun("clang-format", []string{
		fmt.Sprintf("--style=file:%v", filepath.Join(p.RootDirectory(), ".clang-format")),
		"--Werror",
		"-i",
		filepath.Join(p.RootDirectory(), "file1.go"),
		filepath.Join(p.RootDirectory(), "file2.go"),
	}).
		Return(nil)

	err = tool.FormatAll(p, []string{})
	assert.NoError(t, err)

	runMock.AssertExpectations(t)
}

func TestClangFormat_FormatFiles(t *testing.T) {
	runMock := test.Executable{}

	tool := NewClangFormat(runMock.LazyExecutable("clang-format"))

	p := internal.MakeProject(t.TempDir(), "build", nil, internal.Workflow{})

	runMock.OnRun("clang-format", []string{
		"--Werror",
		"-i",
		filepath.Join(p.RootDirectory(), "file1.go"),
		filepath.Join(p.RootDirectory(), "file3.go"),
	}).
		Return(nil)

	err := tool.FormatFiles(p, []string{"file1.go", filepath.Join(p.RootDirectory(), "file3.go")}, []string{})
	assert.NoError(t, err)

	runMock.AssertExpectations(t)
}

func TestClangFormat_FormatFiles_Failed(t *testing.T) {
	runMock := test.Executable{}

	tool := NewClangFormat(runMock.LazyExecutable("clang-format"))

	p := internal.MakeProject("project", "build", nil, internal.Workflow{})

	runMock.OnRun("clang-format", []string{
		"--Werror",
		"-i",
		filepath.Join(p.RootDirectory(), "file1.go"),
	}).
		Return(errors.New("failed"))

	err := tool.FormatFiles(p, []string{"file1.go"}, []string{})
	assert.EqualError(t, err, "failed")

	runMock.AssertExpectations(t)
}

func TestClangFormat_FormatFiles_CustomConfig(t *testing.T) {
	runMock := test.Executable{}

	tool := NewClangFormat(runMock.LazyExecutable("clang-format"))

	p := internal.MakeProject(t.TempDir(), "build", nil, internal.Workflow{})

	_, err := os.Create(filepath.Join(p.RootDirectory(), ".clang-format"))
	assert.NoError(t, err)

	runMock.OnRun("clang-format", []string{
		fmt.Sprintf("--style=file:%v", filepath.Join(p.RootDirectory(), ".clang-format")),
		"--Werror",
		"-i",
		filepath.Join(p.RootDirectory(), "file1.go"),
	}).
		Return(nil)

	err = tool.FormatFiles(p, []string{"file1.go"}, []string{})
	assert.NoError(t, err)

	runMock.AssertExpectations(t)
}

func TestClangFormat_FormatCheckAll(t *testing.T) {
	runMock := test.Executable{}

	tool := NewClangFormat(runMock.LazyExecutable("clang-format"))

	p := internal.MakeProject("project", "build", &internal.FixedProjectStructureProvider{
		ProjectStructure: internal.ProjectStructure{
			Targets: map[string]internal.ProjectTarget{
				"target1": {
					Files: []string{"file1.go", "file2.go"},
				},
			},
		},
	}, internal.Workflow{})

	runMock.OnRun("clang-format", []string{
		"--Werror",
		"--dry-run",
		filepath.Join(p.RootDirectory(), "file1.go"),
		filepath.Join(p.RootDirectory(), "file2.go"),
	}).
		Return(nil)

	err := tool.FormatCheckAll(p, []string{})
	assert.NoError(t, err)

	runMock.AssertExpectations(t)
}

func TestClangFormat_FormatCheckAll_Failed(t *testing.T) {
	runMock := test.Executable{}

	tool := NewClangFormat(runMock.LazyExecutable("clang-format"))

	p := internal.MakeProject("project", "build", &internal.FixedProjectStructureProvider{
		ProjectStructure: internal.ProjectStructure{
			Targets: map[string]internal.ProjectTarget{
				"target1": {
					Files: []string{"file1.go", "file2.go"},
				},
			},
		},
	}, internal.Workflow{})

	runMock.OnRun("clang-format", []string{
		"--Werror",
		"--dry-run",
		filepath.Join(p.RootDirectory(), "file1.go"),
		filepath.Join(p.RootDirectory(), "file2.go"),
	}).
		Return(errors.New("failed"))

	err := tool.FormatCheckAll(p, []string{})
	assert.EqualError(t, err, "failed")

	runMock.AssertExpectations(t)
}

func TestClangFormat_FormatCheckAll_CustomConfig(t *testing.T) {
	runMock := test.Executable{}

	tool := NewClangFormat(runMock.LazyExecutable("clang-format"))

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

	runMock.OnRun("clang-format", []string{
		fmt.Sprintf("--style=file:%v", filepath.Join(p.RootDirectory(), ".clang-format")),
		"--Werror",
		"--dry-run",
		filepath.Join(p.RootDirectory(), "file1.go"),
		filepath.Join(p.RootDirectory(), "file2.go"),
	}).
		Return(nil)

	err = tool.FormatCheckAll(p, []string{})
	assert.NoError(t, err)

	runMock.AssertExpectations(t)
}

func TestClangFormat_FormatCheckFiles(t *testing.T) {
	runMock := test.Executable{}

	tool := NewClangFormat(runMock.LazyExecutable("clang-format"))

	p := internal.MakeProject(t.TempDir(), "build", nil, internal.Workflow{})

	runMock.OnRun("clang-format", []string{
		"--Werror",
		"--dry-run",
		filepath.Join(p.RootDirectory(), "file1.go"),
		filepath.Join(p.RootDirectory(), "file3.go"),
	}).
		Return(nil)

	err := tool.FormatCheckFiles(p, []string{"file1.go", filepath.Join(p.RootDirectory(), "file3.go")}, []string{})
	assert.NoError(t, err)

	runMock.AssertExpectations(t)
}

func TestClangFormat_FormatCheckFiles_Failed(t *testing.T) {
	runMock := test.Executable{}

	tool := NewClangFormat(runMock.LazyExecutable("clang-format"))

	p := internal.MakeProject("project", "build", nil, internal.Workflow{})

	runMock.OnRun("clang-format", []string{
		"--Werror",
		"--dry-run",
		filepath.Join(p.RootDirectory(), "file1.go"),
	}).
		Return(errors.New("failed"))

	err := tool.FormatCheckFiles(p, []string{"file1.go"}, []string{})
	assert.EqualError(t, err, "failed")

	runMock.AssertExpectations(t)
}

func TestClangFormat_FormatCheckFiles_CustomConfig(t *testing.T) {
	runMock := test.Executable{}

	tool := NewClangFormat(runMock.LazyExecutable("clang-format"))

	p := internal.MakeProject(t.TempDir(), "build", nil, internal.Workflow{})

	_, err := os.Create(filepath.Join(p.RootDirectory(), ".clang-format"))
	assert.NoError(t, err)

	runMock.OnRun("clang-format", []string{
		fmt.Sprintf("--style=file:%v", filepath.Join(p.RootDirectory(), ".clang-format")),
		"--Werror",
		"--dry-run",
		filepath.Join(p.RootDirectory(), "file1.go"),
	}).
		Return(nil)

	err = tool.FormatCheckFiles(p, []string{"file1.go"}, []string{})
	assert.NoError(t, err)

	runMock.AssertExpectations(t)
}

func TestClangFormat_Run(t *testing.T) {
	runMock := test.Executable{}

	tool := NewClangFormat(runMock.LazyExecutable("clang-format"))

	p := internal.MakeProject("project", "build", nil, internal.Workflow{})

	runMock.OnRun("clang-format", []string{}).
		Return(nil)

	err := tool.Run(p, []string{})
	assert.NoError(t, err)

	runMock.AssertExpectations(t)
}

func TestClangFormat_Run_Failed(t *testing.T) {
	runMock := test.Executable{}

	tool := NewClangFormat(runMock.LazyExecutable("clang-format"))

	p := internal.MakeProject("project", "build", nil, internal.Workflow{})

	runMock.OnRun("clang-format", []string{}).
		Return(errors.New("failed"))

	err := tool.Run(p, []string{})
	assert.EqualError(t, err, "failed")

	runMock.AssertExpectations(t)
}
