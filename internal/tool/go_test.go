package tool

import (
	"cdt/internal"
	"cdt/internal/test"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestGo_DetectGo(t *testing.T) {
	environment := test.Environment{}
	environment.On("FindExecutable", "go").Return(environment.MakeExecutable("/bin/go"))

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
	environment.On("FindExecutable", "go").Return(nil)

	tool := DetectGo(&environment)
	assert.NotNil(t, tool)
	assert.Equal(t, "go", tool.Id())
	assert.False(t, tool.IsAvailable())
	assert.Nil(t, tool.ExecutablePath())

	environment.AssertExpectations(t)
}

func TestGo_Structure(t *testing.T) {
	environment := test.Environment{}

	tool := NewGo(environment.DetectExecutable("go"))

	p := internal.MakeProject("project", "", tool, internal.Workflow{})

	environment.On("RunExecutable", mock.Anything, mock.Anything, "go", []string{"list", "-json=ImportPath,GoFiles", "./..."}).
		Return(nil).
		Run(func(args mock.Arguments) {
			ctx := args.Get(1).(internal.RunOptions)
			_, _ = ctx.Output.Write([]byte(
				`{"ImportPath": "target1","GoFiles":["file1.go"]}{"ImportPath": "target2","GoFiles":["file2.go", "file3.go"]}`,
			))
		})

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

	environment.AssertExpectations(t)
}

func TestGo_Structure_Failed(t *testing.T) {
	environment := test.Environment{}

	tool := NewGo(environment.DetectExecutable("go"))

	p := internal.MakeProject("project", "", tool, internal.Workflow{})

	environment.On("RunExecutable", mock.Anything, mock.Anything, "go", []string{"list", "-json=ImportPath,GoFiles", "./..."}).
		Return(errors.New("failed"))

	structure, err := tool.Structure(p)
	assert.EqualError(t, err, "failed")
	assert.Nil(t, structure)

	environment.AssertExpectations(t)
}

func TestGo_Structure_InvalidJson(t *testing.T) {
	environment := test.Environment{}

	tool := NewGo(environment.DetectExecutable("go"))

	p := internal.MakeProject("project", "", tool, internal.Workflow{})

	environment.On("RunExecutable", mock.Anything, mock.Anything, "go", []string{"list", "-json=ImportPath,GoFiles", "./..."}).
		Return(nil).
		Run(func(args mock.Arguments) {
			ctx := args.Get(1).(internal.RunOptions)
			_, _ = ctx.Output.Write([]byte(
				`{]}`,
			))
		})

	structure, err := tool.Structure(p)
	assert.ErrorContains(t, err, "json decode failed:")
	assert.Nil(t, structure)

	environment.AssertExpectations(t)
}

func TestGo_BuildAll(t *testing.T) {
	environment := test.Environment{}

	tool := NewGo(environment.DetectExecutable("go"))

	p := internal.MakeProject("project", "", tool, internal.Workflow{})

	environment.OnRunSuccess(p, "go", []string{"build"})

	err := tool.BuildAll(p, []string{})
	assert.NoError(t, err)

	environment.AssertExpectations(t)
}

func TestGo_BuildAll_Failed(t *testing.T) {
	environment := test.Environment{}

	tool := NewGo(environment.DetectExecutable("go"))

	p := internal.MakeProject("project", "", tool, internal.Workflow{})

	environment.OnRunError(p, "go", []string{"build"}, errors.New("failed"))

	err := tool.BuildAll(p, []string{})
	assert.EqualError(t, err, "failed")

	environment.AssertExpectations(t)
}

func TestGo_BuildTargets(t *testing.T) {
	environment := test.Environment{}

	tool := NewGo(environment.DetectExecutable("go"))

	p := internal.MakeProject("project", "", tool, internal.Workflow{})

	environment.OnRunSuccess(p, "go", []string{"build", "target1", "target2"})

	err := tool.BuildTargets(p, []string{"target1", "target2"}, []string{})
	assert.NoError(t, err)

	environment.AssertExpectations(t)
}

func TestGo_BuildTargets_Failed(t *testing.T) {
	environment := test.Environment{}

	tool := NewGo(environment.DetectExecutable("go"))

	p := internal.MakeProject("project", "", tool, internal.Workflow{})

	environment.OnRunError(p, "go", []string{"build", "target1", "target2"}, errors.New("failed"))

	err := tool.BuildTargets(p, []string{"target1", "target2"}, []string{})
	assert.EqualError(t, err, "failed")

	environment.AssertExpectations(t)
}

func TestGo_RunTarget(t *testing.T) {
	environment := test.Environment{}

	tool := NewGo(environment.DetectExecutable("go"))

	p := internal.MakeProject("project", "", tool, internal.Workflow{})

	environment.OnRunSuccess(p, "go", []string{"run", "target1"})

	err := tool.RunTarget(p, "target1", []string{})
	assert.NoError(t, err)

	environment.AssertExpectations(t)
}

func TestGo_RunTarget_Failed(t *testing.T) {
	environment := test.Environment{}

	tool := NewGo(environment.DetectExecutable("go"))

	p := internal.MakeProject("project", "", tool, internal.Workflow{})

	environment.OnRunError(p, "go", []string{"run", "target1"}, errors.New("failed"))

	err := tool.RunTarget(p, "target1", []string{})
	assert.EqualError(t, err, "failed")

	environment.AssertExpectations(t)
}

func TestGo_TestAll(t *testing.T) {
	environment := test.Environment{}

	tool := NewGo(environment.DetectExecutable("go"))

	p := internal.MakeProject("project", "", tool, internal.Workflow{})

	environment.OnRunSuccess(p, "go", []string{"test", "./..."})

	err := tool.TestAll(p, []string{})
	assert.NoError(t, err)

	environment.AssertExpectations(t)
}

func TestGo_TestAll_Failed(t *testing.T) {
	environment := test.Environment{}

	tool := NewGo(environment.DetectExecutable("go"))

	p := internal.MakeProject("project", "", tool, internal.Workflow{})

	environment.OnRunError(p, "go", []string{"test", "./..."}, errors.New("failed"))

	err := tool.TestAll(p, []string{})
	assert.EqualError(t, err, "failed")

	environment.AssertExpectations(t)
}

func TestGo_Test(t *testing.T) {
	environment := test.Environment{}

	tool := NewGo(environment.DetectExecutable("go"))

	p := internal.MakeProject("project", "", tool, internal.Workflow{})

	environment.OnRunSuccess(p, "go", []string{"test", "test1"})

	err := tool.Test(p, "test1", []string{})
	assert.NoError(t, err)

	environment.AssertExpectations(t)
}

func TestGo_Test_Failed(t *testing.T) {
	environment := test.Environment{}

	tool := NewGo(environment.DetectExecutable("go"))

	p := internal.MakeProject("project", "", tool, internal.Workflow{})

	environment.OnRunError(p, "go", []string{"test", "test1"}, errors.New("failed"))

	err := tool.Test(p, "test1", []string{})
	assert.EqualError(t, err, "failed")

	environment.AssertExpectations(t)
}

func TestGo_FormatAll(t *testing.T) {
	environment := test.Environment{}

	tool := NewGo(environment.DetectExecutable("go"))

	p := internal.MakeProject("project", "", tool, internal.Workflow{})

	environment.OnRunSuccess(p, "go", []string{"fmt", "./..."})

	err := tool.FormatAll(p, []string{})
	assert.NoError(t, err)

	environment.AssertExpectations(t)
}

func TestGo_FormatAll_Failed(t *testing.T) {
	environment := test.Environment{}

	tool := NewGo(environment.DetectExecutable("go"))

	p := internal.MakeProject("project", "", tool, internal.Workflow{})

	environment.OnRunError(p, "go", []string{"fmt", "./..."}, errors.New("failed"))

	err := tool.FormatAll(p, []string{})
	assert.EqualError(t, err, "failed")

	environment.AssertExpectations(t)
}

func TestGo_FormatFiles(t *testing.T) {
	environment := test.Environment{}

	tool := NewGo(environment.DetectExecutable("go"))

	p := internal.MakeProject("project", "", tool, internal.Workflow{})

	environment.OnRunSuccess(p, "go", []string{"fmt", "file1"})

	err := tool.FormatFiles(p, []string{"file1"}, []string{})
	assert.NoError(t, err)

	environment.AssertExpectations(t)
}

func TestGo_FormatFiles_Failed(t *testing.T) {
	environment := test.Environment{}

	tool := NewGo(environment.DetectExecutable("go"))

	p := internal.MakeProject("project", "", tool, internal.Workflow{})

	environment.OnRunError(p, "go", []string{"fmt", "file1"}, errors.New("failed"))

	err := tool.FormatFiles(p, []string{"file1"}, []string{})
	assert.EqualError(t, err, "failed")

	environment.AssertExpectations(t)
}

func TestGo_FormatCheckAll(t *testing.T) {
	environment := test.Environment{}

	tool := NewGo(environment.DetectExecutable("go"))

	p := internal.MakeProject("project", "", tool, internal.Workflow{})

	err := tool.FormatCheckAll(p, []string{})
	assert.EqualError(t, err, "go fmt doesn't support check mode")

	environment.AssertExpectations(t)
}

func TestGo_FormatCheckFiles(t *testing.T) {
	environment := test.Environment{}

	tool := NewGo(environment.DetectExecutable("go"))

	p := internal.MakeProject("project", "", tool, internal.Workflow{})

	err := tool.FormatCheckFiles(p, []string{"file1"}, []string{})
	assert.EqualError(t, err, "go fmt doesn't support check mode")

	environment.AssertExpectations(t)
}
