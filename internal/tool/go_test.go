package tool

import (
	"cdt/internal"
	"cdt/internal/test"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGo_DetectGo(t *testing.T) {
	env := test.NewEnvironment(t)
	env.OnFindExecutable("go").
		Return(env.NewExecutable("/bin/go"))

	tool := DetectGo(env)
	assert.NotNil(t, tool)
	assert.Equal(t, "go", tool.Id())
	assert.True(t, tool.IsAvailable())

	if executable := tool.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/go", executable.Path)
	}

	env.AssertExpectations(t)
}

func TestGo_DetectGo_NotFound(t *testing.T) {
	env := test.NewEnvironment(t)
	env.OnFindExecutable("go").
		Return(nil)

	tool := DetectGo(env)
	assert.NotNil(t, tool)
	assert.Equal(t, "go", tool.Id())
	assert.False(t, tool.IsAvailable())
	assert.Nil(t, tool.Executable())

	env.AssertExpectations(t)
}

func TestGo_Structure(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewGo(exec.LazyExecutable("go"))

	p := internal.ProjectInfo{Directory: "."}

	exec.OnRunOutput(
		"go", []string{"list", "-json=ImportPath,GoFiles", "./..."},
		`{"ImportPath": "target1","GoFiles":["file1.go"]}{"ImportPath": "target2","GoFiles":["file2.go", "file3.go"]}`,
	).
		Return(nil)

	structure, err := tool.Structure(context.Background(), p)
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

	exec.AssertExpectations(t)
}

func TestGo_Structure_Failed(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewGo(exec.LazyExecutable("go"))

	p := internal.ProjectInfo{Directory: "."}

	exec.OnRun("go", []string{"list", "-json=ImportPath,GoFiles", "./..."}).
		Return(errors.New("failed"))

	structure, err := tool.Structure(context.Background(), p)
	assert.EqualError(t, err, "failed")
	assert.Nil(t, structure)

	exec.AssertExpectations(t)
}

func TestGo_Structure_InvalidJson(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewGo(exec.LazyExecutable("go"))

	p := internal.ProjectInfo{Directory: "."}

	exec.OnRunOutput("go", []string{"list", "-json=ImportPath,GoFiles", "./..."}, `{]}`).
		Return(nil)

	structure, err := tool.Structure(context.Background(), p)
	assert.ErrorContains(t, err, "json decode failed:")
	assert.Nil(t, structure)

	exec.AssertExpectations(t)
}

func TestGo_BuildAll(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewGo(exec.LazyExecutable("go"))

	p := internal.ProjectInfo{Directory: "."}

	exec.OnRun("go", []string{"build"}).
		Return(nil)

	err := tool.BuildAll(context.Background(), p, []string{})
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestGo_BuildAll_Failed(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewGo(exec.LazyExecutable("go"))

	p := internal.ProjectInfo{Directory: "."}

	exec.OnRun("go", []string{"build"}).
		Return(errors.New("failed"))

	err := tool.BuildAll(context.Background(), p, []string{})
	assert.EqualError(t, err, "failed")

	exec.AssertExpectations(t)
}

func TestGo_BuildTargets(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewGo(exec.LazyExecutable("go"))

	p := internal.ProjectInfo{Directory: "."}

	exec.OnRun("go", []string{"build", "target1", "target2"}).
		Return(nil)

	err := tool.BuildTargets(context.Background(), p, []string{"target1", "target2"}, []string{})
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestGo_BuildTargets_Failed(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewGo(exec.LazyExecutable("go"))

	p := internal.ProjectInfo{Directory: "."}

	exec.OnRun("go", []string{"build", "target1", "target2"}).
		Return(errors.New("failed"))

	err := tool.BuildTargets(context.Background(), p, []string{"target1", "target2"}, []string{})
	assert.EqualError(t, err, "failed")

	exec.AssertExpectations(t)
}

func TestGo_RunTarget(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewGo(exec.LazyExecutable("go"))

	p := internal.ProjectInfo{Directory: "."}

	exec.OnRun("go", []string{"run", "target1"}).
		Return(nil)

	err := tool.RunTarget(context.Background(), p, "target1", []string{})
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestGo_RunTarget_Failed(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewGo(exec.LazyExecutable("go"))

	p := internal.ProjectInfo{Directory: "."}

	exec.OnRun("go", []string{"run", "target1"}).
		Return(errors.New("failed"))

	err := tool.RunTarget(context.Background(), p, "target1", []string{})
	assert.EqualError(t, err, "failed")

	exec.AssertExpectations(t)
}

func TestGo_TestAll(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewGo(exec.LazyExecutable("go"))

	p := internal.ProjectInfo{Directory: "."}

	exec.OnRun("go", []string{"test", "./..."}).
		Return(nil)

	err := tool.TestAll(context.Background(), p, []string{})
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestGo_TestAll_Failed(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewGo(exec.LazyExecutable("go"))

	p := internal.ProjectInfo{Directory: "."}

	exec.OnRun("go", []string{"test", "./..."}).
		Return(errors.New("failed"))

	err := tool.TestAll(context.Background(), p, []string{})
	assert.EqualError(t, err, "failed")

	exec.AssertExpectations(t)
}

func TestGo_Test(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewGo(exec.LazyExecutable("go"))

	p := internal.ProjectInfo{Directory: "."}

	exec.OnRun("go", []string{"test", "test1"}).
		Return(nil)

	err := tool.Test(context.Background(), p, "test1", []string{})
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestGo_Test_Failed(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewGo(exec.LazyExecutable("go"))

	p := internal.ProjectInfo{Directory: "."}

	exec.OnRun("go", []string{"test", "test1"}).
		Return(errors.New("failed"))

	err := tool.Test(context.Background(), p, "test1", []string{})
	assert.EqualError(t, err, "failed")

	exec.AssertExpectations(t)
}

func TestGo_FormatAll(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewGo(exec.LazyExecutable("go"))

	p := internal.ProjectInfo{Directory: "."}

	exec.OnRun("go", []string{"fmt", "./..."}).
		Return(nil)

	err := tool.FormatAll(context.Background(), p, []string{})
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestGo_FormatAll_Failed(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewGo(exec.LazyExecutable("go"))

	p := internal.ProjectInfo{Directory: "."}

	exec.OnRun("go", []string{"fmt", "./..."}).
		Return(errors.New("failed"))

	err := tool.FormatAll(context.Background(), p, []string{})
	assert.EqualError(t, err, "failed")

	exec.AssertExpectations(t)
}

func TestGo_FormatFiles(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewGo(exec.LazyExecutable("go"))

	p := internal.ProjectInfo{Directory: "."}

	exec.OnRun("go", []string{"fmt", "file1"}).
		Return(nil)

	err := tool.FormatFiles(context.Background(), p, []string{"file1"}, []string{})
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestGo_FormatFiles_Failed(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewGo(exec.LazyExecutable("go"))

	p := internal.ProjectInfo{Directory: "."}

	exec.OnRun("go", []string{"fmt", "file1"}).
		Return(errors.New("failed"))

	err := tool.FormatFiles(context.Background(), p, []string{"file1"}, []string{})
	assert.EqualError(t, err, "failed")

	exec.AssertExpectations(t)
}

func TestGo_FormatCheckAll(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewGo(exec.LazyExecutable("go"))

	p := internal.ProjectInfo{Directory: "."}

	err := tool.FormatCheckAll(context.Background(), p, []string{})
	assert.EqualError(t, err, "go fmt doesn't support check mode")

	exec.AssertExpectations(t)
}

func TestGo_FormatCheckFiles(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewGo(exec.LazyExecutable("go"))

	p := internal.ProjectInfo{Directory: "."}

	err := tool.FormatCheckFiles(context.Background(), p, []string{"file1"}, []string{})
	assert.EqualError(t, err, "go fmt doesn't support check mode")

	exec.AssertExpectations(t)
}

func TestGo_Go_LintAll(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewGo(exec.LazyExecutable("lint"))

	p := internal.ProjectInfo{Directory: "."}

	exec.OnRun("lint", []string{"vet", "./..."}).
		Return(nil)

	err := tool.LintAll(context.Background(), p, []string{})
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestGo_Go_Lint(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewGo(exec.LazyExecutable("lint"))

	p := internal.ProjectInfo{Directory: "."}

	exec.OnRun("lint", []string{"vet", "mod1"}).
		Return(nil)

	err := tool.LintFiles(context.Background(), p, []string{"mod1"}, []string{})
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}
