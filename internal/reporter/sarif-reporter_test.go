package reporter_test

import (
	"bytes"
	"cdt/internal/reporter"
	"errors"
	"fmt"
	"testing"

	"github.com/owenrumney/go-sarif/v3/pkg/report/v210/sarif"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSarifReporter_AddSarif(t *testing.T) {
	report := reporter.SarifReporter{}

	sarifReport := sarif.NewReport()
	require.NoError(t, report.AddSarif(sarifReport))
	require.NoError(t, report.AddSarif(sarifReport))

	assert.Equal(t, sarifReport.Version, report.Report.Version)
}

func TestSarifReporter_AddSarif_ErrorNil(t *testing.T) {
	report := reporter.SarifReporter{}

	require.Error(t, report.AddSarif(nil))
}

func TestSarifReporter_Write(t *testing.T) {
	report := reporter.SarifReporter{}

	sarifReport := sarif.NewReport()
	require.NoError(t, report.AddSarif(sarifReport))

	buffer := bytes.Buffer{}

	require.NoError(t, report.Write(&buffer))

	assert.NotEmpty(t, buffer.String())
}

func TestSarifReporter_Write_Empty(t *testing.T) {
	report := reporter.SarifReporter{}

	buffer := bytes.Buffer{}

	require.NoError(t, report.Write(&buffer))

	assert.Empty(t, buffer.String())
}

func TestSimpleSarifReporter_Finish_NoReport(t *testing.T) {
	buffer := bytes.Buffer{}

	report := reporter.NewSimpleSarifReporter(&buffer)

	// Flush result
	require.NoError(t, report.Finish())

	require.Equal(t, "\nTotal issues: 0\n\n", buffer.String())
}

func TestSimpleSarifReporter_Finish_EmptyReport(t *testing.T) {
	buffer := bytes.Buffer{}

	report := reporter.NewSimpleSarifReporter(&buffer)

	sarifReport := sarif.NewReport()
	require.NoError(t, report.AddSarif(sarifReport))

	// Flush result
	require.NoError(t, report.Finish())

	require.Equal(t, "\nTotal issues: 0\n\n", buffer.String())
}

func TestSimpleSarifReporter_Finish_EmptyRun(t *testing.T) {
	buffer := bytes.Buffer{}

	report := reporter.NewSimpleSarifReporter(&buffer)

	sarifReport := sarif.NewReport()
	run := sarif.NewRun()
	sarifReport.AddRun(run)

	require.NoError(t, report.AddSarif(sarifReport))

	// Flush result
	require.NoError(t, report.Finish())

	require.Equal(t, "\nTotal issues: 0\n\n", buffer.String())
}

func TestSimpleSarifReporter_Finish_NonEmptyReport(t *testing.T) {
	buffer := bytes.Buffer{}

	report := reporter.NewSimpleSarifReporter(&buffer)

	sarifReport := buildSarifReport()

	require.NoError(t, report.AddSarif(sarifReport))

	// Flush result
	require.NoError(t, report.Finish())

	require.Equal(
		t,
		"[tool]\nfile1.go:10:failed 1\nfile2.go:10:failed 2\n\n\nTotal issues: 2\n\n",
		buffer.String(),
	)
}

type limitedErrWriter struct {
	limit int
	used  int
}

func (w *limitedErrWriter) Write(data []byte) (int, error) {
	if w.used+len(data) > w.limit {
		return 0, errors.New("buffer full")
	}

	w.used += len(data)

	return len(data), nil
}

func TestSimpleSarifReporter_Finish_NonEmptyReport_Fail(t *testing.T) {
	// Test write failure for different buffer sizes
	for _, limit := range []int{0, 16, 49, 64} {
		t.Run(fmt.Sprintf("%d", limit), func(t *testing.T) {
			report := reporter.NewSimpleSarifReporter(&limitedErrWriter{limit: limit})

			sarifReport := buildSarifReport()

			require.NoError(t, report.AddSarif(sarifReport))

			// Flush result
			require.Error(t, report.Finish())
		})
	}
}

func TestFileSarifReporter_Finish(t *testing.T) {
	buffer := bytes.Buffer{}

	report := reporter.NewFileSarifReporter(&buffer)

	sarifReport := buildSarifReport()

	require.NoError(t, report.AddSarif(sarifReport))

	// Flush result
	require.NoError(t, report.Finish())

	require.NotEmpty(t, buffer.String())
}

func buildSarifReport() *sarif.Report {
	sarifReport := sarif.NewReport()
	run := sarif.NewRun()

	run.Tool = &sarif.Tool{
		Driver: &sarif.ToolComponent{
			Name: new("tool"),
		},
	}

	run.CreateResultForRule("").
		WithMessage(sarif.NewTextMessage("failed 1")).
		AddLocation(sarif.NewLocation().WithPhysicalLocation(
			sarif.NewPhysicalLocation().
				WithArtifactLocation(sarif.NewArtifactLocation().WithURI("file1.go")).
				WithRegion(sarif.NewRegion().WithStartLine(10)),
		))

	run.CreateResultForRule("").
		WithMessage(sarif.NewTextMessage("failed 2")).
		AddLocation(sarif.NewLocation().WithPhysicalLocation(
			sarif.NewPhysicalLocation().
				WithArtifactLocation(sarif.NewArtifactLocation().WithURI("file2.go")).
				WithRegion(sarif.NewRegion().WithStartLine(10)),
		))

	sarifReport.AddRun(run)

	return sarifReport
}
