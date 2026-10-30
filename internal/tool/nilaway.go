package tool

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"cdt/internal"

	"github.com/owenrumney/go-sarif/v3/pkg/report/v210/sarif"
)

const IdNilAway = "nilaway"

// A NilAway is a tool that wraps golang main tool `nilaway`.
type NilAway struct {
	internal.ExecutableTool
}

// NewNilAway creates a go tool from a custom executable.
func NewNilAway(detect internal.ExecutableToolDetectFunc) *NilAway {
	return &NilAway{
		internal.MakeExecutableTool(
			IdNilAway,
			"NilAway",
			"NilAway is a static analysis tool that seeks to help developers avoid nil panics in production.",
			internal.Tags{internal.ToolTagGo, internal.ToolTagAnalyze},
			detect,
		),
	}
}

// DetectNilAway create nilaway tool can be used in the project.
func DetectNilAway(
	ctx context.Context,
	options internal.DetectOptions,
) *NilAway {
	return NewNilAway(func() (*internal.Executable, error) {
		return options.Environment.FindExecutable(ctx, options.GetToolPath(IdNilAway, "nilaway"))
	})
}

func (n *NilAway) AnalyzeFiles(
	ctx context.Context,
	options internal.ProjectAnalyzerOptions,
) error {
	args := appendFiles(options.ExtraArgs, options.Filenames, new("./..."))

	if options.Output.Reporter != nil {
		return n.analyzeFilesReporter(ctx, options, options.Output.Reporter, args)
	}

	return n.RunForProjectWithWriter(ctx, options.ProjectInfo, options.Output.Output, args)
}

func (n *NilAway) analyzeFilesReporter(
	ctx context.Context,
	options internal.ProjectAnalyzerOptions,
	reporter internal.FileIssueReporter,
	args []string,
) error {
	// Process JSON output
	args = append([]string{"-json"}, args...)

	outputBuffer := bytes.Buffer{}

	err := n.RunForProjectWithWriter(ctx, options.ProjectInfo, &outputBuffer, args)

	if err != nil && !checkExitCode(err, []int{3}) {
		return err
	} else if outputBuffer.Len() == 0 {
		// Nothing returned
		return nil
	}

	return errors.Join(err, n.reportJsonOutput(&outputBuffer, reporter))
}

type nilAwayJsonReport map[string]struct {
	NilAway []struct {
		Posn    string `json:"posn"`
		End     string `json:"end"`
		Message string `json:"message"`
	} `json:"nilaway"`
}

func (n *NilAway) reportJsonOutput(
	reader io.Reader,
	analyzerReporter internal.FileIssueReporter,
) error {
	var report nilAwayJsonReport
	if err := json.NewDecoder(reader).Decode(&report); err != nil {
		return err
	}

	sarifReport := sarif.NewReport()

	run := sarif.NewRunWithInformationURI(
		n.Name(),
		"https://github.com/uber-go/nilaway",
	)

	for _, group := range report {
		for _, reportRecord := range group.NilAway {
			position := strings.SplitN(reportRecord.Posn, ":", 3)

			if len(position) != 3 {
				return fmt.Errorf("invalid position: %s", reportRecord.Posn)
			}

			line, _ := strconv.Atoi(position[1])
			column, _ := strconv.Atoi(position[2])

			run.CreateResultForRule("").
				WithMessage(sarif.NewTextMessage(reportRecord.Message)).
				AddLocation(sarif.NewLocationWithPhysicalLocation(
					sarif.NewPhysicalLocation().
						WithArtifactLocation(sarif.NewSimpleArtifactLocation(position[0])).
						WithRegion(sarif.NewRegion().WithStartLine(line).WithStartColumn(column)),
				))
		}
	}

	sarifReport.AddRun(run)

	return analyzerReporter.AddSarif(sarifReport)
}
