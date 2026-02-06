package tool

import (
	"cdt/internal"
	"context"
)

const IdPythonCoverage = "python-coverage"

type PythonCoverage struct {
	internal.ExecutableTool
	detectRunModule func() string
}

// DetectPythonCoverage create a tool for python coverage.
func DetectPythonCoverage(
	ctx context.Context,
	options DetectOptions,
) *PythonCoverage {
	return NewPythonCoverage(func() (*internal.Executable, error) {
		return options.Environment.FindExecutable(
			ctx,
			options.GetToolPath(IdPythonCoverage, "coverage"),
		)
	}, func() string {
		executable, _ := options.Environment.FindExecutable(ctx, "pytest")

		if executable != nil {
			return "pytest"
		}

		return "unittest"
	})
}

// NewPythonCoverage creates a python coverage tool from a custom executable.
func NewPythonCoverage(
	detect internal.ExecutableToolDetectFunc,
	detectRunModule func() string,
) *PythonCoverage {
	if detectRunModule == nil {
		detectRunModule = func() string { return "unittest" }
	}

	return &PythonCoverage{
		ExecutableTool: internal.MakeExecutableTool(
			IdPythonCoverage,
			"Python coverage",
			"Coverage.py is a tool for measuring code coverage of Python programs.",
			internal.Tags{internal.ToolTagPython},
			detect,
		),
		detectRunModule: detectRunModule,
	}
}

// DetectRunModule returns the module to use for test run, e.g., unittest or pytest.
func (p *PythonCoverage) DetectRunModule() string {
	return p.detectRunModule()
}

func (p *PythonCoverage) CollectCoverage(
	ctx context.Context,
	options internal.ProjectCoverageCollectorOptions,
) error {
	args := []string{"run", "-m", p.detectRunModule()}

	if options.Pattern != nil {
		args = append(args, *options.Pattern)
	}

	err := p.RunForProject(
		ctx,
		options.ProjectInfo,
		args,
	)

	if err != nil {
		return err
	}

	return p.RunForProject(
		ctx,
		options.ProjectInfo,
		append([]string{"report"}, options.ExtraArgs...),
	)
}
