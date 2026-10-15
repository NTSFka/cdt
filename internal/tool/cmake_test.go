package tool_test

import (
	"context"
	"errors"
	"testing"

	"cdt/internal"
	"cdt/internal/test"
	"cdt/internal/tool"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCMake_CMakeDetect(t *testing.T) {
	cmake := test.RunDetectToolLastFound(t, cmakeDetect, []string{"cmake"})

	assert.Equal(t, tool.IdCMake, cmake.Id())
}

func TestCMake_CMakeDetect_NotFound(t *testing.T) {
	cmake := test.RunDetectToolNotFound(t, cmakeDetect, []string{"cmake"})

	assert.Equal(t, tool.IdCMake, cmake.Id())
}

func TestCMake_CMakeDetect_Config(t *testing.T) {
	cmake := test.RunDetectToolConfig(t, cmakeDetect, "cmake")

	assert.Equal(t, tool.IdCMake, cmake.Id())
}

func TestCMake_Configure(t *testing.T) {
	buildDir := t.TempDir()

	options := internal.ProjectConfiguratorOptions{
		ProjectInfo: internal.ProjectInfo{
			Directory:       "project",
			OutputDirectory: &buildDir,
		},
	}

	test.InvokeConfigureSuccess(t, cmakeConfigurator, options, []string{
		"-DCMAKE_EXPORT_COMPILE_COMMANDS=ON",
		"-S", ".",
		"-B", buildDir,
	})
}

func TestCMake_Configure_Failed(t *testing.T) {
	buildDir := t.TempDir()

	options := internal.ProjectConfiguratorOptions{
		ProjectInfo: internal.ProjectInfo{
			Directory:       "project",
			OutputDirectory: &buildDir,
		},
	}

	test.InvokeConfigureFailed(t, cmakeConfigurator, options, []string{
		"-DCMAKE_EXPORT_COMPILE_COMMANDS=ON",
		"-S", ".",
		"-B", buildDir,
	})
}

func TestCMake_Structure_NoOutputDirectory(t *testing.T) {
	test.InvokeStructure(
		t,
		cmakeStructure,
		internal.ProjectInfo{
			Directory: "project",
		},
		func(exec *test.Executable, executablePath string) {
		},
		func(structure *internal.ProjectStructure, err error) {
			assert.Nil(t, structure)
			require.ErrorIs(t, err, internal.ErrNoOutputDirectory)
		},
	)
}

func TestCMake_Structure_ConfigureFailed(t *testing.T) {
	buildDir := t.TempDir()

	test.InvokeStructure(
		t,
		cmakeStructure,
		internal.ProjectInfo{
			Directory:       "project",
			OutputDirectory: &buildDir,
		},
		func(exec *test.Executable, executablePath string) {
			exec.OnRun(executablePath, []string{
				"-DCMAKE_EXPORT_COMPILE_COMMANDS=ON",
				"-S", ".",
				"-B", buildDir,
			}).
				Return(errors.New("failed"))
		},
		func(structure *internal.ProjectStructure, err error) {
			assert.Nil(t, structure)
			require.EqualError(t, err, "failed")
		},
	)
}

func TestCMake_BuildTargets_All(t *testing.T) {
	buildDir := t.TempDir()

	test.InvokeBuildTargets(
		t,
		cmakeBuilder,
		internal.ProjectBuilderOptions{
			ProjectInfo: internal.ProjectInfo{
				Directory:       "project",
				OutputDirectory: &buildDir,
			},
		},
		func(exec *test.Executable, executablePath string) {
			// Configure
			exec.OnRun(executablePath, []string{
				"-DCMAKE_EXPORT_COMPILE_COMMANDS=ON",
				"-S", ".",
				"-B", buildDir,
			}).
				Return(nil)

			// Build
			exec.OnRun(executablePath, []string{"--build", buildDir}).
				Return(nil)
		},
		func(err error) {
			require.NoError(t, err)
		},
	)
}

func TestCMake_BuildTargets_All_Failed(t *testing.T) {
	buildDir := t.TempDir()

	test.InvokeBuildTargets(
		t,
		cmakeBuilder,
		internal.ProjectBuilderOptions{
			ProjectInfo: internal.ProjectInfo{
				Directory:       "project",
				OutputDirectory: &buildDir,
			},
		},
		func(exec *test.Executable, executablePath string) {
			// Configure
			exec.OnRun(executablePath, []string{
				"-DCMAKE_EXPORT_COMPILE_COMMANDS=ON",
				"-S", ".",
				"-B", buildDir,
			}).
				Return(nil)

			// Build
			exec.OnRun(executablePath, []string{"--build", buildDir}).
				Return(errors.New("failed"))
		},
		func(err error) {
			require.EqualError(t, err, "failed")
		},
	)
}

func TestCMake_BuildTargets_All_ConfigureFailed(t *testing.T) {
	buildDir := t.TempDir()

	test.InvokeBuildTargets(
		t,
		cmakeBuilder,
		internal.ProjectBuilderOptions{
			ProjectInfo: internal.ProjectInfo{
				Directory:       "project",
				OutputDirectory: &buildDir,
			},
		},
		func(exec *test.Executable, executablePath string) {
			// Configure
			exec.OnRun(executablePath, []string{
				"-DCMAKE_EXPORT_COMPILE_COMMANDS=ON",
				"-S", ".",
				"-B", buildDir,
			}).
				Return(errors.New("failed"))
		},
		func(err error) {
			require.EqualError(t, err, "failed")
		},
	)
}

func TestCMake_BuildTargets(t *testing.T) {
	buildDir := t.TempDir()

	test.InvokeBuildTargets(
		t,
		cmakeBuilder,
		internal.ProjectBuilderOptions{
			ProjectInfo: internal.ProjectInfo{
				Directory:       "project",
				OutputDirectory: &buildDir,
			},
			Targets: &[]string{"target1", "target2"},
		},
		func(exec *test.Executable, executablePath string) {
			// Configure
			exec.OnRun(executablePath, []string{
				"-DCMAKE_EXPORT_COMPILE_COMMANDS=ON",
				"-S", ".",
				"-B", buildDir,
			}).
				Return(nil)

			// Build
			exec.OnRun(executablePath, []string{
				"--build", buildDir,
				"--target", "target1", "target2",
			}).
				Return(nil)
		},
		func(err error) {
			require.NoError(t, err)
		},
	)
}

func TestCMake_BuildTargets_Failed(t *testing.T) {
	buildDir := t.TempDir()

	test.InvokeBuildTargets(
		t,
		cmakeBuilder,
		internal.ProjectBuilderOptions{
			ProjectInfo: internal.ProjectInfo{
				Directory:       "project",
				OutputDirectory: &buildDir,
			},
			Targets: &[]string{"target1", "target2"},
		},
		func(exec *test.Executable, executablePath string) {
			// Configure
			exec.OnRun(executablePath, []string{
				"-DCMAKE_EXPORT_COMPILE_COMMANDS=ON",
				"-S", ".",
				"-B", buildDir,
			}).
				Return(nil)

			// Build
			exec.OnRun(executablePath, []string{
				"--build", buildDir,
				"--target", "target1", "target2",
			}).
				Return(errors.New("failed"))
		},
		func(err error) {
			require.EqualError(t, err, "failed")
		},
	)
}

func TestCMake_BuildTargets_ConfigureFailed(t *testing.T) {
	buildDir := t.TempDir()

	test.InvokeBuildTargets(
		t,
		cmakeBuilder,
		internal.ProjectBuilderOptions{
			ProjectInfo: internal.ProjectInfo{
				Directory:       "project",
				OutputDirectory: &buildDir,
			},
			Targets: &[]string{"target1", "target2"},
		},
		func(exec *test.Executable, executablePath string) {
			// Configure
			exec.OnRun(executablePath, []string{
				"-DCMAKE_EXPORT_COMPILE_COMMANDS=ON",
				"-S", ".",
				"-B", buildDir,
			}).
				Return(errors.New("failed"))
		},
		func(err error) {
			require.EqualError(t, err, "failed")
		},
	)
}

func TestCMake_RunTarget_Failed(t *testing.T) {
	buildDir := t.TempDir()

	options := internal.ProjectRunnerOptions{
		ProjectInfo: internal.ProjectInfo{
			OutputDirectory: &buildDir,
		},
	}

	test.InvokeRunTarget(
		t,
		cmakeRunner,
		options,
		"target1",
		func(exec *test.Executable, executablePath string) {
			// Configure
			exec.OnRun(executablePath, []string{
				"-DCMAKE_EXPORT_COMPILE_COMMANDS=ON",
				"-S", ".",
				"-B", buildDir,
			}).
				Return(nil)

			// Build
			exec.OnRun(executablePath, []string{
				"--build", buildDir,
				"--target", "target1",
			}).
				Return(errors.New("failed"))
		},
		func(err error) {
			require.EqualError(t, err, "failed")
		},
	)
}

func cmakeDetect(ctx context.Context, options internal.DetectOptions) internal.Tool {
	return tool.DetectCMake(ctx, options)
}

func cmakeStructure(
	executable func() (*internal.Executable, error),
) internal.ProjectStructureProvider {
	return tool.NewCMake(executable)
}

func cmakeConfigurator(
	executable func() (*internal.Executable, error),
) internal.ProjectConfigurator {
	return tool.NewCMake(executable)
}

func cmakeBuilder(executable func() (*internal.Executable, error)) internal.ProjectBuilder {
	return tool.NewCMake(executable)
}

func cmakeRunner(executable func() (*internal.Executable, error)) internal.ProjectRunner {
	return tool.NewCMake(executable)
}
