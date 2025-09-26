package tool

import (
	"cdt/internal"
	"context"
)

// InitTools initializes all supported tools for a given environment.
func InitTools(ctx context.Context, environment internal.Environment) internal.Tools {
	return internal.Tools{
		// C/C++
		DetectClangFormat(ctx, environment),
		DetectClangTidy(ctx, environment),
		DetectCMake(ctx, environment),
		DetectCTest(ctx, environment),

		// Go
		DetectGo(ctx, environment),
		DetectGolangCILint(ctx, environment),
		DetectNilAway(ctx, environment),

		// PHP
		DetectPHP(ctx, environment),
		DetectPHPUnit(ctx, environment),
		DetectParaTest(ctx, environment),
		DetectPHPStan(ctx, environment),
		DetectPHPCSFixer(ctx, environment),
		DetectComposer(ctx, environment),

		// Python
		DetectPython(ctx, environment),
		DetectPyTest(ctx, environment),
		DetectPip(ctx, environment),
		DetectPylint(ctx, environment),
		DetectFlake8(ctx, environment),
		DetectMyPy(ctx, environment),
		DetectRuff(ctx, environment),
		DetectBandit(ctx, environment),
		DetectBlack(ctx, environment),

		// Other
		DetectDocker(ctx, environment),
		DetectDockerCompose(ctx, environment),
	}
}

// InitEnvironmentProviders initializes environment providers.
func InitEnvironmentProviders(ctx context.Context, environment internal.Environment) internal.EnvironmentProviders {
	return internal.EnvironmentProviders{
		internal.SystemEnvironmentProvider,
		DetectDocker(ctx, environment),
		DetectDockerCompose(ctx, environment),
		DetectPython(ctx, environment),
	}
}
