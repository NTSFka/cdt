package tool

import (
	"cdt/internal"
	"cdt/internal/test"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClangTidy_DetectClangTidy(t *testing.T) {
	env := test.NewEnvironment(t)
	env.OnFindExecutable("clang-tidy").
		Return(env.NewExecutable("/bin/clang-tidy"))

	tool := DetectClangTidy(env)
	assert.NotNil(t, tool)
	assert.Equal(t, "clang-tidy", tool.Id())
	assert.True(t, tool.IsAvailable())

	if executable := tool.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/clang-tidy", executable.Path)
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
	assert.Nil(t, tool.Executable())

	env.AssertExpectations(t)
}

func TestClangTidy_LintAll(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewClangTidy(exec.LazyExecutable("clang-tidy"))

	desc := internal.Project{
		Directory:             "project",
		IntermediateDirectory: internal.StrPtr("build"),
		StructureProvider: &internal.FixedProjectStructureProvider{
			ProjectStructure: internal.ProjectStructure{
				Targets: map[string]internal.ProjectTarget{
					"target1": {
						Files: []string{"file1.go", "file2.go"},
					},
				},
			},
		},
	}

	exec.OnRun("clang-tidy", []string{
		"-p", *desc.IntermediateDirectory,
		filepath.Join(desc.Directory, "file1.go"),
		filepath.Join(desc.Directory, "file2.go"),
	}).
		Return(nil)

	err := tool.LintAll(desc, []string{})
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestClangTidy_LintAll_Failed(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewClangTidy(exec.LazyExecutable("clang-tidy"))

	desc := internal.Project{
		Directory:             "project",
		IntermediateDirectory: internal.StrPtr("build"),
		StructureProvider: &internal.FixedProjectStructureProvider{
			ProjectStructure: internal.ProjectStructure{
				Targets: map[string]internal.ProjectTarget{
					"target1": {
						Files: []string{"file1.go", "file2.go"},
					},
				},
			},
		},
	}

	exec.OnRun("clang-tidy", []string{
		"-p", *desc.IntermediateDirectory,
		filepath.Join(desc.Directory, "file1.go"),
		filepath.Join(desc.Directory, "file2.go"),
	}).
		Return(errors.New("failed"))

	err := tool.LintAll(desc, []string{})
	assert.EqualError(t, err, "failed")

	exec.AssertExpectations(t)
}

func TestClangTidy_LintAll_CustomConfig(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewClangTidy(exec.LazyExecutable("clang-tidy"))

	desc := internal.Project{
		Directory:             t.TempDir(),
		IntermediateDirectory: internal.StrPtr("build"),
		StructureProvider: &internal.FixedProjectStructureProvider{
			ProjectStructure: internal.ProjectStructure{
				Targets: map[string]internal.ProjectTarget{
					"target1": {
						Files: []string{"file1.go", "file2.go"},
					},
				},
			},
		},
	}

	_, err := os.Create(filepath.Join(desc.Directory, ".clang-tidy"))
	assert.NoError(t, err)

	exec.OnRun("clang-tidy", []string{
		fmt.Sprintf("--config-file=%v", filepath.Join(desc.Directory, ".clang-tidy")),
		"-p", *desc.IntermediateDirectory,
		filepath.Join(desc.Directory, "file1.go"),
		filepath.Join(desc.Directory, "file2.go"),
	}).
		Return(nil)

	err = tool.LintAll(desc, []string{})
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestClangTidy_LintFiles(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewClangTidy(exec.LazyExecutable("clang-tidy"))

	desc := internal.Project{Directory: t.TempDir(), IntermediateDirectory: internal.StrPtr("build")}

	exec.OnRun("clang-tidy", []string{
		"-p", *desc.IntermediateDirectory,
		filepath.Join(desc.Directory, "file1.go"),
		filepath.Join(desc.Directory, "file3.go"),
	}).
		Return(nil)

	err := tool.LintFiles(desc, []string{"file1.go", filepath.Join(desc.Directory, "file3.go")}, []string{})
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestClangTidy_LintFiles_Failed(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewClangTidy(exec.LazyExecutable("clang-tidy"))

	desc := internal.Project{Directory: "project", IntermediateDirectory: internal.StrPtr("build")}

	exec.OnRun("clang-tidy", []string{
		"-p", *desc.IntermediateDirectory,
		filepath.Join(desc.Directory, "file1.go"),
	}).
		Return(errors.New("failed"))

	err := tool.LintFiles(desc, []string{"file1.go"}, []string{})
	assert.EqualError(t, err, "failed")

	exec.AssertExpectations(t)
}

func TestClangTidy_LintFiles_CustomConfig(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewClangTidy(exec.LazyExecutable("clang-tidy"))

	desc := internal.Project{Directory: t.TempDir(), IntermediateDirectory: internal.StrPtr("build")}

	_, err := os.Create(filepath.Join(desc.Directory, ".clang-tidy"))
	assert.NoError(t, err)

	exec.OnRun("clang-tidy", []string{
		fmt.Sprintf("--config-file=%v", filepath.Join(desc.Directory, ".clang-tidy")),
		"-p", *desc.IntermediateDirectory,
		filepath.Join(desc.Directory, "file1.go"),
	}).
		Return(nil)

	err = tool.LintFiles(desc, []string{"file1.go"}, []string{})
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestClangTidy_Run(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewClangTidy(exec.LazyExecutable("clang-tidy"))

	desc := internal.Project{Directory: "build"}

	exec.OnRun("clang-tidy", []string{desc.Directory}).
		Return(nil)

	err := tool.RunForProject(desc, []string{})
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestClangTidy_Run_Failed(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewClangTidy(exec.LazyExecutable("clang-tidy"))

	desc := internal.Project{Directory: "build"}

	exec.OnRun("clang-tidy", []string{desc.Directory}).
		Return(errors.New("failed"))

	err := tool.RunForProject(desc, []string{})
	assert.EqualError(t, err, "failed")

	exec.AssertExpectations(t)
}
