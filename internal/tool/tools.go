package tool

import (
	"context"

	"cdt/internal"
)

// DetectOptions is a set of options for tool(s) detection.
type DetectOptions struct {
	// ProjectDirectory is a path to the project directory.
	ProjectDirectory string
	// Environment is an environment where tools should be look for.
	Environment internal.Environment
	// ToolsPaths is a mapping of tool IDs to executable paths.
	ToolsPaths map[string]string
}

// GetToolPath returns a tool executable path by its ID.
func (c DetectOptions) GetToolPath(id string, def string) string {
	if path, ok := c.ToolsPaths[id]; ok {
		return path
	} else {
		return def
	}
}

// InitTools initializes all supported tools for a given environment.
func InitTools(
	ctx context.Context,
	options DetectOptions,
) internal.Tools {
	return internal.Tools{
		// C/C++
		DetectClangFormat(ctx, options),
		DetectClangTidy(ctx, options),
		DetectCMake(ctx, options),
		DetectCTest(ctx, options),

		// Go
		DetectGo(ctx, options),
		DetectGolangCILint(ctx, options),
		DetectNilAway(ctx, options),

		// PHP
		DetectPHP(ctx, options),
		DetectPHPUnit(ctx, options),
		DetectParaTest(ctx, options),
		DetectPHPStan(ctx, options),
		DetectPHPCSFixer(ctx, options),
		DetectComposer(ctx, options),

		// Python
		DetectPython(ctx, options),
		DetectPyTest(ctx, options),
		DetectPip(ctx, options),
		DetectPylint(ctx, options),
		DetectFlake8(ctx, options),
		DetectMyPy(ctx, options),
		DetectRuff(ctx, options),
		DetectBandit(ctx, options),
		DetectBlack(ctx, options),

		// Other
		DetectDocker(ctx, options),
		DetectDockerCompose(ctx, options),
	}
}

// InitEnvironmentProviders initializes environment providers.
func InitEnvironmentProviders(
	ctx context.Context,
	options DetectOptions,
) internal.EnvironmentProviders {
	return internal.EnvironmentProviders{
		internal.SystemEnvironmentProvider,
		DetectDocker(ctx, options),
		DetectDockerCompose(ctx, options),
		DetectPython(ctx, options),
	}
}
