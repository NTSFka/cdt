package reporter

import (
	"fmt"
	"io"
)

// LineReporter reports results in a simple line-oriented format.
type LineReporter struct {
	Output io.Writer
}

func (r *LineReporter) StartGroup(name string) error {
	return r.WriteLine("[%s]", name)
}

func (r *LineReporter) FinishGroup() error {
	return r.WriteLine("")
}

func (r *LineReporter) WriteLine(format string, args ...any) error {
	if _, err := fmt.Fprintf(r.Output, format+"\n", args...); err != nil {
		return err
	}

	return nil
}
