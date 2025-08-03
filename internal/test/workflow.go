package test

import (
	"cdt/internal"
	"github.com/stretchr/testify/mock"
	"testing"
)

type ProjectConfigurator struct {
	mock.Mock
}

func NewProjectConfigurator(t *testing.T) *ProjectConfigurator {
	configurator := ProjectConfigurator{}
	configurator.Test(t)
	return &configurator
}

func (t *ProjectConfigurator) Configure(project internal.Project, args []string) error {
	return t.Called(project, args).Error(0)
}

type ProjectBuilder struct {
	mock.Mock
}

func NewProjectBuilder(t *testing.T) *ProjectBuilder {
	builder := ProjectBuilder{}
	builder.Test(t)
	return &builder
}

func (t *ProjectBuilder) BuildAll(project internal.Project, args []string) error {
	return t.Called(project, args).Error(0)
}

func (t *ProjectBuilder) BuildTargets(project internal.Project, targets []string, args []string) error {
	return t.Called(project, targets, args).Error(0)
}

type ProjectFormatter struct {
	mock.Mock
}

func NewProjectFormatter(t *testing.T) *ProjectFormatter {
	formatter := ProjectFormatter{}
	formatter.Test(t)
	return &formatter
}

func (t *ProjectFormatter) FormatAll(project internal.Project, args []string) error {
	return t.Called(project, args).Error(0)
}

func (t *ProjectFormatter) FormatFiles(project internal.Project, filenames []string, args []string) error {
	return t.Called(project, filenames, args).Error(0)
}

func (t *ProjectFormatter) FormatCheckAll(project internal.Project, args []string) error {
	return t.Called(project, args).Error(0)
}

func (t *ProjectFormatter) FormatCheckFiles(project internal.Project, filenames []string, args []string) error {
	return t.Called(project, filenames, args).Error(0)
}

type ProjectLinter struct {
	mock.Mock
}

func NewProjectLinter(t *testing.T) *ProjectLinter {
	linter := ProjectLinter{}
	linter.Test(t)
	return &linter
}

func (t *ProjectLinter) LintAll(project internal.Project, args []string) error {
	return t.Called(project, args).Error(0)
}

func (t *ProjectLinter) LintFiles(project internal.Project, filenames []string, args []string) error {
	return t.Called(project, filenames, args).Error(0)
}

type ProjectTester struct {
	mock.Mock
}

func NewProjectTester(t *testing.T) *ProjectTester {
	tester := ProjectTester{}
	tester.Mock.Test(t)
	return &tester
}

func (t *ProjectTester) TestAll(project internal.Project, args []string) error {
	return t.Called(project, args).Error(0)
}

func (t *ProjectTester) Test(project internal.Project, pattern string, args []string) error {
	return t.Called(project, pattern, args).Error(0)
}

type ProjectRunner struct {
	mock.Mock
}

func NewProjectRunner(t *testing.T) *ProjectRunner {
	runner := ProjectRunner{}
	runner.Test(t)
	return &runner
}

func (t *ProjectRunner) RunTarget(project internal.Project, target string, args []string) error {
	return t.Called(project, target, args).Error(0)
}

type StructureProvider struct {
	mock.Mock
}

func NewStructureProvider(t *testing.T) *StructureProvider {
	provider := StructureProvider{}
	provider.Test(t)
	return &provider
}

func (t *StructureProvider) Structure(project internal.Project) (*internal.ProjectStructure, error) {
	called := t.Called(project)
	return called.Get(0).(*internal.ProjectStructure), called.Error(1)
}

type DependencyManager struct {
	mock.Mock
}

func NewDependencyManager(t *testing.T) *DependencyManager {
	manager := DependencyManager{}
	manager.Test(t)
	return &manager
}

func (d *DependencyManager) AddDependencies(project internal.Project, dependencies []string, dev bool) error {
	return d.Called(project, dependencies, dev).Error(0)
}

func (d *DependencyManager) RemoveDependencies(project internal.Project, dependencies []string, dev bool) error {
	return d.Called(project, dependencies, dev).Error(0)
}

func (d *DependencyManager) UpdateDependencies(project internal.Project, dependencies []string) error {
	return d.Called(project, dependencies).Error(0)
}

func (d *DependencyManager) FetchDependencies(project internal.Project, noDev bool) error {
	return d.Called(project, noDev).Error(0)
}

func (d *DependencyManager) ListDependencies(project internal.Project) error {
	return d.Called(project).Error(0)
}

func (d *DependencyManager) AuditDependencies(project internal.Project) error {
	return d.Called(project).Error(0)
}
