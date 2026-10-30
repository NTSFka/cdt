package tool

import (
	"cdt/internal"
	"cdt/internal/utils/cobertura"
	"cdt/internal/utils/junit"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/ctrf-io/go-ctrf-json-reporter/ctrf"
)

const IdPHPUnit = "phpunit"

type PHPUnit struct {
	internal.ExecutableTool
}

// DetectPHPUnit create a tool for phpunit.
func DetectPHPUnit(
	ctx context.Context,
	options internal.DetectOptions,
) *PHPUnit {
	if path, ok := options.ToolsPaths[IdPHPUnit]; ok {
		return NewPHPUnit(func() (*internal.Executable, error) {
			return options.Environment.FindExecutable(ctx, path)
		})
	}

	return NewPHPUnit(internal.DetectExecutableChain(
		[]string{filepath.Join(options.ProjectDirectory, "vendor/bin/phpunit"), "phpunit"},
		func(name string) (*internal.Executable, error) {
			return options.Environment.FindExecutable(ctx, name)
		},
	))
}

// NewPHPUnit creates a phpunit tool from a custom executable.
func NewPHPUnit(detect internal.ExecutableToolDetectFunc) *PHPUnit {
	return &PHPUnit{
		ExecutableTool: internal.MakeExecutableTool(
			IdPHPUnit,
			"PHPUnit",
			"PHPUnit is a programmer-oriented testing framework for PHP.",
			internal.Tags{internal.ToolTagPhp, internal.ToolTagTest},
			detect,
		),
	}
}

func (p *PHPUnit) RunTests(ctx context.Context, options internal.ProjectTesterOptions) error {
	args := options.ExtraArgs

	if options.Pattern != nil {
		args = append(args, *options.Pattern)
	}

	if options.Output.Reporter != nil {
		return p.runTestsReporter(ctx, options, options.Output.Reporter, args)
	}

	return p.RunForProjectWithWriter(ctx, options.ProjectInfo, options.Output.Output, args)
}

func (p *PHPUnit) CollectCoverage(
	ctx context.Context,
	options internal.ProjectCoverageCollectorOptions,
) error {
	args := options.ExtraArgs

	if options.Output.Reporter != nil {
		return p.collectCoverageReporter(ctx, options, options.Output.Reporter, args)
	}

	args = append(args, "--coverage-text")

	if options.Pattern != nil {
		args = append(args, *options.Pattern)
	}

	return p.RunForProjectWithEnv(
		ctx,
		options.ProjectInfo,
		[]string{"XDEBUG_MODE=coverage"},
		args,
	)
}

func (p *PHPUnit) runTestsReporter(
	ctx context.Context,
	options internal.ProjectTesterOptions,
	reporter internal.TestsReporter,
	args []string,
) error {
	file, err := os.CreateTemp("", "phpunit-report.*.xml")
	if err != nil {
		return fmt.Errorf("failed to create temporary file for report: %w", err)
	}
	defer func() { _ = os.Remove(file.Name()) }()
	defer func() { _ = file.Close() }()

	args = append(args, "--log-junit", file.Name())

	if options.Pattern != nil {
		args = append(args, *options.Pattern)
	}

	if err = p.RunForProject(ctx, options.ProjectInfo, args); err != nil {
		return fmt.Errorf("failed to run PHPUnit: %w", err)
	}

	return p.reportJUnitOutput(file, reporter)
}

func (p *PHPUnit) reportJUnitOutput(
	reader io.Reader,
	testsReporter internal.TestsReporter,
) error {
	// FIXME: recursive test suites
	var junitReport junit.Testsuites
	if err := xml.NewDecoder(reader).Decode(&junitReport); err != nil {
		return err
	}

	var testResults []*ctrf.TestResult

	for _, suite := range junitReport.Suites {
		for _, testCase := range suite.Testcases {
			var (
				status  ctrf.TestStatus
				message string
			)

			switch {
			case testCase.Failure != nil:
				status = ctrf.TestFailed
				message = testCase.Failure.Message
			case testCase.Skipped != nil:
				status = ctrf.TestSkipped
				message = testCase.Skipped.Message
			default:
				status = ctrf.TestPassed
			}

			testResults = append(testResults, &ctrf.TestResult{
				Name:    testCase.Name,
				Status:  status,
				Message: message,
			})
		}
	}

	report := ctrf.Report{
		Results: &ctrf.Results{
			Tests: testResults,
		},
	}

	return testsReporter.AddCtrf(&report)
}

func (p *PHPUnit) collectCoverageReporter(
	ctx context.Context,
	options internal.ProjectCoverageCollectorOptions,
	reporter internal.CoverageReporter,
	args []string,
) error {
	file, err := os.CreateTemp("", "phpunit-coverage.*.xml")
	if err != nil {
		return fmt.Errorf("failed to create temporary file for report: %w", err)
	}
	defer func() { _ = os.Remove(file.Name()) }()
	defer func() { _ = file.Close() }()

	args = append(args, "--coverage-cobertura", file.Name())

	if options.Pattern != nil {
		args = append(args, *options.Pattern)
	}

	err = p.RunForProjectWithEnv(
		ctx,
		options.ProjectInfo,
		[]string{"XDEBUG_MODE=coverage"},
		args,
	)
	if err != nil {
		return fmt.Errorf("failed to run PHPUnit with coverage: %w", err)
	}

	return p.reportCoberturaOutput(file, reporter)
}

func (p *PHPUnit) reportCoberturaOutput(
	reader io.Reader,
	coverageReporter internal.CoverageReporter,
) error {
	decoder := xml.NewDecoder(reader)

	var report cobertura.Coverage
	if err := decoder.Decode(&report); err != nil {
		return fmt.Errorf("failed to decode Cobertura report: %w", err)
	}

	return coverageReporter.AddCobertura(&report)
}
