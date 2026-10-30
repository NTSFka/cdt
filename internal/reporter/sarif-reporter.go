package reporter

import (
	"cdt/internal"
	"errors"
	"io"

	"github.com/owenrumney/go-sarif/v3/pkg/report/merge"
	"github.com/owenrumney/go-sarif/v3/pkg/report/v210/sarif"
)

// SarifReporter reports results in a SARIF format.
type SarifReporter struct {
	Report *sarif.Report
}

func (r *SarifReporter) AddSarif(report *sarif.Report) error {
	if report == nil {
		return errors.New("report cannot be nil")
	}

	if r.Report == nil {
		r.Report = report

		return nil
	}

	// Merge failed only when any of the reports are nil, which is not true at this point
	r.Report, _ = merge.V210(r.Report, report)

	return nil
}

func (r *SarifReporter) Write(output io.Writer) error {
	if r.Report == nil {
		return nil
	}

	return r.Report.Write(output)
}

// SimpleSarifReporter reports SARIF results in a simple format.
type SimpleSarifReporter struct {
	LineReporter
	SarifReporter
}

// NewSimpleSarifReporter creates a new SimpleSarifReporter.
func NewSimpleSarifReporter(output io.Writer) *SimpleSarifReporter {
	return &SimpleSarifReporter{
		LineReporter:  LineReporter{Output: output},
		SarifReporter: SarifReporter{},
	}
}

func (r *SimpleSarifReporter) Finish() error {
	count, err := r.printReport(r.Report)
	if err != nil {
		return err
	}

	return r.WriteLine("\nTotal issues: %d\n", count)
}

func (r *SimpleSarifReporter) printReport(report *sarif.Report) (int, error) {
	if report == nil {
		return 0, nil
	}

	count := 0

	for _, run := range report.Runs {
		runCount, err := r.printRun(run)
		if err != nil {
			return 0, err
		}

		count += runCount
	}

	return count, nil
}

func (r *SimpleSarifReporter) printRun(run *sarif.Run) (int, error) {
	if len(run.Results) == 0 {
		return 0, nil
	}

	var groupName *string = nil
	if run.Tool != nil && run.Tool.Driver != nil && run.Tool.Driver.Name != nil {
		groupName = run.Tool.Driver.Name
	}

	if groupName != nil {
		if err := r.StartGroup(*groupName); err != nil {
			return 0, err
		}
	}

	count, err := r.printResults(run.Results)
	if err != nil {
		return 0, err
	}

	if groupName != nil {
		if err := r.FinishGroup(); err != nil {
			return 0, err
		}
	}

	return count, nil
}

func (r *SimpleSarifReporter) printResults(results []*sarif.Result) (int, error) {
	count := 0

	for _, result := range results {
		count++

		message := ""

		if result.Message != nil && result.Message.Text != nil {
			message = *result.Message.Text
		}

		// Use last location
		location := result.Locations[len(result.Locations)-1]

		filename := internal.DefaultValue(
			location.PhysicalLocation.ArtifactLocation.URI,
			"",
		)
		line := internal.DefaultValue(location.PhysicalLocation.Region.StartLine, 0)

		if err := r.WriteLine("%s:%d:%s", filename, line, message); err != nil {
			return 0, err
		}
	}

	return count, nil
}

// FileSarifReporter reports results in a SARIF format.
type FileSarifReporter struct {
	output io.Writer
	SarifReporter
}

// NewFileSarifReporter creates a new FileSarifReporter.
func NewFileSarifReporter(output io.Writer) *FileSarifReporter {
	return &FileSarifReporter{output: output}
}

func (r *FileSarifReporter) Finish() error {
	return r.Write(r.output)
}
