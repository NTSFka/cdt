package reporter

import (
	"fmt"
	"io"

	"github.com/ctrf-io/go-ctrf-json-reporter/ctrf"
)

// CtrfReporter reports results in a CTRF format.
type CtrfReporter struct {
	Reports []*ctrf.Report
}

func (r *CtrfReporter) AddCtrf(report *ctrf.Report) error {
	if report == nil {
		return fmt.Errorf("report cannot be nil")
	}

	r.Reports = append(r.Reports, report)

	return nil
}

func (r *CtrfReporter) Write(output io.Writer) error {
	// The reports cannot be merged into one report, so it needs to be separate JSONs
	for _, report := range r.Reports {
		if err := report.Write(output, false); err != nil {
			return err
		}
	}

	return nil
}

// SimpleCtrfReporter reports CRTF results in a simple format.
type SimpleCtrfReporter struct {
	LineReporter
	CtrfReporter
}

// NewSimpleCtrfReporter creates a new SimpleCtrfReporter.
func NewSimpleCtrfReporter(output io.Writer) *SimpleCtrfReporter {
	return &SimpleCtrfReporter{
		LineReporter: LineReporter{Output: output},
	}
}

func (r *SimpleCtrfReporter) Finish() error {
	failed := 0
	passed := 0
	skipped := 0
	pending := 0
	other := 0

	for _, report := range r.Reports {
		if report.Results != nil {
			for _, test := range report.Results.Tests {
				switch test.Status {
				case ctrf.TestPassed:
					passed++
				case ctrf.TestSkipped:
					skipped++
				case ctrf.TestFailed:
					failed++

					if err := r.writeFailed(test); err != nil {
						return err
					}
				case ctrf.TestPending:
					pending++
				case ctrf.TestOther:
					other++
				}
			}
		}
	}

	_, err := fmt.Fprintf(
		r.Output,
		"\nPassed: %d, Failed: %d, Skipped: %d, Pending: %d, Other: %d\n",
		passed,
		failed,
		skipped,
		pending,
		other,
	)

	return err
}

func (r *SimpleCtrfReporter) writeFailed(test *ctrf.TestResult) error {
	return r.WriteLine("%s [%s] %s", test.Filepath, test.Name, test.Message)
}

// FileCtrfReporter reports results in a CTRF format.
type FileCtrfReporter struct {
	output io.Writer
	CtrfReporter
}

// NewFileCtrfReporter creates a new FileCtrfReporter.
func NewFileCtrfReporter(output io.Writer) *FileCtrfReporter {
	return &FileCtrfReporter{output: output}
}

func (r *FileCtrfReporter) Finish() error {
	return r.Write(r.output)
}
