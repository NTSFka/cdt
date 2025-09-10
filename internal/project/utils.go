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

func (f *ConfiguratorFallback) Configure(ctx context.Context, info internal.ProjectInfo, args []string) error {
	return runFirstAvailable(*f, "configurator", func(tool configuratorTool) error {
		return tool.Configure(ctx, info, args)
	})
}

type builderTool interface {
	internal.ProjectBuilder
	internal.Tool
}

type BuilderFallback []builderTool

func (f *BuilderFallback) BuildAll(ctx context.Context, info internal.ProjectInfo, args []string) error {
	return runFirstAvailable(*f, "builder", func(tool builderTool) error {
		return tool.BuildAll(ctx, info, args)
	})
}

func (f *BuilderFallback) BuildTargets(ctx context.Context, info internal.ProjectInfo, targets []string, args []string) error {
	return runFirstAvailable(*f, "builder", func(tool builderTool) error {
		return tool.BuildTargets(ctx, info, targets, args)
	})
}

type testerTool interface {
	internal.ProjectTester
	internal.Tool
}

type TesterFallback []testerTool

func (f *TesterFallback) TestAll(ctx context.Context, info internal.ProjectInfo, args []string) error {
	return runFirstAvailable(*f, "tester", func(tool testerTool) error {
		return tool.TestAll(ctx, info, args)
	})
}

func (f *TesterFallback) Test(ctx context.Context, info internal.ProjectInfo, pattern string, args []string) error {
	return runFirstAvailable(*f, "tester", func(tool testerTool) error {
		return tool.Test(ctx, info, pattern, args)
	})
}

type formatterTool interface {
	internal.ProjectFormatter
	internal.Tool
}

type FormatterFallback []formatterTool

func (f *FormatterFallback) FormatAll(ctx context.Context, info internal.ProjectInfo, args []string) error {
	return runFirstAvailable(*f, "formatter", func(tool formatterTool) error {
		return tool.FormatAll(ctx, info, args)
	})
}

func (f *FormatterFallback) FormatFiles(ctx context.Context, info internal.ProjectInfo, filenames []string, args []string) error {
	return runFirstAvailable(*f, "formatter", func(tool formatterTool) error {
		return tool.FormatFiles(ctx, info, filenames, args)
	})
}

func (f *FormatterFallback) FormatCheckAll(ctx context.Context, info internal.ProjectInfo, args []string) error {
	return runFirstAvailable(*f, "formatter", func(tool formatterTool) error {
		return tool.FormatCheckAll(ctx, info, args)
	})
}

func (f *FormatterFallback) FormatCheckFiles(ctx context.Context, info internal.ProjectInfo, filenames []string, args []string) error {
	return runFirstAvailable(*f, "formatter", func(tool formatterTool) error {
		return tool.FormatCheckFiles(ctx, info, filenames, args)
	})
}

type linterTool interface {
	internal.ProjectLinter
	internal.Tool
}

type LinterList []linterTool

func (f *LinterList) LintAll(ctx context.Context, info internal.ProjectInfo, args []string) error {
	return runAllAvailable(*f, "linter", func(tool linterTool) error {
		return tool.LintAll(ctx, info, args)
	})
}

func (f *LinterList) LintFiles(ctx context.Context, info internal.ProjectInfo, filenames []string, args []string) error {
	return runAllAvailable(*f, "linter", func(tool linterTool) error {
		return tool.LintFiles(ctx, info, filenames, args)
	})
}

type LinterFallback []linterTool

func (f *LinterFallback) LintAll(ctx context.Context, info internal.ProjectInfo, args []string) error {
	return runFirstAvailable(*f, "linter", func(tool linterTool) error {
		return tool.LintAll(ctx, info, args)
	})
}

func (f *LinterFallback) LintFiles(ctx context.Context, info internal.ProjectInfo, filenames []string, args []string) error {
	return runFirstAvailable(*f, "linter", func(tool linterTool) error {
		return tool.LintFiles(ctx, info, filenames, args)
	})
}

type runnerTool interface {
	internal.ProjectRunner
	internal.Tool
}

type RunnerFallback []runnerTool

func (f *RunnerFallback) RunTarget(ctx context.Context, info internal.ProjectInfo, target string, args []string) error {
	return runFirstAvailable(*f, "runner", func(tool runnerTool) error {
		return tool.RunTarget(ctx, info, target, args)
	})
}

type dependencyManagerTool interface {
	internal.ProjectDependencyManager
	internal.Tool
}

type DependencyManagerFallback []dependencyManagerTool

func (f *DependencyManagerFallback) AddDependencies(ctx context.Context, info internal.ProjectInfo, dependencies []string, dev bool) error {
	return runFirstAvailable(*f, "dependency management", func(tool dependencyManagerTool) error {
		return tool.AddDependencies(ctx, info, dependencies, dev)
	})
}

func (f *DependencyManagerFallback) RemoveDependencies(ctx context.Context, info internal.ProjectInfo, dependencies []string, dev bool) error {
	return runFirstAvailable(*f, "dependency management", func(tool dependencyManagerTool) error {
		return tool.RemoveDependencies(ctx, info, dependencies, dev)
	})
}

func (f *DependencyManagerFallback) UpdateDependencies(ctx context.Context, info internal.ProjectInfo, dependencies []string) error {
	return runFirstAvailable(*f, "dependency management", func(tool dependencyManagerTool) error {
		return tool.UpdateDependencies(ctx, info, dependencies)
	})
}

func (f *DependencyManagerFallback) FetchDependencies(ctx context.Context, info internal.ProjectInfo, noDev bool) error {
	return runFirstAvailable(*f, "dependency management", func(tool dependencyManagerTool) error {
		return tool.FetchDependencies(ctx, info, noDev)
	})
}

func (f *DependencyManagerFallback) ListDependencies(ctx context.Context, info internal.ProjectInfo) error {
	return runFirstAvailable(*f, "dependency management", func(tool dependencyManagerTool) error {
		return tool.ListDependencies(ctx, info)
	})
}

func (f *DependencyManagerFallback) AuditDependencies(ctx context.Context, info internal.ProjectInfo) error {
	return runFirstAvailable(*f, "dependency management", func(tool dependencyManagerTool) error {
		return tool.AuditDependencies(ctx, info)
	})
}
