package tool

import (
	"cdt/internal"
	"cdt/internal/test"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGo_DetectGo(t *testing.T) {
	environment := test.Environment{}
	environment.OnFindExecutable("go").
		Return(environment.MakeExecutable("/bin/go"))

	tool := DetectGo(&environment)
	assert.NotNil(t, tool)
	assert.Equal(t, "go", tool.Id())
	assert.True(t, tool.IsAvailable())

	if path := tool.ExecutablePath(); assert.NotNil(t, path) {
		assert.Equal(t, "/bin/go", *path)
	}

	environment.AssertExpectations(t)
}

func TestGo_DetectGo_NotFound(t *testing.T) {
	environment := test.Environment{}
	environment.OnFindExecutable("go").
		Return(nil)

	tool := DetectGo(&environment)
	assert.NotNil(t, tool)
	assert.Equal(t, "go", tool.Id())
	assert.False(t, tool.IsAvailable())
	assert.Nil(t, tool.ExecutablePath())

	environment.AssertExpectations(t)
}

func TestGo_Structure(t *testing.T) {
	runMock := test.Executable{}

	tool := NewGo(runMock.LazyExecutable("go"))

	p := internal.MakeProject("project", "", tool, internal.Workflow{})

	runMock.OnRunOutput(
		"go", []string{"list", "-json=ImportPath,GoFiles", "./..."},
		`{"ImportPath": "target1","GoFiles":["file1.go"]}{"ImportPath": "target2","GoFiles":["file2.go", "file3.go"]}`,
	).
		Return(nil)

	structure, err := tool.Structure(p)
	assert.NoError(t, err)
	if assert.NotNil(t, structure) {
		assert.Equal(t,
			internal.ProjectStructure{
				Targets: map[string]internal.ProjectTarget{
					"target1": {
						Files: []string{"file1.go"},
					},
					"target2": {
						Files: []string{"file2.go", "file3.go"},
					},
				},
			},
			*structure,
		)
	}

	runMock.AssertExpectations(t)
}

func TestGo_Structure_Failed(t *testing.T) {
	runMock := test.Executable{}

	tool := NewGo(runMock.LazyExecutable("go"))

	p := internal.MakeProject("project", "", tool, internal.Workflow{})

	runMock.OnRun("go", []string{"list", "-json=ImportPath,GoFiles", "./..."}).
		Return(errors.New("failed"))

	structure, err := tool.Structure(p)
	assert.EqualError(t, err, "failed")
	assert.Nil(t, structure)

	runMock.AssertExpectations(t)
}

func TestGo_Structure_InvalidJson(t *testing.T) {
	runMock := test.Executable{}

	tool := NewGo(runMock.LazyExecutable("go"))

	p := internal.MakeProject("project", "", tool, internal.Workflow{})

	runMock.OnRunOutput("go", []string{"list", "-json=ImportPath,GoFiles", "./..."}, `{]}`).
		Return(nil)

	structure, err := tool.Structure(p)
	assert.ErrorContains(t, err, "json decode failed:")
	assert.Nil(t, structure)

	runMock.AssertExpectations(t)
}

func TestGo_BuildAll(t *testing.T) {
	runMock := test.Executable{}

	tool := NewGo(runMock.LazyExecutable("go"))

	p := internal.MakeProject("project", "", tool, internal.Workflow{})

	runMock.OnRun("go", []string{"build"}).
		Return(nil)

	err := tool.BuildAll(p, []string{})
	assert.NoError(t, err)

	runMock.AssertExpectations(t)
}

func TestGo_BuildAll_Failed(t *testing.T) {
	runMock := test.Executable{}

	tool := NewGo(runMock.LazyExecutable("go"))

	p := internal.MakeProject("project", "", tool, internal.Workflow{})

	runMock.OnRun("go", []string{"build"}).
		Return(errors.New("failed"))

	err := tool.BuildAll(p, []string{})
	assert.EqualError(t, err, "failed")

	runMock.AssertExpectations(t)
}

func TestGo_BuildTargets(t *testing.T) {
	runMock := test.Executable{}

	tool := NewGo(runMock.LazyExecutable("go"))

	p := internal.MakeProject("project", "", tool, internal.Workflow{})

	runMock.OnRun("go", []string{"build", "target1", "target2"}).
		Return(nil)

	err := tool.BuildTargets(p, []string{"target1", "target2"}, []string{})
	assert.NoError(t, err)

	runMock.AssertExpectations(t)
}

func TestGo_BuildTargets_Failed(t *testing.T) {
	runMock := test.Executable{}

	tool := NewGo(runMock.LazyExecutable("go"))

	p := internal.MakeProject("project", "", tool, internal.Workflow{})

	runMock.OnRun("go", []string{"build", "target1", "target2"}).
		Return(errors.New("failed"))

	err := tool.BuildTargets(p, []string{"target1", "target2"}, []string{})
	assert.EqualError(t, err, "failed")

	runMock.AssertExpectations(t)
}

func TestGo_RunTarget(t *testing.T) {
	runMock := test.Executable{}

	tool := NewGo(runMock.LazyExecutable("go"))

	p := internal.MakeProject("project", "", tool, internal.Workflow{})

	runMock.OnRun("go", []string{"run", "target1"}).
		Return(nil)

	err := tool.RunTarget(p, "target1", []string{})
	assert.NoError(t, err)

	runMock.AssertExpectations(t)
}

func TestGo_RunTarget_Failed(t *testing.T) {
	runMock := test.Executable{}

	tool := NewGo(runMock.LazyExecutable("go"))

	p := internal.MakeProject("project", "", tool, internal.Workflow{})

	runMock.OnRun("go", []string{"run", "target1"}).
		Return(errors.New("failed"))

	err := tool.RunTarget(p, "target1", []string{})
	assert.EqualError(t, err, "failed")

	runMock.AssertExpectations(t)
}

func TestGo_TestAll(t *testing.T) {
	runMock := test.Executable{}

	tool := NewGo(runMock.LazyExecutable("go"))

	p := internal.MakeProject("project", "", tool, internal.Workflow{})

	runMock.OnRun("go", []string{"test", "./..."}).
		Return(nil)

	err := tool.TestAll(p, []string{})
	assert.NoError(t, err)

	runMock.AssertExpectations(t)
}

func TestGo_TestAll_Failed(t *testing.T) {
	runMock := test.Executable{}

	tool := NewGo(runMock.LazyExecutable("go"))

	p := internal.MakeProject("project", "", tool, internal.Workflow{})

	runMock.OnRun("go", []string{"test", "./..."}).
		Return(errors.New("failed"))

	err := tool.TestAll(p, []string{})
	assert.EqualError(t, err, "failed")

	runMock.AssertExpectations(t)
}

func TestGo_Test(t *testing.T) {
	runMock := test.Executable{}

	tool := NewGo(runMock.LazyExecutable("go"))

	p := internal.MakeProject("project", "", tool, internal.Workflow{})

	runMock.OnRun("go", []string{"test", "test1"}).
		Return(nil)

	err := tool.Test(p, "test1", []string{})
	assert.NoError(t, err)

	runMock.AssertExpectations(t)
}

func TestGo_Test_Failed(t *testing.T) {
	runMock := test.Executable{}

	tool := NewGo(runMock.LazyExecutable("go"))

	p := internal.MakeProject("project", "", tool, internal.Workflow{})

	runMock.OnRun("go", []string{"test", "test1"}).
		Return(errors.New("failed"))

	err := tool.Test(p, "test1", []string{})
	assert.EqualError(t, err, "failed")

	runMock.AssertExpectations(t)
}

func TestGo_FormatAll(t *testing.T) {
	runMock := test.Executable{}

	tool := NewGo(runMock.LazyExecutable("go"))

	p := internal.MakeProject("project", "", tool, internal.Workflow{})

	runMock.OnRun("go", []string{"fmt", "./..."}).
		Return(nil)

	err := tool.FormatAll(p, []string{})
	assert.NoError(t, err)

	runMock.AssertExpectations(t)
}

func TestGo_FormatAll_Failed(t *testing.T) {
	runMock := test.Executable{}

	tool := NewGo(runMock.LazyExecutable("go"))

	p := internal.MakeProject("project", "", tool, internal.Workflow{})

	runMock.OnRun("go", []string{"fmt", "./..."}).
		Return(errors.New("failed"))

	err := tool.FormatAll(p, []string{})
	assert.EqualError(t, err, "failed")

	runMock.AssertExpectations(t)
}

func TestGo_FormatFiles(t *testing.T) {
	runMock := test.Executable{}

	tool := NewGo(runMock.LazyExecutable("go"))

	p := internal.MakeProject("project", "", tool, internal.Workflow{})

	runMock.OnRun("go", []string{"fmt", "file1"}).
		Return(nil)

	err := tool.FormatFiles(p, []string{"file1"}, []string{})
	assert.NoError(t, err)

	runMock.AssertExpectations(t)
}

func TestGo_FormatFiles_Failed(t *testing.T) {
	runMock := test.Executable{}

	tool := NewGo(runMock.LazyExecutable("go"))

	p := internal.MakeProject("project", "", tool, internal.Workflow{})

	runMock.OnRun("go", []string{"fmt", "file1"}).
		Return(errors.New("failed"))

	err := tool.FormatFiles(p, []string{"file1"}, []string{})
	assert.EqualError(t, err, "failed")

	runMock.AssertExpectations(t)
}

func TestGo_FormatCheckAll(t *testing.T) {
	runMock := test.Executable{}

	tool := NewGo(runMock.LazyExecutable("go"))

	p := internal.MakeProject("project", "", tool, internal.Workflow{})

	err := tool.FormatCheckAll(p, []string{})
	assert.EqualError(t, err, "go fmt doesn't support check mode")

	runMock.AssertExpectations(t)
}

func TestGo_FormatCheckFiles(t *testing.T) {
	runMock := test.Executable{}

	tool := NewGo(runMock.LazyExecutable("go"))

	p := internal.MakeProject("project", "", tool, internal.Workflow{})

	err := tool.FormatCheckFiles(p, []string{"file1"}, []string{})
	assert.EqualError(t, err, "go fmt doesn't support check mode")

	runMock.AssertExpectations(t)
}
