package pkg

import (
	"cdt/internal"
	"cdt/internal/tool"
	"context"
)

// InitTools initializes all supported tools for a given environment.
func InitTools(
	ctx context.Context,
	options internal.DetectOptions,
) internal.Tools {
	return internal.Tools{
		// C/C++
		tool.DetectClangFormat(ctx, options),
		tool.DetectClangTidy(ctx, options),
		tool.DetectCMake(ctx, options),
		tool.DetectCTest(ctx, options),

		// Go
		tool.DetectGo(ctx, options),
		tool.DetectGolangCILint(ctx, options),
		tool.DetectNilAway(ctx, options),

		// PHP
		tool.DetectPHP(ctx, options),
		tool.DetectPHPUnit(ctx, options),
		tool.DetectParaTest(ctx, options),
		tool.DetectPHPStan(ctx, options),
		tool.DetectPHPCSFixer(ctx, options),
		tool.DetectComposer(ctx, options),

		// Python
		tool.DetectPython(ctx, options),
		tool.DetectPyTest(ctx, options),
		tool.DetectPip(ctx, options),
		tool.DetectPylint(ctx, options),
		tool.DetectFlake8(ctx, options),
		tool.DetectMyPy(ctx, options),
		tool.DetectRuff(ctx, options),
		tool.DetectBandit(ctx, options),
		tool.DetectBlack(ctx, options),
		tool.DetectPythonCoverage(ctx, options),

		// Other
		tool.DetectDocker(ctx, options),
		tool.DetectDockerCompose(ctx, options),
	}
}

// InitEnvironmentProviders initializes environment providers.
func InitEnvironmentProviders(
	ctx context.Context,
	options internal.DetectOptions,
) internal.EnvironmentProviders {
	return internal.EnvironmentProviders{
		internal.SystemEnvironmentProvider,
		tool.DetectDocker(ctx, options),
		tool.DetectDockerCompose(ctx, options),
		tool.DetectPython(ctx, options),
	}
}
