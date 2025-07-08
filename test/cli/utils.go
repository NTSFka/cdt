package cli

import (
	"cdt/internal"
	"cdt/pkg"
	"errors"
	"github.com/stretchr/testify/mock"
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

func runMain(contextBuilder func(config internal.Config) (*internal.Context, error), args ...string) error {
	return pkg.RunMain(contextBuilder, append([]string{"cdt"}, args...))
}

// Run main function and return obtained configuration
func runMainGetConfig(args ...string) (config internal.Config) {
	_ = runMain(func(cfg internal.Config) (*internal.Context, error) {
		config = cfg
		return nil, errors.New("not used")
	}, args...)

	return
}

func runMainWithEnvironment(environment internal.Environment, args ...string) error {
	return runMain(func(config internal.Config) (*internal.Context, error) {
		return &internal.Context{
			Config:      internal.Config{},
			Project:     internal.Project{},
			Tools:       nil,
			Environment: environment,
		}, nil
	}, args...)
}

func runMainWithProject(project internal.Project, args ...string) error {
	return runMain(func(config internal.Config) (*internal.Context, error) {
		return &internal.Context{
			Config:  internal.Config{},
			Project: project,
			Tools:   nil,
		}, nil
	}, args...)
}

func runMainWithWorkflow(workflow internal.Workflow, args ...string) error {
	return runMain(func(config internal.Config) (*internal.Context, error) {
		return &internal.Context{
			Config:  internal.Config{},
			Project: internal.Project{Workflow: workflow},
			Tools:   nil,
		}, nil
	}, args...)
}

func runMainWithTools(tools internal.Tools, args ...string) error {
	return runMain(func(config internal.Config) (*internal.Context, error) {
		return &internal.Context{
			Config:  internal.Config{},
			Project: internal.Project{},
			Tools:   tools,
		}, nil
	}, args...)
}
