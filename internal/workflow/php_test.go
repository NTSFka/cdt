package workflow

import (
	"cdt/internal"
	"cdt/internal/test"
	"cdt/internal/tool"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPHPType_Detect_NoModFile(t *testing.T) {
	workflowType := PHP{}

	res := workflowType.Detect("dir1")

	assert.False(t, res)
}

func TestPHPType_Detect_ModFile(t *testing.T) {
	workflowType := PHP{}

	dir := t.TempDir()

	_, err := os.Create(filepath.Join(dir, "composer.json"))
	require.NoError(t, err)

	res := workflowType.Detect(dir)

	assert.True(t, res)
}

func TestPHPType_Create(t *testing.T) {
	workflowType := PHP{}

	tools := internal.Tools{
		tool.NewPHP(func() *internal.Executable { return &internal.Executable{Path: "php-test"} }),
		tool.NewPHPUnit(func() *internal.Executable { return &internal.Executable{Path: "phpunit-test"} }),
		tool.NewParaTest(func() *internal.Executable { return &internal.Executable{Path: "paratest-test"} }),
		tool.NewPHPStan(func() *internal.Executable { return &internal.Executable{Path: "phpstan-test"} }),
		tool.NewPHPCSFixer(func() *internal.Executable { return &internal.Executable{Path: "php-cs-fixer-test"} }),
		tool.NewComposer(func() *internal.Executable { return &internal.Executable{Path: "composer-test"} }),
	}

	project := workflowType.Create(Config{Directory: "dir1"}, tools)

	require.NotNil(t, project)
	assert.Nil(t, project.Workflow.Configurator)
	assert.Nil(t, project.Workflow.Builder)
	assert.NotNil(t, project.Workflow.Tester)
	assert.NotNil(t, project.Workflow.Linter)
	assert.NotNil(t, project.Workflow.Formatter)
	assert.NotNil(t, project.Workflow.Runner)
}

func TestPHPType_Project_TestAll_Paratest(t *testing.T) {
	workflowType := PHP{}
	paratestMock := test.NewExecutable(t)
	phpunitMock := test.NewExecutable(t)

	tools := internal.Tools{
		tool.NewPHP(func() *internal.Executable { return &internal.Executable{Path: "php-test"} }),
		tool.NewPHPUnit(phpunitMock.LazyExecutable("phpunit-test")),
		tool.NewParaTest(paratestMock.LazyExecutable("paratest-test")),
		tool.NewPHPStan(func() *internal.Executable { return &internal.Executable{Path: "phpstan-test"} }),
		tool.NewPHPCSFixer(func() *internal.Executable { return &internal.Executable{Path: "php-cs-fixer-test"} }),
		tool.NewComposer(func() *internal.Executable { return &internal.Executable{Path: "composer-test"} }),
	}

	dir := t.TempDir()

	p := workflowType.Create(Config{Directory: dir}, tools)

	require.NotNil(t, p.Workflow.Tester)
	paratestMock.OnRunAnything("paratest-test").Return(nil)

	err := p.Workflow.Tester.TestAll(t.Context(), internal.ProjectTesterOptions{ProjectInfo: p.Info})
	require.NoError(t, err)

	paratestMock.AssertExpectations(t)
	phpunitMock.AssertExpectations(t)
}

func TestPHPType_Project_TestAll_PHPUnit(t *testing.T) {
	workflowType := PHP{}
	paratestMock := test.NewExecutable(t)
	phpunitMock := test.NewExecutable(t)

	tools := internal.Tools{
		tool.NewPHP(func() *internal.Executable { return &internal.Executable{Path: "php-test"} }),
		tool.NewPHPUnit(phpunitMock.LazyExecutable("phpunit-test")),
		tool.NewParaTest(func() *internal.Executable { return nil }),
		tool.NewPHPStan(func() *internal.Executable { return &internal.Executable{Path: "phpstan-test"} }),
		tool.NewPHPCSFixer(func() *internal.Executable { return &internal.Executable{Path: "php-cs-fixer-test"} }),
		tool.NewComposer(func() *internal.Executable { return &internal.Executable{Path: "composer-test"} }),
	}

	dir := t.TempDir()

	p := workflowType.Create(Config{Directory: dir}, tools)

	require.NotNil(t, p.Workflow.Tester)
	phpunitMock.OnRunAnything("phpunit-test").Return(nil)

	err := p.Workflow.Tester.TestAll(t.Context(), internal.ProjectTesterOptions{ProjectInfo: p.Info})
	require.NoError(t, err)

	paratestMock.AssertExpectations(t)
	phpunitMock.AssertExpectations(t)
}

func TestPHPType_Project_Test_Paratest(t *testing.T) {
	workflowType := PHP{}
	paratestMock := test.NewExecutable(t)
	phpunitMock := test.NewExecutable(t)

	tools := internal.Tools{
		tool.NewPHP(func() *internal.Executable { return &internal.Executable{Path: "php-test"} }),
		tool.NewPHPUnit(phpunitMock.LazyExecutable("phpunit-test")),
		tool.NewParaTest(paratestMock.LazyExecutable("paratest-test")),
		tool.NewPHPStan(func() *internal.Executable { return &internal.Executable{Path: "phpstan-test"} }),
		tool.NewPHPCSFixer(func() *internal.Executable { return &internal.Executable{Path: "php-cs-fixer-test"} }),
		tool.NewComposer(func() *internal.Executable { return &internal.Executable{Path: "composer-test"} }),
	}

	dir := t.TempDir()

	p := workflowType.Create(Config{Directory: dir}, tools)

	require.NotNil(t, p.Workflow.Tester)
	paratestMock.OnRunAnything("paratest-test").Return(nil)

	err := p.Workflow.Tester.TestPattern(t.Context(), internal.ProjectTesterOptions{ProjectInfo: p.Info}, "my-test")
	require.NoError(t, err)

	paratestMock.AssertExpectations(t)
	phpunitMock.AssertExpectations(t)
}

func TestPHPType_Project_Test_PHPUnit(t *testing.T) {
	workflowType := PHP{}
	paratestMock := test.NewExecutable(t)
	phpunitMock := test.NewExecutable(t)

	tools := internal.Tools{
		tool.NewPHP(func() *internal.Executable { return &internal.Executable{Path: "php-test"} }),
		tool.NewPHPUnit(phpunitMock.LazyExecutable("phpunit-test")),
		tool.NewParaTest(func() *internal.Executable { return nil }),
		tool.NewPHPStan(func() *internal.Executable { return &internal.Executable{Path: "phpstan-test"} }),
		tool.NewPHPCSFixer(func() *internal.Executable { return &internal.Executable{Path: "php-cs-fixer-test"} }),
		tool.NewComposer(func() *internal.Executable { return &internal.Executable{Path: "composer-test"} }),
	}

	dir := t.TempDir()

	p := workflowType.Create(Config{Directory: dir}, tools)

	require.NotNil(t, p.Workflow.Tester)
	phpunitMock.OnRunAnything("phpunit-test").Return(nil)

	err := p.Workflow.Tester.TestPattern(t.Context(), internal.ProjectTesterOptions{ProjectInfo: p.Info}, "my-test")
	require.NoError(t, err)

	paratestMock.AssertExpectations(t)
	phpunitMock.AssertExpectations(t)
}
