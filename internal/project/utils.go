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

func (f *ConfiguratorFallback) Configure(project internal.Project, args []string) error {
	return runFirstAvailable(*f, "configurator", func(tool configuratorTool) error {
		return tool.Configure(project, args)
	})
}

type builderTool interface {
	internal.ProjectBuilder
	internal.Tool
}

type BuilderFallback []builderTool

func (f *BuilderFallback) BuildAll(project internal.Project, args []string) error {
	return runFirstAvailable(*f, "builder", func(tool builderTool) error {
		return tool.BuildAll(project, args)
	})
}

func (f *BuilderFallback) BuildTargets(project internal.Project, targets []string, args []string) error {
	return runFirstAvailable(*f, "builder", func(tool builderTool) error {
		return tool.BuildTargets(project, targets, args)
	})
}

type testerTool interface {
	internal.ProjectTester
	internal.Tool
}

type TesterFallback []testerTool

func (f *TesterFallback) TestAll(project internal.Project, args []string) error {
	return runFirstAvailable(*f, "tester", func(tool testerTool) error {
		return tool.TestAll(project, args)
	})
}

func (f *TesterFallback) Test(project internal.Project, pattern string, args []string) error {
	return runFirstAvailable(*f, "tester", func(tool testerTool) error {
		return tool.Test(project, pattern, args)
	})
}

type formatterTool interface {
	internal.ProjectFormatter
	internal.Tool
}

type FormatterFallback []formatterTool

func (f *FormatterFallback) FormatAll(project internal.Project, args []string) error {
	return runFirstAvailable(*f, "formatter", func(tool formatterTool) error {
		return tool.FormatAll(project, args)
	})
}

func (f *FormatterFallback) FormatFiles(project internal.Project, filenames []string, args []string) error {
	return runFirstAvailable(*f, "formatter", func(tool formatterTool) error {
		return tool.FormatFiles(project, filenames, args)
	})
}

func (f *FormatterFallback) FormatCheckAll(project internal.Project, args []string) error {
	return runFirstAvailable(*f, "formatter", func(tool formatterTool) error {
		return tool.FormatCheckAll(project, args)
	})
}

func (f *FormatterFallback) FormatCheckFiles(project internal.Project, filenames []string, args []string) error {
	return runFirstAvailable(*f, "formatter", func(tool formatterTool) error {
		return tool.FormatCheckFiles(project, filenames, args)
	})
}

type linterTool interface {
	internal.ProjectLinter
	internal.Tool
}

type LinterList []linterTool

func (f *LinterList) LintAll(project internal.Project, args []string) error {
	return runAllAvailable(*f, "linter", func(tool linterTool) error {
		return tool.LintAll(project, args)
	})
}

func (f *LinterList) LintFiles(project internal.Project, filenames []string, args []string) error {
	return runAllAvailable(*f, "linter", func(tool linterTool) error {
		return tool.LintFiles(project, filenames, args)
	})
}

type LinterFallback []linterTool

func (f *LinterFallback) LintAll(project internal.Project, args []string) error {
	return runFirstAvailable(*f, "linter", func(tool linterTool) error {
		return tool.LintAll(project, args)
	})
}

func (f *LinterFallback) LintFiles(project internal.Project, filenames []string, args []string) error {
	return runFirstAvailable(*f, "linter", func(tool linterTool) error {
		return tool.LintFiles(project, filenames, args)
	})
}

type runnerTool interface {
	internal.ProjectRunner
	internal.Tool
}

type RunnerFallback []runnerTool

func (f *RunnerFallback) RunTarget(project internal.Project, target string, args []string) error {
	return runFirstAvailable(*f, "runner", func(tool runnerTool) error {
		return tool.RunTarget(project, target, args)
	})
}

type dependencyManagerTool interface {
	internal.ProjectDependencyManager
	internal.Tool
}

type DependencyManagerFallback []dependencyManagerTool

func (f *DependencyManagerFallback) AddDependencies(project internal.Project, dependencies []string, dev bool) error {
	return runFirstAvailable(*f, "dependency management", func(tool dependencyManagerTool) error {
		return tool.AddDependencies(project, dependencies, dev)
	})
}

func (f *DependencyManagerFallback) RemoveDependencies(project internal.Project, dependencies []string, dev bool) error {
	return runFirstAvailable(*f, "dependency management", func(tool dependencyManagerTool) error {
		return tool.RemoveDependencies(project, dependencies, dev)
	})
}

func (f *DependencyManagerFallback) UpdateDependencies(project internal.Project, dependencies []string) error {
	return runFirstAvailable(*f, "dependency management", func(tool dependencyManagerTool) error {
		return tool.UpdateDependencies(project, dependencies)
	})
}

func (f *DependencyManagerFallback) FetchDependencies(project internal.Project, noDev bool) error {
	return runFirstAvailable(*f, "dependency management", func(tool dependencyManagerTool) error {
		return tool.FetchDependencies(project, noDev)
	})
}

func (f *DependencyManagerFallback) ListDependencies(project internal.Project) error {
	return runFirstAvailable(*f, "dependency management", func(tool dependencyManagerTool) error {
		return tool.ListDependencies(project)
	})
}

func (f *DependencyManagerFallback) AuditDependencies(project internal.Project) error {
	return runFirstAvailable(*f, "dependency management", func(tool dependencyManagerTool) error {
		return tool.AuditDependencies(project)
	})
}
