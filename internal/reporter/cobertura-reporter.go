package reporter

import (
	"cdt/internal/utils/cobertura"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
)

// CoberturaReporter reports results in a Cobertura format.
type CoberturaReporter struct {
	Coverage *cobertura.Coverage
}

func (r *CoberturaReporter) AddCobertura(report *cobertura.Coverage) error {
	if report == nil {
		return fmt.Errorf("report cannot be nil")
	}

	if r.Coverage == nil {
		r.Coverage = report

		return nil
	}

	// TODO: merge
	return errors.New("merging Cobertura reports is not supported")
}

func (r *CoberturaReporter) Write(output io.Writer) error {
	if r.Coverage == nil {
		return nil
	}

	encoder := xml.NewEncoder(output)
	defer func() {
		_ = encoder.Close()
	}()

	if err := encoder.Encode(r.Coverage); err != nil {
		return err
	}

	return nil
}

// SimpleCoberturaReporter reports CRTF results in a simple format.
type SimpleCoberturaReporter struct {
	LineReporter
	CoberturaReporter
}

// NewSimpleCoberturaReporter creates a new SimpleCoberturaReporter.
func NewSimpleCoberturaReporter(output io.Writer) *SimpleCoberturaReporter {
	return &SimpleCoberturaReporter{
		LineReporter: LineReporter{Output: output},
	}
}

func (r *SimpleCoberturaReporter) Finish() error {
	if r.Coverage == nil {
		return nil
	}

	for _, pack := range r.Coverage.Packages {
		for _, class := range pack.Classes {
			if err := r.WriteLine("%s: %.2f%%", class.Filename, class.LineRate*100); err != nil {
				return err
			}
		}
	}

	return r.WriteLine("\nTotal: %.2f%%", r.Coverage.LineRate*100)
}

// FileCoberturaReporter reports results in a CTRF format.
type FileCoberturaReporter struct {
	output io.Writer
	CoberturaReporter
}

// NewFileCoberturaReporter creates a new FileCoberturaReporter.
func NewFileCoberturaReporter(output io.Writer) *FileCoberturaReporter {
	return &FileCoberturaReporter{output: output}
}

func (r *FileCoberturaReporter) Finish() error {
	return r.Write(r.output)
}
