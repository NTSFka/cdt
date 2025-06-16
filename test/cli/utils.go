package cli

import (
	"cdt/internal"
	. "cdt/pkg"
	"github.com/stretchr/testify/mock"
	"testing"
)

type testProjectConfigurator struct {
	mock.Mock
}

func (t *testProjectConfigurator) Configure(project internal.Project, args []string) error {
	return t.Called(project, args).Error(0)
}

type testProjectBuilder struct {
	mock.Mock
}

func (t *testProjectBuilder) BuildAll(project internal.Project, args []string) error {
	return t.Called(project, args).Error(0)
}

func (t *testProjectBuilder) BuildTargets(project internal.Project, targets []string, args []string) error {
	return t.Called(project, targets, args).Error(0)
}

type testProjectFormatter struct {
	mock.Mock
}

func (t *testProjectFormatter) FormatAll(project internal.Project, args []string) error {
	return t.Called(project, args).Error(0)
}

func (t *testProjectFormatter) FormatFiles(project internal.Project, filenames []string, args []string) error {
	return t.Called(project, filenames, args).Error(0)
}

func (t *testProjectFormatter) FormatCheckAll(project internal.Project, args []string) error {
	return t.Called(project, args).Error(0)
}

func (t *testProjectFormatter) FormatCheckFiles(project internal.Project, filenames []string, args []string) error {
	return t.Called(project, filenames, args).Error(0)
}

type testProjectLinter struct {
	mock.Mock
}

func (t *testProjectLinter) LintAll(project internal.Project, args []string) error {
	return t.Called(project, args).Error(0)
}

func (t *testProjectLinter) LintFiles(project internal.Project, filenames []string, args []string) error {
	return t.Called(project, filenames, args).Error(0)
}

type testProjectTester struct {
	mock.Mock
}

func (t *testProjectTester) TestAll(project internal.Project, args []string) error {
	return t.Called(project, args).Error(0)
}

func (t *testProjectTester) Test(project internal.Project, pattern string, args []string) error {
	return t.Called(project, pattern, args).Error(0)
}

type testProjectRunner struct {
	mock.Mock
}

func (t *testProjectRunner) RunTarget(project internal.Project, target string, args []string) error {
	return t.Called(project, target, args).Error(0)
}

type testStructureProvider struct {
	mock.Mock
}

func (t *testStructureProvider) Structure(project internal.Project) (*internal.ProjectStructure, error) {
	called := t.Called(project)
	return called.Get(0).(*internal.ProjectStructure), called.Error(1)
}

func runMain(contextBuilder func(config internal.Config) internal.Context, args ...string) error {
	return RunMain(contextBuilder, append([]string{"cdt"}, args...))
}

// Run main function and return obtained configuration
func runMainGetConfig(args ...string) (config internal.Config) {
	_ = runMain(func(cfg internal.Config) internal.Context {
		config = cfg
		return internal.Context{}
	}, args...)

	return
}

func runMainWithProject(project internal.Project, args ...string) error {
	return runMain(func(config internal.Config) internal.Context {
		return internal.Context{
			Config:  internal.Config{},
			Project: project,
			Tools:   nil,
		}
	}, args...)
}

func runMainWithWorkflow(workflow internal.Workflow, args ...string) error {
	return runMain(func(config internal.Config) internal.Context {
		return internal.Context{
			Config:  internal.Config{},
			Project: internal.Project{Workflow: workflow},
			Tools:   nil,
		}
	}, args...)
}

func runMainWithTools(tools internal.Tools, args ...string) error {
	return runMain(func(config internal.Config) internal.Context {
		return internal.Context{
			Config:  internal.Config{},
			Project: internal.Project{},
			Tools:   tools,
		}
	}, args...)
}

// Check if a tool exists in the current environment
func checkTool(t *testing.T, toolName string) {
	if executable := internal.FindExecutable(toolName); executable == nil {
		t.Skipf("unable to find tool: %v", toolName)
	}
}
