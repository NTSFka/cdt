package tool

import (
	"cdt/internal"
)

// InitTools initializes all supported tools for a given environment
func InitTools(environment internal.Environment) internal.Tools {
	return internal.Tools{
		DetectClangFormat(environment),
		DetectClangTidy(environment),
		DetectCMake(environment),
		DetectCTest(environment),
		DetectGo(environment),
		DetectGolangCILint(environment),
		DetectDocker(environment),
		DetectDockerCompose(environment),
		DetectPHP(environment),
		DetectPHPUnit(environment),
		DetectParaTest(environment),
		DetectPHPStan(environment),
		DetectPHPCSFixer(environment),
		DetectComposer(environment),
		DetectPython(environment),
	}
}

// InitEnvironmentProviders initializes environment providers
func InitEnvironmentProviders(environment internal.Environment) internal.EnvironmentProviders {
	return internal.EnvironmentProviders{
		internal.SystemEnvironmentProvider,
		DetectDocker(environment),
		DetectDockerCompose(environment),
		DetectPython(environment),
	}
}
