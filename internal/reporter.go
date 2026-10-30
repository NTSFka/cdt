package internal

import (
	"cdt/internal/utils/cobertura"

	"github.com/ctrf-io/go-ctrf-json-reporter/ctrf"
	"github.com/owenrumney/go-sarif/v3/pkg/report/v210/sarif"
)

// A Reporter is a basic interface for processing reports.
type Reporter interface {
	// Finish is called when the report is complete.
	Finish() error
}

// A FileIssueReporter is an interface for processing file issue reports.
type FileIssueReporter interface {
	Reporter

	// AddSarif adds a report in SARIF format.
	AddSarif(report *sarif.Report) error
}

// An TestsReporter is an interface for reporting test results.
type TestsReporter interface {
	Reporter

	// AddCtrf adds a report in CTRF format.
	AddCtrf(report *ctrf.Report) error
}

// An CoverageReporter is an interface for reporting code coverage.
type CoverageReporter interface {
	Reporter

	// AddCobertura adds coverage in Cobertura format.
	AddCobertura(cobertura *cobertura.Coverage) error
}
