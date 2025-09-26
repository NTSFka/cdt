package test

import (
	"cdt/internal"
	"context"
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

func (t *ProjectConfigurator) Configure(ctx context.Context, options internal.ProjectConfiguratorOptions) error {
	return t.Called(ctx, options).Error(0)
}

type ProjectBuilder struct {
	mock.Mock
}

func NewProjectBuilder(t *testing.T) *ProjectBuilder {
	builder := ProjectBuilder{}
	builder.Test(t)

	return &builder
}

func (t *ProjectBuilder) BuildAll(ctx context.Context, options internal.ProjectBuilderOptions) error {
	return t.Called(ctx, options).Error(0)
}

func (t *ProjectBuilder) BuildTargets(ctx context.Context, options internal.ProjectBuilderOptions, targets []string) error {
	return t.Called(ctx, options, targets).Error(0)
}

type ProjectFormatter struct {
	mock.Mock
}

func NewProjectFormatter(t *testing.T) *ProjectFormatter {
	formatter := ProjectFormatter{}
	formatter.Test(t)

	return &formatter
}

func (t *ProjectFormatter) FormatAll(ctx context.Context, options internal.ProjectFormatterOptions) error {
	return t.Called(ctx, options).Error(0)
}

func (t *ProjectFormatter) FormatFiles(ctx context.Context, options internal.ProjectFormatterOptions, filenames []string) error {
	return t.Called(ctx, options, filenames).Error(0)
}

func (t *ProjectFormatter) FormatCheckAll(ctx context.Context, options internal.ProjectFormatterOptions) error {
	return t.Called(ctx, options).Error(0)
}

func (t *ProjectFormatter) FormatCheckFiles(ctx context.Context, options internal.ProjectFormatterOptions, filenames []string) error {
	return t.Called(ctx, options, filenames).Error(0)
}

type ProjectLinter struct {
	mock.Mock
}

func NewProjectLinter(t *testing.T) *ProjectLinter {
	linter := ProjectLinter{}
	linter.Test(t)

	return &linter
}

func (t *ProjectLinter) LintAll(ctx context.Context, options internal.ProjectLinterOptions) error {
	return t.Called(ctx, options).Error(0)
}

func (t *ProjectLinter) LintFiles(ctx context.Context, options internal.ProjectLinterOptions, filenames []string) error {
	return t.Called(ctx, options, filenames).Error(0)
}

type ProjectTester struct {
	mock.Mock
}

func NewProjectTester(t *testing.T) *ProjectTester {
	tester := ProjectTester{}
	tester.Test(t)

	return &tester
}

func (t *ProjectTester) TestAll(ctx context.Context, options internal.ProjectTesterOptions) error {
	return t.Called(ctx, options).Error(0)
}

func (t *ProjectTester) TestPattern(ctx context.Context, options internal.ProjectTesterOptions, pattern string) error {
	return t.Called(ctx, options, pattern).Error(0)
}

type ProjectRunner struct {
	mock.Mock
}

func NewProjectRunner(t *testing.T) *ProjectRunner {
	runner := ProjectRunner{}
	runner.Test(t)

	return &runner
}

func (t *ProjectRunner) RunTarget(ctx context.Context, options internal.ProjectRunnerOptions, target string) error {
	return t.Called(ctx, options, target).Error(0)
}

type StructureProvider struct {
	mock.Mock
}

func NewStructureProvider(t *testing.T) *StructureProvider {
	provider := StructureProvider{}
	provider.Test(t)

	return &provider
}

func (t *StructureProvider) Structure(ctx context.Context, info internal.ProjectInfo) (*internal.ProjectStructure, error) {
	called := t.Called(ctx, info)

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

func (d *DependencyManager) AddDependencies(ctx context.Context, options internal.ProjectDependencyManagerOptions, dependencies []string, dev bool) error {
	return d.Called(ctx, options, dependencies, dev).Error(0)
}

func (d *DependencyManager) RemoveDependencies(ctx context.Context, options internal.ProjectDependencyManagerOptions, dependencies []string, dev bool) error {
	return d.Called(ctx, options, dependencies, dev).Error(0)
}

func (d *DependencyManager) UpdateDependencies(ctx context.Context, options internal.ProjectDependencyManagerOptions, dependencies []string) error {
	return d.Called(ctx, options, dependencies).Error(0)
}

func (d *DependencyManager) FetchDependencies(ctx context.Context, options internal.ProjectDependencyManagerOptions, noDev bool) error {
	return d.Called(ctx, options, noDev).Error(0)
}

func (d *DependencyManager) ListDependencies(ctx context.Context, options internal.ProjectDependencyManagerOptions) error {
	return d.Called(ctx, options).Error(0)
}

func (d *DependencyManager) AuditDependencies(ctx context.Context, options internal.ProjectDependencyManagerOptions) error {
	return d.Called(ctx, options).Error(0)
}
