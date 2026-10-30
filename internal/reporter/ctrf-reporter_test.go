package reporter_test

import (
	"bytes"
	"cdt/internal/reporter"
	"fmt"
	"testing"

	"github.com/ctrf-io/go-ctrf-json-reporter/ctrf"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCtrfReporter_AddCtrf(t *testing.T) {
	report := reporter.CtrfReporter{}

	ctrfReport := ctrf.NewReport("tool1", nil)
	require.NoError(t, report.AddCtrf(ctrfReport))
	require.NoError(t, report.AddCtrf(ctrfReport))

	assert.Len(t, report.Reports, 2)
}

func TestCtrfReporter_AddCtrf_ErrorNil(t *testing.T) {
	report := reporter.CtrfReporter{}

	require.Error(t, report.AddCtrf(nil))
}

func TestCtrfReporter_Write(t *testing.T) {
	report := reporter.CtrfReporter{}

	ctrfReport := ctrf.NewReport("tool1", nil)
	require.NoError(t, report.AddCtrf(ctrfReport))

	buffer := bytes.Buffer{}

	require.NoError(t, report.Write(&buffer))

	assert.NotEmpty(t, buffer.String())
}

func TestCtrfReporter_Write_Empty(t *testing.T) {
	report := reporter.CtrfReporter{}

	buffer := bytes.Buffer{}

	require.NoError(t, report.Write(&buffer))

	assert.Empty(t, buffer.String())
}

func TestSimpleCtrfReporter_Finish_NoReport(t *testing.T) {
	buffer := bytes.Buffer{}

	report := reporter.NewSimpleCtrfReporter(&buffer)

	// Flush result
	require.NoError(t, report.Finish())

	require.Equal(t, "\nPassed: 0, Failed: 0, Skipped: 0, Pending: 0, Other: 0\n", buffer.String())
}

func TestSimpleCtrfReporter_Finish_EmptyReport(t *testing.T) {
	buffer := bytes.Buffer{}

	report := reporter.NewSimpleCtrfReporter(&buffer)

	ctrfReport := ctrf.NewReport("tool1", nil)
	require.NoError(t, report.AddCtrf(ctrfReport))

	// Flush result
	require.NoError(t, report.Finish())

	require.Equal(t, "\nPassed: 0, Failed: 0, Skipped: 0, Pending: 0, Other: 0\n", buffer.String())
}

func TestSimpleCtrfReporter_Finish_NonEmptyReport(t *testing.T) {
	buffer := bytes.Buffer{}

	report := reporter.NewSimpleCtrfReporter(&buffer)

	ctrfReport := ctrf.NewReport("tool1", nil)
	ctrfReport.Results = &ctrf.Results{
		Tests: []*ctrf.TestResult{
			{
				Name:     "Test1",
				Status:   ctrf.TestPassed,
				Message:  "Ok",
				Filepath: "file1.go",
			},
			{
				Name:     "Test2",
				Status:   ctrf.TestFailed,
				Message:  "Fail",
				Filepath: "file2.go",
			},
			{
				Name:     "Test3",
				Status:   ctrf.TestSkipped,
				Message:  "Skipped",
				Filepath: "file3.go",
			},
			{
				Name:     "Test4",
				Status:   ctrf.TestOther,
				Message:  "Other",
				Filepath: "file4.go",
			},
			{
				Name:     "Test5",
				Status:   ctrf.TestPending,
				Message:  "Pending",
				Filepath: "file5.go",
			},
		},
	}
	require.NoError(t, report.AddCtrf(ctrfReport))

	// Flush result
	require.NoError(t, report.Finish())

	require.Equal(
		t,
		"file2.go [Test2] Fail\n\nPassed: 1, Failed: 1, Skipped: 1, Pending: 1, Other: 1\n",
		buffer.String(),
	)
}

func TestSimpleCtrfReporter_Finish_NonEmptyReport_Fail(t *testing.T) {
	// Test write failure for different buffer sizes
	for _, limit := range []int{0, 16, 49} {
		t.Run(fmt.Sprintf("%d", limit), func(t *testing.T) {
			report := reporter.NewSimpleCtrfReporter(&limitedErrWriter{limit: limit})

			ctrfReport := ctrf.NewReport("tool1", nil)
			require.NoError(t, report.AddCtrf(ctrfReport))

			// Flush result
			require.Error(t, report.Finish())
		})
	}
}

func TestFileCtrfReporter_Finish(t *testing.T) {
	buffer := bytes.Buffer{}

	report := reporter.NewFileCtrfReporter(&buffer)

	ctrfReport := ctrf.NewReport("tool1", nil)
	require.NoError(t, report.AddCtrf(ctrfReport))

	// Flush result
	require.NoError(t, report.Finish())

	require.NotEmpty(t, buffer.String())
}

func TestFileCtrfReporter_Finish_NonEmptyReport_Fail(t *testing.T) {
	// Test write failure for different buffer sizes
	for _, limit := range []int{0, 16, 49} {
		t.Run(fmt.Sprintf("%d", limit), func(t *testing.T) {
			report := reporter.NewFileCtrfReporter(&limitedErrWriter{limit: limit})

			ctrfReport := ctrf.NewReport("tool1", nil)
			require.NoError(t, report.AddCtrf(ctrfReport))

			// Flush result
			require.Error(t, report.Finish())
		})
	}
}
