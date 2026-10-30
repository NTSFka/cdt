package reporter_test

import (
	"bytes"
	"cdt/internal/reporter"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLineReporter_AddSarif_Empty(t *testing.T) {
	buffer := bytes.Buffer{}

	_ = reporter.LineReporter{Output: &buffer}

	require.Empty(t, buffer.String())
}

func TestLineReporter_ReportLine(t *testing.T) {
	buffer := bytes.Buffer{}

	report := reporter.LineReporter{Output: &buffer}

	assert.NoError(t, report.WriteLine("%s:%d:%s", "file1.go", 10, "failed 1"))
	assert.NoError(t, report.WriteLine("%s:%d:%s", "file2.go", 10, "failed 2"))

	require.Equal(
		t,
		"file1.go:10:failed 1\nfile2.go:10:failed 2\n",
		buffer.String(),
	)
}

type errWriter struct{}

func (errWriter) Write(p []byte) (int, error) {
	return 0, errors.New("write failed")
}

func TestLineReporter_ReportLine_Error(t *testing.T) {
	report := reporter.LineReporter{Output: &errWriter{}}

	assert.Error(t, report.WriteLine("%s:%d:%s", "file1.go", 10, "failed 1"))
}

func TestLineReporter_Group(t *testing.T) {
	buffer := bytes.Buffer{}

	report := reporter.LineReporter{Output: &buffer}

	assert.NoError(t, report.StartGroup("group1"))
	assert.NoError(t, report.WriteLine("%s:%d:%s", "file1.go", 10, "failed 1"))
	assert.NoError(t, report.WriteLine("%s:%d:%s", "file2.go", 10, "failed 2"))
	assert.NoError(t, report.FinishGroup())

	require.Equal(
		t,
		"[group1]\nfile1.go:10:failed 1\nfile2.go:10:failed 2\n\n",
		buffer.String(),
	)
}
