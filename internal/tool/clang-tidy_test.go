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
	environment := test.Environment{}
	environment.OnFindExecutable("clang-tidy").
		Return(environment.MakeExecutable("/bin/clang-tidy"))

	tool := DetectClangTidy(&environment, nil)
	assert.NotNil(t, tool)
	assert.Equal(t, "clang-tidy", tool.Id())
	assert.True(t, tool.IsAvailable())

	if path := tool.ExecutablePath(); assert.NotNil(t, path) {
		assert.Equal(t, "/bin/clang-tidy", *path)
	}

	environment.AssertExpectations(t)
}

func TestClangTidy_DetectClangTidy_Version(t *testing.T) {
	environment := test.Environment{}
	environment.OnFindExecutable("clang-tidy").
		Return(nil)
	environment.OnFindExecutable(fmt.Sprintf("clang-tidy-%v", clangTidyMaxVersion)).
		Return(environment.MakeExecutable("/bin/clang-tidy"))

	tool := DetectClangTidy(&environment, nil)
	assert.NotNil(t, tool)
	assert.Equal(t, "clang-tidy", tool.Id())
	assert.True(t, tool.IsAvailable())

	if path := tool.ExecutablePath(); assert.NotNil(t, path) {
		assert.Equal(t, "/bin/clang-tidy", *path)
	}

	environment.AssertExpectations(t)
}

func TestClangTidy_DetectClangTidy_Preferred(t *testing.T) {
	environment := test.Environment{}
	environment.OnFindExecutable("clang-tidy-20").
		Return(environment.MakeExecutable("/bin/clang-tidy-20"))

	version := 20
	tool := DetectClangTidy(&environment, &version)
	assert.NotNil(t, tool)
	assert.Equal(t, "clang-tidy", tool.Id())
	assert.True(t, tool.IsAvailable())

	if path := tool.ExecutablePath(); assert.NotNil(t, path) {
		assert.Equal(t, "/bin/clang-tidy-20", *path)
	}

	environment.AssertExpectations(t)
}

func TestClangTidy_DetectClangTidy_NotFound(t *testing.T) {
	environment := test.Environment{}
	environment.OnFindExecutable("clang-tidy").
		Return(nil)

	for version := clangTidyMaxVersion; version >= clangTidyMinVersion; version-- {
		environment.OnFindExecutable(fmt.Sprintf("clang-tidy-%v", version)).Return(nil)
	}

	tool := DetectClangTidy(&environment, nil)
	assert.NotNil(t, tool)
	assert.Equal(t, "clang-tidy", tool.Id())
	assert.False(t, tool.IsAvailable())
	assert.Nil(t, tool.ExecutablePath())

	environment.AssertExpectations(t)
}

func TestClangTidy_LintAll(t *testing.T) {
	runMock := test.Executable{}

	tool := NewClangTidy(runMock.LazyExecutable("clang-tidy"))

	p := internal.MakeProject("project", "build", &internal.FixedProjectStructureProvider{
		ProjectStructure: internal.ProjectStructure{
			Targets: map[string]internal.ProjectTarget{
				"target1": {
					Files: []string{"file1.go", "file2.go"},
				},
			},
		},
	}, internal.Workflow{})

	runMock.OnRun("clang-tidy", []string{
		"-p", p.BuildDirectory(),
		filepath.Join(p.RootDirectory(), "file1.go"),
		filepath.Join(p.RootDirectory(), "file2.go"),
	}).
		Return(nil)

	err := tool.LintAll(p, []string{})
	assert.NoError(t, err)

	runMock.AssertExpectations(t)
}

func TestClangTidy_LintAll_Failed(t *testing.T) {
	runMock := test.Executable{}

	tool := NewClangTidy(runMock.LazyExecutable("clang-tidy"))

	p := internal.MakeProject("project", "build", &internal.FixedProjectStructureProvider{
		ProjectStructure: internal.ProjectStructure{
			Targets: map[string]internal.ProjectTarget{
				"target1": {
					Files: []string{"file1.go", "file2.go"},
				},
			},
		},
	}, internal.Workflow{})

	runMock.OnRun("clang-tidy", []string{
		"-p", p.BuildDirectory(),
		filepath.Join(p.RootDirectory(), "file1.go"),
		filepath.Join(p.RootDirectory(), "file2.go"),
	}).
		Return(errors.New("failed"))

	err := tool.LintAll(p, []string{})
	assert.EqualError(t, err, "failed")

	runMock.AssertExpectations(t)
}

func TestClangTidy_LintAll_CustomConfig(t *testing.T) {
	runMock := test.Executable{}

	tool := NewClangTidy(runMock.LazyExecutable("clang-tidy"))

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

	runMock.OnRun("clang-tidy", []string{
		fmt.Sprintf("--config-file=%v", filepath.Join(p.RootDirectory(), ".clang-tidy")),
		"-p", p.BuildDirectory(),
		filepath.Join(p.RootDirectory(), "file1.go"),
		filepath.Join(p.RootDirectory(), "file2.go"),
	}).
		Return(nil)

	err = tool.LintAll(p, []string{})
	assert.NoError(t, err)

	runMock.AssertExpectations(t)
}

func TestClangTidy_LintFiles(t *testing.T) {
	runMock := test.Executable{}

	tool := NewClangTidy(runMock.LazyExecutable("clang-tidy"))

	p := internal.MakeProject(t.TempDir(), "build", nil, internal.Workflow{})

	runMock.OnRun("clang-tidy", []string{
		"-p", p.BuildDirectory(),
		filepath.Join(p.RootDirectory(), "file1.go"),
		filepath.Join(p.RootDirectory(), "file3.go"),
	}).
		Return(nil)

	err := tool.LintFiles(p, []string{"file1.go", filepath.Join(p.RootDirectory(), "file3.go")}, []string{})
	assert.NoError(t, err)

	runMock.AssertExpectations(t)
}

func TestClangTidy_LintFiles_Failed(t *testing.T) {
	runMock := test.Executable{}

	tool := NewClangTidy(runMock.LazyExecutable("clang-tidy"))

	p := internal.MakeProject("project", "build", nil, internal.Workflow{})

	runMock.OnRun("clang-tidy", []string{
		"-p", p.BuildDirectory(),
		filepath.Join(p.RootDirectory(), "file1.go"),
	}).
		Return(errors.New("failed"))

	err := tool.LintFiles(p, []string{"file1.go"}, []string{})
	assert.EqualError(t, err, "failed")

	runMock.AssertExpectations(t)
}

func TestClangTidy_LintFiles_CustomConfig(t *testing.T) {
	runMock := test.Executable{}

	tool := NewClangTidy(runMock.LazyExecutable("clang-tidy"))

	p := internal.MakeProject(t.TempDir(), "build", nil, internal.Workflow{})

	_, err := os.Create(filepath.Join(p.RootDirectory(), ".clang-tidy"))
	assert.NoError(t, err)

	runMock.OnRun("clang-tidy", []string{
		fmt.Sprintf("--config-file=%v", filepath.Join(p.RootDirectory(), ".clang-tidy")),
		"-p", p.BuildDirectory(),
		filepath.Join(p.RootDirectory(), "file1.go"),
	}).
		Return(nil)

	err = tool.LintFiles(p, []string{"file1.go"}, []string{})
	assert.NoError(t, err)

	runMock.AssertExpectations(t)
}

func TestClangTidy_Run(t *testing.T) {
	runMock := test.Executable{}

	tool := NewClangTidy(runMock.LazyExecutable("clang-tidy"))

	p := internal.MakeProject("project", "build", nil, internal.Workflow{})

	runMock.OnRun("clang-tidy", []string{p.RootDirectory()}).
		Return(nil)

	err := tool.RunForProject(p, []string{})
	assert.NoError(t, err)

	runMock.AssertExpectations(t)
}

func TestClangTidy_Run_Failed(t *testing.T) {
	runMock := test.Executable{}

	tool := NewClangTidy(runMock.LazyExecutable("clang-tidy"))

	p := internal.MakeProject("project", "build", nil, internal.Workflow{})

	runMock.OnRun("clang-tidy", []string{p.RootDirectory()}).
		Return(errors.New("failed"))

	err := tool.RunForProject(p, []string{})
	assert.EqualError(t, err, "failed")

	runMock.AssertExpectations(t)
}
