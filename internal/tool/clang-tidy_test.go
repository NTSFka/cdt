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

func TestClangTidy_DetectClangTidy(t *testing.T) {
	env := test.NewEnvironment(t)
	env.OnFindExecutable("clang-tidy").
		Return(env.NewExecutable("/bin/clang-tidy"))

	tool := DetectClangTidy(env)
	assert.NotNil(t, tool)
	assert.Equal(t, "clang-tidy", tool.Id())
	assert.True(t, tool.IsAvailable())

	if path := tool.ExecutablePath(); assert.NotNil(t, path) {
		assert.Equal(t, "/bin/clang-tidy", *path)
	}

	env.AssertExpectations(t)
}

func TestClangTidy_DetectClangTidy_NotFound(t *testing.T) {
	env := test.NewEnvironment(t)
	env.OnFindExecutable("clang-tidy").
		Return(nil)

	tool := DetectClangTidy(env)
	assert.NotNil(t, tool)
	assert.Equal(t, "clang-tidy", tool.Id())
	assert.False(t, tool.IsAvailable())
	assert.Nil(t, tool.ExecutablePath())

	env.AssertExpectations(t)
}

func TestClangTidy_LintAll(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewClangTidy(exec.LazyExecutable("clang-tidy"))

	p := internal.MakeProject("project", "build", &internal.FixedProjectStructureProvider{
		ProjectStructure: internal.ProjectStructure{
			Targets: map[string]internal.ProjectTarget{
				"target1": {
					Files: []string{"file1.go", "file2.go"},
				},
			},
		},
	}, internal.Workflow{})

	exec.OnRun("clang-tidy", []string{
		"-p", p.BuildDirectory(),
		filepath.Join(p.RootDirectory(), "file1.go"),
		filepath.Join(p.RootDirectory(), "file2.go"),
	}).
		Return(nil)

	err := tool.LintAll(p, []string{})
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestClangTidy_LintAll_Failed(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewClangTidy(exec.LazyExecutable("clang-tidy"))

	p := internal.MakeProject("project", "build", &internal.FixedProjectStructureProvider{
		ProjectStructure: internal.ProjectStructure{
			Targets: map[string]internal.ProjectTarget{
				"target1": {
					Files: []string{"file1.go", "file2.go"},
				},
			},
		},
	}, internal.Workflow{})

	exec.OnRun("clang-tidy", []string{
		"-p", p.BuildDirectory(),
		filepath.Join(p.RootDirectory(), "file1.go"),
		filepath.Join(p.RootDirectory(), "file2.go"),
	}).
		Return(errors.New("failed"))

	err := tool.LintAll(p, []string{})
	assert.EqualError(t, err, "failed")

	exec.AssertExpectations(t)
}

func TestClangTidy_LintAll_CustomConfig(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewClangTidy(exec.LazyExecutable("clang-tidy"))

	p := internal.MakeProject(t.TempDir(), "build", &internal.FixedProjectStructureProvider{
		ProjectStructure: internal.ProjectStructure{
			Targets: map[string]internal.ProjectTarget{
				"target1": {
					Files: []string{"file1.go", "file2.go"},
				},
			},
		},
	}, internal.Workflow{})

	_, err := os.Create(filepath.Join(p.RootDirectory(), ".clang-tidy"))
	assert.NoError(t, err)

	exec.OnRun("clang-tidy", []string{
		fmt.Sprintf("--config-file=%v", filepath.Join(p.RootDirectory(), ".clang-tidy")),
		"-p", p.BuildDirectory(),
		filepath.Join(p.RootDirectory(), "file1.go"),
		filepath.Join(p.RootDirectory(), "file2.go"),
	}).
		Return(nil)

	err = tool.LintAll(p, []string{})
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestClangTidy_LintFiles(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewClangTidy(exec.LazyExecutable("clang-tidy"))

	p := internal.MakeProject(t.TempDir(), "build", nil, internal.Workflow{})

	exec.OnRun("clang-tidy", []string{
		"-p", p.BuildDirectory(),
		filepath.Join(p.RootDirectory(), "file1.go"),
		filepath.Join(p.RootDirectory(), "file3.go"),
	}).
		Return(nil)

	err := tool.LintFiles(p, []string{"file1.go", filepath.Join(p.RootDirectory(), "file3.go")}, []string{})
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestClangTidy_LintFiles_Failed(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewClangTidy(exec.LazyExecutable("clang-tidy"))

	p := internal.MakeProject("project", "build", nil, internal.Workflow{})

	exec.OnRun("clang-tidy", []string{
		"-p", p.BuildDirectory(),
		filepath.Join(p.RootDirectory(), "file1.go"),
	}).
		Return(errors.New("failed"))

	err := tool.LintFiles(p, []string{"file1.go"}, []string{})
	assert.EqualError(t, err, "failed")

	exec.AssertExpectations(t)
}

func TestClangTidy_LintFiles_CustomConfig(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewClangTidy(exec.LazyExecutable("clang-tidy"))

	p := internal.MakeProject(t.TempDir(), "build", nil, internal.Workflow{})

	_, err := os.Create(filepath.Join(p.RootDirectory(), ".clang-tidy"))
	assert.NoError(t, err)

	exec.OnRun("clang-tidy", []string{
		fmt.Sprintf("--config-file=%v", filepath.Join(p.RootDirectory(), ".clang-tidy")),
		"-p", p.BuildDirectory(),
		filepath.Join(p.RootDirectory(), "file1.go"),
	}).
		Return(nil)

	err = tool.LintFiles(p, []string{"file1.go"}, []string{})
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestClangTidy_Run(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewClangTidy(exec.LazyExecutable("clang-tidy"))

	p := internal.MakeProject("project", "build", nil, internal.Workflow{})

	exec.OnRun("clang-tidy", []string{p.RootDirectory()}).
		Return(nil)

	err := tool.RunForProject(p, []string{})
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestClangTidy_Run_Failed(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewClangTidy(exec.LazyExecutable("clang-tidy"))

	p := internal.MakeProject("project", "build", nil, internal.Workflow{})

	exec.OnRun("clang-tidy", []string{p.RootDirectory()}).
		Return(errors.New("failed"))

	err := tool.RunForProject(p, []string{})
	assert.EqualError(t, err, "failed")

	exec.AssertExpectations(t)
}
