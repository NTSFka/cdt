package workflow

import (
	"cdt/internal"
)

// SupportedTypes stores supported workflow types.
var SupportedTypes = []Type{
	&Go{},
	&CMake{},
	&PHP{},
	&Python{},
}

// Config specifies a workflow configuration.
type Config struct {
	// Directory is the root directory of the project
	Directory string

	// IntermediateDirectory is the directory for the project's intermediate files
	IntermediateDirectory *string
}

// Type specifies a predefined workflow type.
type Type interface {
	// Id returns workflow type unique identifier.
	Id() string

	// Detect detects if a workflow of this type is in a given directory.
	Detect(directory string) bool

	// Create a workflow from a directory with a given tool set.
	Create(config Config, tools internal.Tools) internal.Project
}
