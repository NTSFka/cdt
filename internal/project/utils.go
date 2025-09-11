package project

import (
	"cdt/internal"
	"context"
	"errors"
	"fmt"
	"strings"
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
	var names []string
	var err error
	var available = 0

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

type configuratorTool interface {
	internal.ProjectConfigurator
	internal.Tool
}

type ConfiguratorFallback []configuratorTool

func (f *ConfiguratorFallback) Configure(ctx context.Context, options internal.ProjectConfiguratorOptions) error {
	return runFirstAvailable(*f, "configurator", func(tool configuratorTool) error {
		return tool.Configure(ctx, options)
	})
}

type builderTool interface {
	internal.ProjectBuilder
	internal.Tool
}

type BuilderFallback []builderTool

func (f *BuilderFallback) BuildAll(ctx context.Context, options internal.ProjectBuilderOptions) error {
	return runFirstAvailable(*f, "builder", func(tool builderTool) error {
		return tool.BuildAll(ctx, options)
	})
}

func (f *BuilderFallback) BuildTargets(ctx context.Context, options internal.ProjectBuilderOptions, targets []string) error {
	return runFirstAvailable(*f, "builder", func(tool builderTool) error {
		return tool.BuildTargets(ctx, options, targets)
	})
}

type testerTool interface {
	internal.ProjectTester
	internal.Tool
}

type TesterFallback []testerTool

func (f *TesterFallback) TestAll(ctx context.Context, options internal.ProjectTesterOptions) error {
	return runFirstAvailable(*f, "tester", func(tool testerTool) error {
		return tool.TestAll(ctx, options)
	})
}

func (f *TesterFallback) TestPattern(ctx context.Context, options internal.ProjectTesterOptions, pattern string) error {
	return runFirstAvailable(*f, "tester", func(tool testerTool) error {
		return tool.TestPattern(ctx, options, pattern)
	})
}

type formatterTool interface {
	internal.ProjectFormatter
	internal.Tool
}

type FormatterFallback []formatterTool

func (f *FormatterFallback) FormatAll(ctx context.Context, options internal.ProjectFormatterOptions) error {
	return runFirstAvailable(*f, "formatter", func(tool formatterTool) error {
		return tool.FormatAll(ctx, options)
	})
}

func (f *FormatterFallback) FormatFiles(ctx context.Context, options internal.ProjectFormatterOptions, filenames []string) error {
	return runFirstAvailable(*f, "formatter", func(tool formatterTool) error {
		return tool.FormatFiles(ctx, options, filenames)
	})
}

func (f *FormatterFallback) FormatCheckAll(ctx context.Context, options internal.ProjectFormatterOptions) error {
	return runFirstAvailable(*f, "formatter", func(tool formatterTool) error {
		return tool.FormatCheckAll(ctx, options)
	})
}

func (f *FormatterFallback) FormatCheckFiles(ctx context.Context, options internal.ProjectFormatterOptions, filenames []string) error {
	return runFirstAvailable(*f, "formatter", func(tool formatterTool) error {
		return tool.FormatCheckFiles(ctx, options, filenames)
	})
}

type linterTool interface {
	internal.ProjectLinter
	internal.Tool
}

type LinterList []linterTool

func (f *LinterList) LintAll(ctx context.Context, options internal.ProjectLinterOptions) error {
	return runAllAvailable(*f, "linter", func(tool linterTool) error {
		return tool.LintAll(ctx, options)
	})
}

func (f *LinterList) LintFiles(ctx context.Context, options internal.ProjectLinterOptions, filenames []string) error {
	return runAllAvailable(*f, "linter", func(tool linterTool) error {
		return tool.LintFiles(ctx, options, filenames)
	})
}

type LinterFallback []linterTool

func (f *LinterFallback) LintAll(ctx context.Context, options internal.ProjectLinterOptions) error {
	return runFirstAvailable(*f, "linter", func(tool linterTool) error {
		return tool.LintAll(ctx, options)
	})
}

func (f *LinterFallback) LintFiles(ctx context.Context, options internal.ProjectLinterOptions, filenames []string) error {
	return runFirstAvailable(*f, "linter", func(tool linterTool) error {
		return tool.LintFiles(ctx, options, filenames)
	})
}

type runnerTool interface {
	internal.ProjectRunner
	internal.Tool
}

type RunnerFallback []runnerTool

func (f *RunnerFallback) RunTarget(ctx context.Context, options internal.ProjectRunnerOptions, target string) error {
	return runFirstAvailable(*f, "runner", func(tool runnerTool) error {
		return tool.RunTarget(ctx, options, target)
	})
}

type dependencyManagerTool interface {
	internal.ProjectDependencyManager
	internal.Tool
}

type DependencyManagerFallback []dependencyManagerTool

func (f *DependencyManagerFallback) AddDependencies(ctx context.Context, options internal.ProjectDependencyManagerOptions, dependencies []string, dev bool) error {
	return runFirstAvailable(*f, "dependency management", func(tool dependencyManagerTool) error {
		return tool.AddDependencies(ctx, options, dependencies, dev)
	})
}

func (f *DependencyManagerFallback) RemoveDependencies(ctx context.Context, options internal.ProjectDependencyManagerOptions, dependencies []string, dev bool) error {
	return runFirstAvailable(*f, "dependency management", func(tool dependencyManagerTool) error {
		return tool.RemoveDependencies(ctx, options, dependencies, dev)
	})
}

func (f *DependencyManagerFallback) UpdateDependencies(ctx context.Context, options internal.ProjectDependencyManagerOptions, dependencies []string) error {
	return runFirstAvailable(*f, "dependency management", func(tool dependencyManagerTool) error {
		return tool.UpdateDependencies(ctx, options, dependencies)
	})
}

func (f *DependencyManagerFallback) FetchDependencies(ctx context.Context, options internal.ProjectDependencyManagerOptions, noDev bool) error {
	return runFirstAvailable(*f, "dependency management", func(tool dependencyManagerTool) error {
		return tool.FetchDependencies(ctx, options, noDev)
	})
}

func (f *DependencyManagerFallback) ListDependencies(ctx context.Context, options internal.ProjectDependencyManagerOptions) error {
	return runFirstAvailable(*f, "dependency management", func(tool dependencyManagerTool) error {
		return tool.ListDependencies(ctx, options)
	})
}

func (f *DependencyManagerFallback) AuditDependencies(ctx context.Context, options internal.ProjectDependencyManagerOptions) error {
	return runFirstAvailable(*f, "dependency management", func(tool dependencyManagerTool) error {
		return tool.AuditDependencies(ctx, options)
	})
}
