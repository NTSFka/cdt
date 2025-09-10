package test

import (
	"cdt/internal"
	"testing"

	"github.com/stretchr/testify/mock"
)

type ProjectConfigurator struct {
	mock.Mock
}

func NewProjectConfigurator(t *testing.T) *ProjectConfigurator {
	configurator := ProjectConfigurator{}
	configurator.Test(t)
	return &configurator
}

func (t *ProjectConfigurator) Configure(info internal.ProjectInfo, args []string) error {
	return t.Called(info, args).Error(0)
}

type ProjectBuilder struct {
	mock.Mock
}

func NewProjectBuilder(t *testing.T) *ProjectBuilder {
	builder := ProjectBuilder{}
	builder.Test(t)
	return &builder
}

func (t *ProjectBuilder) BuildAll(info internal.ProjectInfo, args []string) error {
	return t.Called(info, args).Error(0)
}

func (t *ProjectBuilder) BuildTargets(info internal.ProjectInfo, targets []string, args []string) error {
	return t.Called(info, targets, args).Error(0)
}

type ProjectFormatter struct {
	mock.Mock
}

func NewProjectFormatter(t *testing.T) *ProjectFormatter {
	formatter := ProjectFormatter{}
	formatter.Test(t)
	return &formatter
}

func (t *ProjectFormatter) FormatAll(info internal.ProjectInfo, args []string) error {
	return t.Called(info, args).Error(0)
}

func (t *ProjectFormatter) FormatFiles(info internal.ProjectInfo, filenames []string, args []string) error {
	return t.Called(info, filenames, args).Error(0)
}

func (t *ProjectFormatter) FormatCheckAll(info internal.ProjectInfo, args []string) error {
	return t.Called(info, args).Error(0)
}

func (t *ProjectFormatter) FormatCheckFiles(info internal.ProjectInfo, filenames []string, args []string) error {
	return t.Called(info, filenames, args).Error(0)
}

type ProjectLinter struct {
	mock.Mock
}

func NewProjectLinter(t *testing.T) *ProjectLinter {
	linter := ProjectLinter{}
	linter.Test(t)
	return &linter
}

func (t *ProjectLinter) LintAll(info internal.ProjectInfo, args []string) error {
	return t.Called(info, args).Error(0)
}

func (t *ProjectLinter) LintFiles(info internal.ProjectInfo, filenames []string, args []string) error {
	return t.Called(info, filenames, args).Error(0)
}

type ProjectTester struct {
	mock.Mock
}

func NewProjectTester(t *testing.T) *ProjectTester {
	tester := ProjectTester{}
	tester.Mock.Test(t)
	return &tester
}

func (t *ProjectTester) TestAll(info internal.ProjectInfo, args []string) error {
	return t.Called(info, args).Error(0)
}

func (t *ProjectTester) Test(info internal.ProjectInfo, pattern string, args []string) error {
	return t.Called(info, pattern, args).Error(0)
}

type ProjectRunner struct {
	mock.Mock
}

func NewProjectRunner(t *testing.T) *ProjectRunner {
	runner := ProjectRunner{}
	runner.Test(t)
	return &runner
}

func (t *ProjectRunner) RunTarget(info internal.ProjectInfo, target string, args []string) error {
	return t.Called(info, target, args).Error(0)
}

type StructureProvider struct {
	mock.Mock
}

func NewStructureProvider(t *testing.T) *StructureProvider {
	provider := StructureProvider{}
	provider.Test(t)
	return &provider
}

func (t *StructureProvider) Structure(info internal.ProjectInfo) (*internal.ProjectStructure, error) {
	called := t.Called(info)
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

func (d *DependencyManager) AddDependencies(info internal.ProjectInfo, dependencies []string, dev bool) error {
	return d.Called(info, dependencies, dev).Error(0)
}

func (d *DependencyManager) RemoveDependencies(info internal.ProjectInfo, dependencies []string, dev bool) error {
	return d.Called(info, dependencies, dev).Error(0)
}

func (d *DependencyManager) UpdateDependencies(info internal.ProjectInfo, dependencies []string) error {
	return d.Called(info, dependencies).Error(0)
}

func (d *DependencyManager) FetchDependencies(info internal.ProjectInfo, noDev bool) error {
	return d.Called(info, noDev).Error(0)
}

func (d *DependencyManager) ListDependencies(info internal.ProjectInfo) error {
	return d.Called(info).Error(0)
}

func (d *DependencyManager) AuditDependencies(info internal.ProjectInfo) error {
	return d.Called(info).Error(0)
}
