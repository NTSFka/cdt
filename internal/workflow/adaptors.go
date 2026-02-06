package workflow

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"cdt/internal"
)

func runFirstAvailable[T internal.Tool](tools []T, name string, run func(T) error) error {
	var names []string

	for _, tool := range tools {
		if tool.IsAvailable() {
			return run(tool)
		}

		names = append(names, tool.Id())
	}

	if len(names) == 0 {
		names = append(names, "none")
	}

	return fmt.Errorf("no %v tool available: %v", name, strings.Join(names, ", "))
}

func runAllAvailable[T internal.Tool](tools []T, name string, run func(T) error) error {
	var (
		names     []string
		err       error
		available = 0
	)

	for _, tool := range tools {
		if tool.IsAvailable() {
			available++
			err = errors.Join(err, run(tool))
		}

		names = append(names, tool.Id())
	}

	if len(names) == 0 {
		names = append(names, "none")
	}

	if available == 0 {
		return fmt.Errorf("no %v tool available: %v", name, strings.Join(names, ", "))
	}

	return err
}

// Adaptor returns information about workflow adaptor.
type Adaptor interface {
	Details() string
}

func adaptorToolIds[T internal.Tool](name string, tools []T) string {
	ids := make([]string, 0, len(tools))

	for _, tool := range tools {
		ids = append(ids, tool.Id())
	}

	return fmt.Sprintf("%v (%v)", name, strings.Join(ids, ", "))
}

type configuratorTool interface {
	internal.ProjectConfigurator
	internal.Tool
}

// ConfiguratorFallback run the first available configurator.
type ConfiguratorFallback []configuratorTool

func (f *ConfiguratorFallback) Details() string {
	return adaptorToolIds("fallback", *f)
}

func (f *ConfiguratorFallback) Configure(
	ctx context.Context,
	options internal.ProjectConfiguratorOptions,
) error {
	return runFirstAvailable(*f, "configurator", func(tool configuratorTool) error {
		return tool.Configure(ctx, options)
	})
}

type builderTool interface {
	internal.ProjectBuilder
	internal.Tool
}

// BuilderFallback run the first available builder.
type BuilderFallback []builderTool

func (f *BuilderFallback) Details() string {
	return adaptorToolIds("fallback", *f)
}

func (f *BuilderFallback) BuildAll(
	ctx context.Context,
	options internal.ProjectBuilderOptions,
) error {
	return runFirstAvailable(*f, "builder", func(tool builderTool) error {
		return tool.BuildAll(ctx, options)
	})
}

func (f *BuilderFallback) BuildTargets(
	ctx context.Context,
	options internal.ProjectBuilderOptions,
	targets []string,
) error {
	return runFirstAvailable(*f, "builder", func(tool builderTool) error {
		return tool.BuildTargets(ctx, options, targets)
	})
}

type testerTool interface {
	internal.ProjectTester
	internal.Tool
}

// TesterFallback run the first available tester.
type TesterFallback []testerTool

func (f *TesterFallback) Details() string {
	return adaptorToolIds("fallback", *f)
}

func (f *TesterFallback) RunTests(
	ctx context.Context,
	options internal.ProjectTesterOptions,
) error {
	return runFirstAvailable(*f, "tester", func(tool testerTool) error {
		return tool.RunTests(ctx, options)
	})
}

type coverageCollectorTool interface {
	internal.ProjectCoverageCollector
	internal.Tool
}

// CoverageCollectorFallback run the first available coverage collector.
type CoverageCollectorFallback []coverageCollectorTool

func (f *CoverageCollectorFallback) Details() string {
	return adaptorToolIds("fallback", *f)
}

func (f *CoverageCollectorFallback) CollectCoverage(
	ctx context.Context,
	options internal.ProjectCoverageCollectorOptions,
) error {
	return runFirstAvailable(*f, "coverage", func(tool coverageCollectorTool) error {
		return tool.CollectCoverage(ctx, options)
	})
}

type formatterTool interface {
	internal.ProjectFormatter
	internal.Tool
}

// FormatterFallback run the first available formatter.
type FormatterFallback []formatterTool

func (f *FormatterFallback) Details() string {
	return adaptorToolIds("fallback", *f)
}

func (f *FormatterFallback) FormatFiles(
	ctx context.Context,
	options internal.ProjectFormatterOptions,
) error {
	return runFirstAvailable(*f, "formatter", func(tool formatterTool) error {
		return tool.FormatFiles(ctx, options)
	})
}

type linterTool interface {
	internal.ProjectLinter
	internal.Tool
}

// LinterList run all available linters.
type LinterList []linterTool

func (f *LinterList) Details() string {
	return adaptorToolIds("list", *f)
}

func (f *LinterList) LintFiles(
	ctx context.Context,
	options internal.ProjectLinterOptions,
) error {
	return runAllAvailable(*f, "linter", func(tool linterTool) error {
		return tool.LintFiles(ctx, options)
	})
}

// LinterFallback run the first available linter.
type LinterFallback []linterTool

func (f *LinterFallback) Details() string {
	return adaptorToolIds("fallback", *f)
}

func (f *LinterFallback) LintFiles(
	ctx context.Context,
	options internal.ProjectLinterOptions,
) error {
	return runFirstAvailable(*f, "linter", func(tool linterTool) error {
		return tool.LintFiles(ctx, options)
	})
}

type runnerTool interface {
	internal.ProjectRunner
	internal.Tool
}

// RunnerFallback run the first available runner.
type RunnerFallback []runnerTool

func (f *RunnerFallback) Details() string {
	return adaptorToolIds("fallback", *f)
}

func (f *RunnerFallback) RunTarget(
	ctx context.Context,
	options internal.ProjectRunnerOptions,
	target string,
) error {
	return runFirstAvailable(*f, "runner", func(tool runnerTool) error {
		return tool.RunTarget(ctx, options, target)
	})
}

type dependencyManagerTool interface {
	internal.ProjectDependencyManager
	internal.Tool
}

// DependencyManagerFallback run the first available dependency manager.
type DependencyManagerFallback []dependencyManagerTool

func (f *DependencyManagerFallback) Details() string {
	return adaptorToolIds("fallback", *f)
}

func (f *DependencyManagerFallback) AddDependencies(
	ctx context.Context,
	options internal.ProjectDependencyManagerOptions,
	dependencies []string,
	dev bool,
) error {
	return runFirstAvailable(*f, "dependency management", func(tool dependencyManagerTool) error {
		return tool.AddDependencies(ctx, options, dependencies, dev)
	})
}

func (f *DependencyManagerFallback) RemoveDependencies(
	ctx context.Context,
	options internal.ProjectDependencyManagerOptions,
	dependencies []string,
	dev bool,
) error {
	return runFirstAvailable(*f, "dependency management", func(tool dependencyManagerTool) error {
		return tool.RemoveDependencies(ctx, options, dependencies, dev)
	})
}

func (f *DependencyManagerFallback) UpdateDependencies(
	ctx context.Context,
	options internal.ProjectDependencyManagerOptions,
	dependencies []string,
) error {
	return runFirstAvailable(*f, "dependency management", func(tool dependencyManagerTool) error {
		return tool.UpdateDependencies(ctx, options, dependencies)
	})
}

func (f *DependencyManagerFallback) FetchDependencies(
	ctx context.Context,
	options internal.ProjectDependencyManagerOptions,
	noDev bool,
) error {
	return runFirstAvailable(*f, "dependency management", func(tool dependencyManagerTool) error {
		return tool.FetchDependencies(ctx, options, noDev)
	})
}

func (f *DependencyManagerFallback) ListDependencies(
	ctx context.Context,
	options internal.ProjectDependencyManagerOptions,
) error {
	return runFirstAvailable(*f, "dependency management", func(tool dependencyManagerTool) error {
		return tool.ListDependencies(ctx, options)
	})
}

func (f *DependencyManagerFallback) AuditDependencies(
	ctx context.Context,
	options internal.ProjectDependencyManagerOptions,
) error {
	return runFirstAvailable(*f, "dependency management", func(tool dependencyManagerTool) error {
		return tool.AuditDependencies(ctx, options)
	})
}
