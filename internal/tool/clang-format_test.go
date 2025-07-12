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
	env := test.NewEnvironment(t)
	env.OnFindExecutable("clang-format").
		Return(env.NewExecutable("/bin/clang-format"))

	tool := DetectClangFormat(env, nil)
	assert.NotNil(t, tool)
	assert.Equal(t, "clang-format", tool.Id())
	assert.True(t, tool.IsAvailable())

	if path := tool.ExecutablePath(); assert.NotNil(t, path) {
		assert.Equal(t, "/bin/clang-format", *path)
	}

	env.AssertExpectations(t)
}

func TestClangFormat_DetectClangFormat_Version(t *testing.T) {
	env := test.NewEnvironment(t)
	env.OnFindExecutable("clang-format").
		Return(nil)
	env.OnFindExecutable(fmt.Sprintf("clang-format-%v", clangTidyMaxVersion)).
		Return(env.NewExecutable("/bin/clang-format"))

	tool := DetectClangFormat(env, nil)
	assert.NotNil(t, tool)
	assert.Equal(t, "clang-format", tool.Id())
	assert.True(t, tool.IsAvailable())

	if path := tool.ExecutablePath(); assert.NotNil(t, path) {
		assert.Equal(t, "/bin/clang-format", *path)
	}

	env.AssertExpectations(t)
}

func TestClangFormat_DetectClangFormat_Preferred(t *testing.T) {
	env := test.NewEnvironment(t)
	env.OnFindExecutable("clang-format-20").
		Return(env.NewExecutable("/bin/clang-format-20"))

	version := 20
	tool := DetectClangFormat(env, &version)
	assert.NotNil(t, tool)
	assert.Equal(t, "clang-format", tool.Id())
	assert.True(t, tool.IsAvailable())

	if path := tool.ExecutablePath(); assert.NotNil(t, path) {
		assert.Equal(t, "/bin/clang-format-20", *path)
	}

	env.AssertExpectations(t)
}

func TestClangFormat_DetectClangFormat_NotFound(t *testing.T) {
	env := test.NewEnvironment(t)
	env.OnFindExecutable("clang-format").
		Return(nil)

	for version := clangTidyMaxVersion; version >= clangTidyMinVersion; version-- {
		env.OnFindExecutable(fmt.Sprintf("clang-format-%v", version)).
			Return(nil)
	}

	tool := DetectClangFormat(env, nil)
	assert.NotNil(t, tool)
	assert.Equal(t, "clang-format", tool.Id())
	assert.False(t, tool.IsAvailable())
	assert.Nil(t, tool.ExecutablePath())

	env.AssertExpectations(t)
}

func TestClangFormat_FormatAll(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewClangFormat(exec.LazyExecutable("clang-format"))

	p := internal.MakeProject("project", "build", &internal.FixedProjectStructureProvider{
		ProjectStructure: internal.ProjectStructure{
			Targets: map[string]internal.ProjectTarget{
				"target1": {
					Files: []string{"file1.go", "file2.go"},
				},
			},
		},
	}, internal.Workflow{})

	exec.OnRun("clang-format", []string{
		"--Werror",
		"-i",
		filepath.Join(p.RootDirectory(), "file1.go"),
		filepath.Join(p.RootDirectory(), "file2.go"),
	}).
		Return(nil)

	err := tool.FormatAll(p, []string{})
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestClangFormat_FormatAll_Failed(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewClangFormat(exec.LazyExecutable("clang-format"))

	p := internal.MakeProject("project", "build", &internal.FixedProjectStructureProvider{
		ProjectStructure: internal.ProjectStructure{
			Targets: map[string]internal.ProjectTarget{
				"target1": {
					Files: []string{"file1.go", "file2.go"},
				},
			},
		},
	}, internal.Workflow{})

	exec.OnRun("clang-format", []string{
		"--Werror",
		"-i",
		filepath.Join(p.RootDirectory(), "file1.go"),
		filepath.Join(p.RootDirectory(), "file2.go"),
	}).
		Return(errors.New("failed"))

	err := tool.FormatAll(p, []string{})
	assert.EqualError(t, err, "failed")

	exec.AssertExpectations(t)
}

func TestClangFormat_FormatAll_CustomConfig(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewClangFormat(exec.LazyExecutable("clang-format"))

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

	exec.OnRun("clang-format", []string{
		fmt.Sprintf("--style=file:%v", filepath.Join(p.RootDirectory(), ".clang-format")),
		"--Werror",
		"-i",
		filepath.Join(p.RootDirectory(), "file1.go"),
		filepath.Join(p.RootDirectory(), "file2.go"),
	}).
		Return(nil)

	err = tool.FormatAll(p, []string{})
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestClangFormat_FormatFiles(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewClangFormat(exec.LazyExecutable("clang-format"))

	p := internal.MakeProject(t.TempDir(), "build", nil, internal.Workflow{})

	exec.OnRun("clang-format", []string{
		"--Werror",
		"-i",
		filepath.Join(p.RootDirectory(), "file1.go"),
		filepath.Join(p.RootDirectory(), "file3.go"),
	}).
		Return(nil)

	err := tool.FormatFiles(p, []string{"file1.go", filepath.Join(p.RootDirectory(), "file3.go")}, []string{})
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestClangFormat_FormatFiles_Failed(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewClangFormat(exec.LazyExecutable("clang-format"))

	p := internal.MakeProject("project", "build", nil, internal.Workflow{})

	exec.OnRun("clang-format", []string{
		"--Werror",
		"-i",
		filepath.Join(p.RootDirectory(), "file1.go"),
	}).
		Return(errors.New("failed"))

	err := tool.FormatFiles(p, []string{"file1.go"}, []string{})
	assert.EqualError(t, err, "failed")

	exec.AssertExpectations(t)
}

func TestClangFormat_FormatFiles_CustomConfig(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewClangFormat(exec.LazyExecutable("clang-format"))

	p := internal.MakeProject(t.TempDir(), "build", nil, internal.Workflow{})

	_, err := os.Create(filepath.Join(p.RootDirectory(), ".clang-format"))
	assert.NoError(t, err)

	exec.OnRun("clang-format", []string{
		fmt.Sprintf("--style=file:%v", filepath.Join(p.RootDirectory(), ".clang-format")),
		"--Werror",
		"-i",
		filepath.Join(p.RootDirectory(), "file1.go"),
	}).
		Return(nil)

	err = tool.FormatFiles(p, []string{"file1.go"}, []string{})
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestClangFormat_FormatCheckAll(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewClangFormat(exec.LazyExecutable("clang-format"))

	p := internal.MakeProject("project", "build", &internal.FixedProjectStructureProvider{
		ProjectStructure: internal.ProjectStructure{
			Targets: map[string]internal.ProjectTarget{
				"target1": {
					Files: []string{"file1.go", "file2.go"},
				},
			},
		},
	}, internal.Workflow{})

	exec.OnRun("clang-format", []string{
		"--Werror",
		"--dry-run",
		filepath.Join(p.RootDirectory(), "file1.go"),
		filepath.Join(p.RootDirectory(), "file2.go"),
	}).
		Return(nil)

	err := tool.FormatCheckAll(p, []string{})
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestClangFormat_FormatCheckAll_Failed(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewClangFormat(exec.LazyExecutable("clang-format"))

	p := internal.MakeProject("project", "build", &internal.FixedProjectStructureProvider{
		ProjectStructure: internal.ProjectStructure{
			Targets: map[string]internal.ProjectTarget{
				"target1": {
					Files: []string{"file1.go", "file2.go"},
				},
			},
		},
	}, internal.Workflow{})

	exec.OnRun("clang-format", []string{
		"--Werror",
		"--dry-run",
		filepath.Join(p.RootDirectory(), "file1.go"),
		filepath.Join(p.RootDirectory(), "file2.go"),
	}).
		Return(errors.New("failed"))

	err := tool.FormatCheckAll(p, []string{})
	assert.EqualError(t, err, "failed")

	exec.AssertExpectations(t)
}

func TestClangFormat_FormatCheckAll_CustomConfig(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewClangFormat(exec.LazyExecutable("clang-format"))

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

	exec.OnRun("clang-format", []string{
		fmt.Sprintf("--style=file:%v", filepath.Join(p.RootDirectory(), ".clang-format")),
		"--Werror",
		"--dry-run",
		filepath.Join(p.RootDirectory(), "file1.go"),
		filepath.Join(p.RootDirectory(), "file2.go"),
	}).
		Return(nil)

	err = tool.FormatCheckAll(p, []string{})
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestClangFormat_FormatCheckFiles(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewClangFormat(exec.LazyExecutable("clang-format"))

	p := internal.MakeProject(t.TempDir(), "build", nil, internal.Workflow{})

	exec.OnRun("clang-format", []string{
		"--Werror",
		"--dry-run",
		filepath.Join(p.RootDirectory(), "file1.go"),
		filepath.Join(p.RootDirectory(), "file3.go"),
	}).
		Return(nil)

	err := tool.FormatCheckFiles(p, []string{"file1.go", filepath.Join(p.RootDirectory(), "file3.go")}, []string{})
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestClangFormat_FormatCheckFiles_Failed(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewClangFormat(exec.LazyExecutable("clang-format"))

	p := internal.MakeProject("project", "build", nil, internal.Workflow{})

	exec.OnRun("clang-format", []string{
		"--Werror",
		"--dry-run",
		filepath.Join(p.RootDirectory(), "file1.go"),
	}).
		Return(errors.New("failed"))

	err := tool.FormatCheckFiles(p, []string{"file1.go"}, []string{})
	assert.EqualError(t, err, "failed")

	exec.AssertExpectations(t)
}

func TestClangFormat_FormatCheckFiles_CustomConfig(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewClangFormat(exec.LazyExecutable("clang-format"))

	p := internal.MakeProject(t.TempDir(), "build", nil, internal.Workflow{})

	_, err := os.Create(filepath.Join(p.RootDirectory(), ".clang-format"))
	assert.NoError(t, err)

	exec.OnRun("clang-format", []string{
		fmt.Sprintf("--style=file:%v", filepath.Join(p.RootDirectory(), ".clang-format")),
		"--Werror",
		"--dry-run",
		filepath.Join(p.RootDirectory(), "file1.go"),
	}).
		Return(nil)

	err = tool.FormatCheckFiles(p, []string{"file1.go"}, []string{})
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestClangFormat_Run(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewClangFormat(exec.LazyExecutable("clang-format"))

	p := internal.MakeProject("project", "build", nil, internal.Workflow{})

	exec.OnRun("clang-format", []string{}).
		Return(nil)

	err := tool.RunForProject(p, []string{})
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestClangFormat_Run_Failed(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewClangFormat(exec.LazyExecutable("clang-format"))

	p := internal.MakeProject("project", "build", nil, internal.Workflow{})

	exec.OnRun("clang-format", []string{}).
		Return(errors.New("failed"))

	err := tool.RunForProject(p, []string{})
	assert.EqualError(t, err, "failed")

	exec.AssertExpectations(t)
}
