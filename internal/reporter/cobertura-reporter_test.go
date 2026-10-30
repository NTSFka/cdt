package reporter_test

import (
	"bytes"
	"cdt/internal/reporter"
	"cdt/internal/utils/cobertura"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCoberturaReporter_AddCobertura(t *testing.T) {
	report := reporter.CoberturaReporter{}

	coverage := &cobertura.Coverage{}
	require.NoError(t, report.AddCobertura(coverage))

	assert.NotNil(t, report.Coverage)
}

func TestCoberturaReporter_AddCobertura_ErrorNil(t *testing.T) {
	report := reporter.CoberturaReporter{}

	require.Error(t, report.AddCobertura(nil))
}

func TestCoberturaReporter_Write(t *testing.T) {
	report := reporter.CoberturaReporter{}

	coverage := &cobertura.Coverage{}
	require.NoError(t, report.AddCobertura(coverage))

	buffer := bytes.Buffer{}

	require.NoError(t, report.Write(&buffer))

	assert.NotEmpty(t, buffer.String())
}

func TestCoberturaReporter_Write_Empty(t *testing.T) {
	report := reporter.CoberturaReporter{}

	buffer := bytes.Buffer{}

	require.NoError(t, report.Write(&buffer))

	assert.Empty(t, buffer.String())
}

func TestSimpleCoberturaReporter_Finish_NoReport(t *testing.T) {
	buffer := bytes.Buffer{}

	report := reporter.NewSimpleCoberturaReporter(&buffer)

	// Flush result
	require.NoError(t, report.Finish())

	require.Empty(t, buffer.String())
}

func TestSimpleCoberturaReporter_Finish_EmptyReport(t *testing.T) {
	buffer := bytes.Buffer{}

	report := reporter.NewSimpleCoberturaReporter(&buffer)

	coverage := &cobertura.Coverage{}
	require.NoError(t, report.AddCobertura(coverage))

	// Flush result
	require.NoError(t, report.Finish())

	require.Equal(t, "\nTotal: 0.00%\n", buffer.String())
}

func TestSimpleCoberturaReporter_Finish_NonEmptyReport(t *testing.T) {
	buffer := bytes.Buffer{}

	report := reporter.NewSimpleCoberturaReporter(&buffer)

	coverage := &cobertura.Coverage{
		LineRate: 0.2,
		Packages: []*cobertura.Package{
			{
				Classes: []*cobertura.Class{
					{
						Name:     "Class1",
						Filename: "file1.go",
						LineRate: 0.59,
					},
				},
			},
		},
	}
	require.NoError(t, report.AddCobertura(coverage))

	// Flush result
	require.NoError(t, report.Finish())

	require.Equal(
		t,
		"file1.go: 59.00%\n\nTotal: 20.00%\n",
		buffer.String(),
	)
}

func TestSimpleCoberturaReporter_Finish_NonEmptyReport_Fail(t *testing.T) {
	// Test write failure for different buffer sizes
	for _, limit := range []int{0, 16} {
		t.Run(fmt.Sprintf("%d", limit), func(t *testing.T) {
			report := reporter.NewSimpleCoberturaReporter(&limitedErrWriter{limit: limit})

			coverage := &cobertura.Coverage{
				LineRate: 0.2,
				Packages: []*cobertura.Package{
					{
						Classes: []*cobertura.Class{
							{
								Name:     "Class1",
								Filename: "file1.go",
								LineRate: 0.59,
							},
						},
					},
				},
			}
			require.NoError(t, report.AddCobertura(coverage))

			// Flush result
			require.Error(t, report.Finish())
		})
	}
}

func TestFileCoberturaReporter_Finish(t *testing.T) {
	buffer := bytes.Buffer{}

	report := reporter.NewFileCoberturaReporter(&buffer)

	coverage := &cobertura.Coverage{}
	require.NoError(t, report.AddCobertura(coverage))

	// Flush result
	require.NoError(t, report.Finish())

	require.NotEmpty(t, buffer.String())
}

func TestFileCoberturaReporter_Finish_Empty(t *testing.T) {
	buffer := bytes.Buffer{}

	report := reporter.NewFileCoberturaReporter(&buffer)

	// Flush result
	require.NoError(t, report.Finish())

	require.Empty(t, buffer.String())
}

func TestFileCoberturaReporter_Finish_NonEmptyReport_Fail(t *testing.T) {
	// Test write failure for different buffer sizes
	for _, limit := range []int{0, 16, 49} {
		t.Run(fmt.Sprintf("%d", limit), func(t *testing.T) {
			report := reporter.NewFileCoberturaReporter(&limitedErrWriter{limit: limit})

			coverage := &cobertura.Coverage{}
			require.NoError(t, report.AddCobertura(coverage))

			// Flush result
			require.Error(t, report.Finish())
		})
	}
}
