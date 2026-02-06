package workflow_test

import (
	"os"
	"path/filepath"
	"testing"

	"cdt/internal"
	"cdt/internal/test"
	"cdt/internal/tool"
	"cdt/internal/workflow"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPHPType_Detect_NoModFile(t *testing.T) {
	workflowType := workflow.PHP{}

	res := workflowType.Detect("dir1")

	assert.False(t, res)
}

func TestPHPType_Detect_ModFile(t *testing.T) {
	workflowType := workflow.PHP{}

	dir := t.TempDir()

	file, err := os.Create(filepath.Join(dir, "composer.json"))
	require.NoError(t, err)
	assert.NoError(t, file.Close())

	res := workflowType.Detect(dir)

	assert.True(t, res)
}

func TestPHPType_Create(t *testing.T) {
	workflowType := workflow.PHP{}

	tools := internal.Tools{
		tool.NewPHP(test.LazyExecutable("php-test")),
		tool.NewPHPUnit(test.LazyExecutable("phpunit-test")),
		tool.NewParaTest(test.LazyExecutable("paratest-test")),
		tool.NewPHPStan(test.LazyExecutable("phpstan-test")),
		tool.NewPHPCSFixer(test.LazyExecutable("php-cs-fixer-test")),
		tool.NewComposer(test.LazyExecutable("composer-test")),
	}

	project := workflowType.Create(workflow.Config{Directory: "dir1"}, tools)

	require.NotNil(t, project)
	assert.Nil(t, project.Workflow.Configurator)
	assert.Nil(t, project.Workflow.Builder)
	assert.NotNil(t, project.Workflow.Tester)
	assert.NotNil(t, project.Workflow.Linter)
	assert.NotNil(t, project.Workflow.Formatter)
	assert.NotNil(t, project.Workflow.Runner)
}

func TestPHPType_Project_RunTests_Paratest(t *testing.T) {
	workflowType := workflow.PHP{}
	paratestMock := test.NewExecutable(t)
	phpunitMock := test.NewExecutable(t)

	tools := internal.Tools{
		tool.NewPHP(test.LazyExecutable("php-test")),
		tool.NewPHPUnit(phpunitMock.LazyExecutable("phpunit-test")),
		tool.NewParaTest(paratestMock.LazyExecutable("paratest-test")),
		tool.NewPHPStan(test.LazyExecutable("phpstan-test")),
		tool.NewPHPCSFixer(test.LazyExecutable("php-cs-fixer-test")),
		tool.NewComposer(test.LazyExecutable("composer-test")),
	}

	dir := t.TempDir()

	project := workflowType.Create(workflow.Config{Directory: dir}, tools)

	require.NotNil(t, project.Workflow.Tester)
	paratestMock.OnRunAnything("paratest-test").Return(nil)

	err := project.Workflow.Tester.RunTests(
		t.Context(),
		internal.ProjectTesterOptions{ProjectInfo: project.Info},
	)
	require.NoError(t, err)

	paratestMock.AssertExpectations(t)
	phpunitMock.AssertExpectations(t)
}

func TestPHPType_Project_RunTests_PHPUnit(t *testing.T) {
	workflowType := workflow.PHP{}
	paratestMock := test.NewExecutable(t)
	phpunitMock := test.NewExecutable(t)

	tools := internal.Tools{
		tool.NewPHP(test.LazyExecutable("php-test")),
		tool.NewPHPUnit(phpunitMock.LazyExecutable("phpunit-test")),
		tool.NewParaTest(test.LazyExecutableNil),
		tool.NewPHPStan(test.LazyExecutable("phpstan-test")),
		tool.NewPHPCSFixer(test.LazyExecutable("php-cs-fixer-test")),
		tool.NewComposer(test.LazyExecutable("composer-test")),
	}

	dir := t.TempDir()

	project := workflowType.Create(workflow.Config{Directory: dir}, tools)

	require.NotNil(t, project.Workflow.Tester)
	phpunitMock.OnRunAnything("phpunit-test").Return(nil)

	err := project.Workflow.Tester.RunTests(
		t.Context(),
		internal.ProjectTesterOptions{ProjectInfo: project.Info},
	)
	require.NoError(t, err)

	paratestMock.AssertExpectations(t)
	phpunitMock.AssertExpectations(t)
}

func TestPHPType_Project_RunTests_Pattern_Paratest(t *testing.T) {
	workflowType := workflow.PHP{}
	paratestMock := test.NewExecutable(t)
	phpunitMock := test.NewExecutable(t)

	tools := internal.Tools{
		tool.NewPHP(test.LazyExecutable("php-test")),
		tool.NewPHPUnit(phpunitMock.LazyExecutable("phpunit-test")),
		tool.NewParaTest(paratestMock.LazyExecutable("paratest-test")),
		tool.NewPHPStan(test.LazyExecutable("phpstan-test")),
		tool.NewPHPCSFixer(test.LazyExecutable("php-cs-fixer-test")),
		tool.NewComposer(test.LazyExecutable("composer-test")),
	}

	dir := t.TempDir()

	project := workflowType.Create(workflow.Config{Directory: dir}, tools)

	require.NotNil(t, project.Workflow.Tester)
	paratestMock.OnRunAnything("paratest-test").Return(nil)

	err := project.Workflow.Tester.RunTests(
		t.Context(),
		internal.ProjectTesterOptions{
			ProjectInfo: project.Info,
			Pattern:     internal.StrPtr("my-test"),
		},
	)
	require.NoError(t, err)

	paratestMock.AssertExpectations(t)
	phpunitMock.AssertExpectations(t)
}

func TestPHPType_Project_RunTests_Pattern_PHPUnit(t *testing.T) {
	workflowType := workflow.PHP{}
	paratestMock := test.NewExecutable(t)
	phpunitMock := test.NewExecutable(t)

	tools := internal.Tools{
		tool.NewPHP(test.LazyExecutable("php-test")),
		tool.NewPHPUnit(phpunitMock.LazyExecutable("phpunit-test")),
		tool.NewParaTest(test.LazyExecutableNil),
		tool.NewPHPStan(test.LazyExecutable("phpstan-test")),
		tool.NewPHPCSFixer(test.LazyExecutable("php-cs-fixer-test")),
		tool.NewComposer(test.LazyExecutable("composer-test")),
	}

	dir := t.TempDir()

	project := workflowType.Create(workflow.Config{Directory: dir}, tools)

	require.NotNil(t, project.Workflow.Tester)
	phpunitMock.OnRunAnything("phpunit-test").Return(nil)

	err := project.Workflow.Tester.RunTests(
		t.Context(),
		internal.ProjectTesterOptions{
			ProjectInfo: project.Info,
			Pattern:     internal.StrPtr("my-test"),
		},
	)
	require.NoError(t, err)

	paratestMock.AssertExpectations(t)
	phpunitMock.AssertExpectations(t)
}
