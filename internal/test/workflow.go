package test

import (
	"cdt/internal"
	"github.com/stretchr/testify/mock"
)

type ProjectConfigurator struct {
	mock.Mock
}

func (t *ProjectConfigurator) Configure(project internal.Project, args []string) error {
	return t.Called(project, args).Error(0)
}

type ProjectBuilder struct {
	mock.Mock
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

func (t *ProjectLinter) LintAll(project internal.Project, args []string) error {
	return t.Called(project, args).Error(0)
}

func (t *ProjectLinter) LintFiles(project internal.Project, filenames []string, args []string) error {
	return t.Called(project, filenames, args).Error(0)
}

type ProjectTester struct {
	mock.Mock
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

func (t *ProjectRunner) RunTarget(project internal.Project, target string, args []string) error {
	return t.Called(project, target, args).Error(0)
}

type StructureProvider struct {
	mock.Mock
}

func (t *StructureProvider) Structure(project internal.Project) (*internal.ProjectStructure, error) {
	called := t.Called(project)
	return called.Get(0).(*internal.ProjectStructure), called.Error(1)
}

type DependencyManager struct {
	mock.Mock
}

func (d *DependencyManager) AddDependencies(project internal.Project, dependencies []string) error {
	return d.Called(project, dependencies).Error(0)
}

func (d *DependencyManager) RemoveDependencies(project internal.Project, dependencies []string) error {
	return d.Called(project, dependencies).Error(0)
}

func (d *DependencyManager) UpdateDependencies(project internal.Project, dependencies []string) error {
	return d.Called(project, dependencies).Error(0)
}

func (d *DependencyManager) FetchDependencies(project internal.Project) error {
	return d.Called(project).Error(0)
}

func (d *DependencyManager) ListDependencies(project internal.Project) error {
	return d.Called(project).Error(0)
}

func (d *DependencyManager) AuditDependencies(project internal.Project) error {
	return d.Called(project).Error(0)
}
