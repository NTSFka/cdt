package project

import (
	"cdt/internal"
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

func (f *ConfiguratorFallback) Configure(info internal.ProjectInfo, args []string) error {
	return runFirstAvailable(*f, "configurator", func(tool configuratorTool) error {
		return tool.Configure(info, args)
	})
}

type builderTool interface {
	internal.ProjectBuilder
	internal.Tool
}

type BuilderFallback []builderTool

func (f *BuilderFallback) BuildAll(info internal.ProjectInfo, args []string) error {
	return runFirstAvailable(*f, "builder", func(tool builderTool) error {
		return tool.BuildAll(info, args)
	})
}

func (f *BuilderFallback) BuildTargets(info internal.ProjectInfo, targets []string, args []string) error {
	return runFirstAvailable(*f, "builder", func(tool builderTool) error {
		return tool.BuildTargets(info, targets, args)
	})
}

type testerTool interface {
	internal.ProjectTester
	internal.Tool
}

type TesterFallback []testerTool

func (f *TesterFallback) TestAll(info internal.ProjectInfo, args []string) error {
	return runFirstAvailable(*f, "tester", func(tool testerTool) error {
		return tool.TestAll(info, args)
	})
}

func (f *TesterFallback) Test(info internal.ProjectInfo, pattern string, args []string) error {
	return runFirstAvailable(*f, "tester", func(tool testerTool) error {
		return tool.Test(info, pattern, args)
	})
}

type formatterTool interface {
	internal.ProjectFormatter
	internal.Tool
}

type FormatterFallback []formatterTool

func (f *FormatterFallback) FormatAll(info internal.ProjectInfo, args []string) error {
	return runFirstAvailable(*f, "formatter", func(tool formatterTool) error {
		return tool.FormatAll(info, args)
	})
}

func (f *FormatterFallback) FormatFiles(info internal.ProjectInfo, filenames []string, args []string) error {
	return runFirstAvailable(*f, "formatter", func(tool formatterTool) error {
		return tool.FormatFiles(info, filenames, args)
	})
}

func (f *FormatterFallback) FormatCheckAll(info internal.ProjectInfo, args []string) error {
	return runFirstAvailable(*f, "formatter", func(tool formatterTool) error {
		return tool.FormatCheckAll(info, args)
	})
}

func (f *FormatterFallback) FormatCheckFiles(info internal.ProjectInfo, filenames []string, args []string) error {
	return runFirstAvailable(*f, "formatter", func(tool formatterTool) error {
		return tool.FormatCheckFiles(info, filenames, args)
	})
}

type linterTool interface {
	internal.ProjectLinter
	internal.Tool
}

type LinterList []linterTool

func (f *LinterList) LintAll(info internal.ProjectInfo, args []string) error {
	return runAllAvailable(*f, "linter", func(tool linterTool) error {
		return tool.LintAll(info, args)
	})
}

func (f *LinterList) LintFiles(info internal.ProjectInfo, filenames []string, args []string) error {
	return runAllAvailable(*f, "linter", func(tool linterTool) error {
		return tool.LintFiles(info, filenames, args)
	})
}

type LinterFallback []linterTool

func (f *LinterFallback) LintAll(info internal.ProjectInfo, args []string) error {
	return runFirstAvailable(*f, "linter", func(tool linterTool) error {
		return tool.LintAll(info, args)
	})
}

func (f *LinterFallback) LintFiles(info internal.ProjectInfo, filenames []string, args []string) error {
	return runFirstAvailable(*f, "linter", func(tool linterTool) error {
		return tool.LintFiles(info, filenames, args)
	})
}

type runnerTool interface {
	internal.ProjectRunner
	internal.Tool
}

type RunnerFallback []runnerTool

func (f *RunnerFallback) RunTarget(info internal.ProjectInfo, target string, args []string) error {
	return runFirstAvailable(*f, "runner", func(tool runnerTool) error {
		return tool.RunTarget(info, target, args)
	})
}

type dependencyManagerTool interface {
	internal.ProjectDependencyManager
	internal.Tool
}

type DependencyManagerFallback []dependencyManagerTool

func (f *DependencyManagerFallback) AddDependencies(info internal.ProjectInfo, dependencies []string, dev bool) error {
	return runFirstAvailable(*f, "dependency management", func(tool dependencyManagerTool) error {
		return tool.AddDependencies(info, dependencies, dev)
	})
}

func (f *DependencyManagerFallback) RemoveDependencies(info internal.ProjectInfo, dependencies []string, dev bool) error {
	return runFirstAvailable(*f, "dependency management", func(tool dependencyManagerTool) error {
		return tool.RemoveDependencies(info, dependencies, dev)
	})
}

func (f *DependencyManagerFallback) UpdateDependencies(info internal.ProjectInfo, dependencies []string) error {
	return runFirstAvailable(*f, "dependency management", func(tool dependencyManagerTool) error {
		return tool.UpdateDependencies(info, dependencies)
	})
}

func (f *DependencyManagerFallback) FetchDependencies(info internal.ProjectInfo, noDev bool) error {
	return runFirstAvailable(*f, "dependency management", func(tool dependencyManagerTool) error {
		return tool.FetchDependencies(info, noDev)
	})
}

func (f *DependencyManagerFallback) ListDependencies(info internal.ProjectInfo) error {
	return runFirstAvailable(*f, "dependency management", func(tool dependencyManagerTool) error {
		return tool.ListDependencies(info)
	})
}

func (f *DependencyManagerFallback) AuditDependencies(info internal.ProjectInfo) error {
	return runFirstAvailable(*f, "dependency management", func(tool dependencyManagerTool) error {
		return tool.AuditDependencies(info)
	})
}
