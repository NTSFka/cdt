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
	environment.On("FindExecutable", "clang-tidy").Return(environment.MakeExecutable("/bin/clang-tidy"))

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
	environment.On("FindExecutable", "clang-tidy").Return(nil)
	environment.On("FindExecutable", fmt.Sprintf("clang-tidy-%v", clangTidyMaxVersion)).
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
	environment.On("FindExecutable", "clang-tidy-20").Return(environment.MakeExecutable("/bin/clang-tidy-20"))

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
	environment.On("FindExecutable", "clang-tidy").Return(nil)

	for version := clangTidyMaxVersion; version >= clangTidyMinVersion; version-- {
		environment.On("FindExecutable", fmt.Sprintf("clang-tidy-%v", version)).Return(nil)
	}

	tool := DetectClangTidy(&environment, nil)
	assert.NotNil(t, tool)
	assert.Equal(t, "clang-tidy", tool.Id())
	assert.False(t, tool.IsAvailable())
	assert.Nil(t, tool.ExecutablePath())

	environment.AssertExpectations(t)
}

func TestClangTidy_LintAll(t *testing.T) {
	environment := test.Environment{}

	tool := NewClangTidy(environment.MakeExecutable("clang-tidy"), nil)

	p := internal.MakeProject("project", "build", &internal.FixedProjectStructureProvider{
		ProjectStructure: internal.ProjectStructure{
			Targets: map[string]internal.ProjectTarget{
				"target1": {
					Files: []string{"file1.go", "file2.go"},
				},
			},
		},
	}, internal.Workflow{})

	environment.OnRunSuccess(p, "clang-tidy", []string{
		"-p", p.BuildDirectory(),
		filepath.Join(p.RootDirectory(), "file1.go"),
		filepath.Join(p.RootDirectory(), "file2.go"),
	})

	err := tool.LintAll(p, []string{})
	assert.NoError(t, err)

	environment.AssertExpectations(t)
}

func TestClangTidy_LintAll_Failed(t *testing.T) {
	environment := test.Environment{}

	tool := NewClangTidy(environment.MakeExecutable("clang-tidy"), nil)

	p := internal.MakeProject("project", "build", &internal.FixedProjectStructureProvider{
		ProjectStructure: internal.ProjectStructure{
			Targets: map[string]internal.ProjectTarget{
				"target1": {
					Files: []string{"file1.go", "file2.go"},
				},
			},
		},
	}, internal.Workflow{})

	environment.OnRunError(p, "clang-tidy", []string{
		"-p", p.BuildDirectory(),
		filepath.Join(p.RootDirectory(), "file1.go"),
		filepath.Join(p.RootDirectory(), "file2.go"),
	}, errors.New("failed"))

	err := tool.LintAll(p, []string{})
	assert.EqualError(t, err, "failed")

	environment.AssertExpectations(t)
}

func TestClangTidy_LintAll_CustomConfig(t *testing.T) {
	environment := test.Environment{}

	tool := NewClangTidy(environment.MakeExecutable("clang-tidy"), nil)

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

	environment.OnRunSuccess(p, "clang-tidy", []string{
		fmt.Sprintf("--config-file=%v", filepath.Join(p.RootDirectory(), ".clang-tidy")),
		"-p", p.BuildDirectory(),
		filepath.Join(p.RootDirectory(), "file1.go"),
		filepath.Join(p.RootDirectory(), "file2.go"),
	})

	err = tool.LintAll(p, []string{})
	assert.NoError(t, err)

	environment.AssertExpectations(t)
}

func TestClangTidy_LintFiles(t *testing.T) {
	environment := test.Environment{}

	tool := NewClangTidy(environment.MakeExecutable("clang-tidy"), nil)

	p := internal.MakeProject(t.TempDir(), "build", nil, internal.Workflow{})

	environment.OnRunSuccess(p, "clang-tidy", []string{
		"-p", p.BuildDirectory(),
		filepath.Join(p.RootDirectory(), "file1.go"),
		filepath.Join(p.RootDirectory(), "file3.go"),
	})

	err := tool.LintFiles(p, []string{"file1.go", filepath.Join(p.RootDirectory(), "file3.go")}, []string{})
	assert.NoError(t, err)

	environment.AssertExpectations(t)
}

func TestClangTidy_LintFiles_Failed(t *testing.T) {
	environment := test.Environment{}

	tool := NewClangTidy(environment.MakeExecutable("clang-tidy"), nil)

	p := internal.MakeProject("project", "build", nil, internal.Workflow{})

	environment.OnRunError(p, "clang-tidy", []string{
		"-p", p.BuildDirectory(),
		filepath.Join(p.RootDirectory(), "file1.go"),
	}, errors.New("failed"))

	err := tool.LintFiles(p, []string{"file1.go"}, []string{})
	assert.EqualError(t, err, "failed")

	environment.AssertExpectations(t)
}

func TestClangTidy_LintFiles_CustomConfig(t *testing.T) {
	environment := test.Environment{}

	tool := NewClangTidy(environment.MakeExecutable("clang-tidy"), nil)

	p := internal.MakeProject(t.TempDir(), "build", nil, internal.Workflow{})

	_, err := os.Create(filepath.Join(p.RootDirectory(), ".clang-tidy"))
	assert.NoError(t, err)

	environment.OnRunSuccess(p, "clang-tidy", []string{
		fmt.Sprintf("--config-file=%v", filepath.Join(p.RootDirectory(), ".clang-tidy")),
		"-p", p.BuildDirectory(),
		filepath.Join(p.RootDirectory(), "file1.go"),
	})

	err = tool.LintFiles(p, []string{"file1.go"}, []string{})
	assert.NoError(t, err)

	environment.AssertExpectations(t)
}

func TestClangTidy_Run(t *testing.T) {
	environment := test.Environment{}

	tool := NewClangTidy(environment.MakeExecutable("clang-tidy"), nil)

	p := internal.MakeProject("project", "build", nil, internal.Workflow{})

	environment.OnRunSuccess(p, "clang-tidy", []string{p.RootDirectory()})

	err := tool.Run(p, []string{})
	assert.NoError(t, err)

	environment.AssertExpectations(t)
}

func TestClangTidy_Run_Failed(t *testing.T) {
	environment := test.Environment{}

	tool := NewClangTidy(environment.MakeExecutable("clang-tidy"), nil)

	p := internal.MakeProject("project", "build", nil, internal.Workflow{})

	environment.OnRunError(p, "clang-tidy", []string{p.RootDirectory()}, errors.New("failed"))

	err := tool.Run(p, []string{})
	assert.EqualError(t, err, "failed")

	environment.AssertExpectations(t)
}
