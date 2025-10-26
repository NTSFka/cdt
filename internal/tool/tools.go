package tool

import (
	"context"

	"cdt/internal"
)

// InitTools initializes all supported tools for a given environment.
func InitTools(
	ctx context.Context,
	config internal.ConfigTools,
	environment internal.Environment,
) internal.Tools {
	return internal.Tools{
		// C/C++
		DetectClangFormat(ctx, config, environment),
		DetectClangTidy(ctx, config, environment),
		DetectCMake(ctx, config, environment),
		DetectCTest(ctx, config, environment),

		// Go
		DetectGo(ctx, config, environment),
		DetectGolangCILint(ctx, config, environment),
		DetectNilAway(ctx, config, environment),

		// PHP
		DetectPHP(ctx, config, environment),
		DetectPHPUnit(ctx, config, environment),
		DetectParaTest(ctx, config, environment),
		DetectPHPStan(ctx, config, environment),
		DetectPHPCSFixer(ctx, config, environment),
		DetectComposer(ctx, config, environment),

		// Python
		DetectPython(ctx, config, environment),
		DetectPyTest(ctx, config, environment),
		DetectPip(ctx, config, environment),
		DetectPylint(ctx, config, environment),
		DetectFlake8(ctx, config, environment),
		DetectMyPy(ctx, config, environment),
		DetectRuff(ctx, config, environment),
		DetectBandit(ctx, config, environment),
		DetectBlack(ctx, config, environment),

		// Other
		DetectDocker(ctx, config, environment),
		DetectDockerCompose(ctx, config, environment),
	}
}

// InitEnvironmentProviders initializes environment providers.
func InitEnvironmentProviders(
	ctx context.Context,
	config internal.ConfigTools,
	environment internal.Environment,
) internal.EnvironmentProviders {
	return internal.EnvironmentProviders{
		internal.SystemEnvironmentProvider,
		DetectDocker(ctx, config, environment),
		DetectDockerCompose(ctx, config, environment),
		DetectPython(ctx, config, environment),
	}
}
