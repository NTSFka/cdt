package command

import (
	"cdt/internal"
	"fmt"
	"io"
	"os"
	"strings"
)

// ParsedOutputOptions is the result of ParseOptionOutput.
type ParsedOutputOptions struct {
	Format   string
	Filename *string
}

// ParseOptionOutput parses the output string into format and filename components.
func ParseOptionOutput(output string) ParsedOutputOptions {
	if len(output) == 0 {
		return ParsedOutputOptions{}
	}

	parts := strings.SplitN(output, ":", 2)
	result := ParsedOutputOptions{
		Format: parts[0], //nolint: nilaway
	}

	if len(parts) == 2 {
		result.Filename = &parts[1]
	}

	return result
}

type noopCloser struct {
	writer io.WriteCloser
}

func (n *noopCloser) Write(p []byte) (int, error) {
	return n.writer.Write(p)
}

// gofix: inline
func (n *noopCloser) Close() error {
	return nil
}

func buildOutput(filename *string) (io.WriteCloser, error) {
	if filename != nil {
		file, err := os.OpenFile(*filename, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to create output file: %w", err)
		}

		return file, nil
	}

	return &noopCloser{os.Stdout}, nil
}

func getTool[T any](toolName string, tools internal.Tools, kind string) (T, error) {
	var none T

	tool := tools.Get(toolName)
	if tool == nil {
		return none, fmt.Errorf("tool '%s' not found", toolName)
	}

	linterTool, ok := tool.(T)

	if ok {
		return linterTool, nil
	}

	return none, fmt.Errorf("tool '%s' doesn't support %s", toolName, kind)
}
