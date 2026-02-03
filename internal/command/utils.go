package command

import (
	"cdt/internal"
	"strings"
)

// ParseOptionOutput parses the output string into format and filename components.
func ParseOptionOutput[T ~string](output string, def T) internal.OutputOptions[T] {
	if len(output) == 0 {
		return internal.OutputOptions[T]{Format: def}
	}

	parts := strings.SplitN(output, ":", 2)
	result := internal.OutputOptions[T]{
		Format: T(parts[0]),
	}

	if len(parts) == 2 {
		result.Filename = &parts[1]
	}

	return result
}
