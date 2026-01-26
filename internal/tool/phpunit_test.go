package tool_test

import (
	"path/filepath"
	"testing"

	"cdt/internal"
	"cdt/internal/test"
	"cdt/internal/tool"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPHPUnit_DetectPHPUnit_Composer(t *testing.T) {
	env := test.NewEnvironment(t)

	// Composer installation
	env.OnFindExecutable(filepath.Join("vendor", "bin", "phpunit")).
		Return(env.NewExecutable("/bin/phpunit"), nil)

	phpUnit := tool.DetectPHPUnit(t.Context(), tool.DetectOptions{Environment: env})
	assert.NotNil(t, phpUnit)
	assert.Equal(t, "phpunit", phpUnit.Id())
	assert.True(t, phpUnit.IsAvailable())

	if executable := phpUnit.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/phpunit", executable.Path)
	}

	env.AssertExpectations(t)
}

func TestPHPUnit_DetectPHPUnit_System(t *testing.T) {
	env := test.NewEnvironment(t)

	// Composer installation
	env.OnFindExecutable(filepath.Join("vendor", "bin", "phpunit")).
		Return(nil, nil)

	// System installation
	env.OnFindExecutable("phpunit").
		Return(env.NewExecutable("/bin/phpunit"), nil)

	phpUnit := tool.DetectPHPUnit(t.Context(), tool.DetectOptions{Environment: env})
	assert.NotNil(t, phpUnit)
	assert.Equal(t, "phpunit", phpUnit.Id())
	assert.True(t, phpUnit.IsAvailable())

	if executable := phpUnit.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/phpunit", executable.Path)
	}

	env.AssertExpectations(t)
}

func TestPHPUnit_DetectPHPUnit_NotFound(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable(filepath.Join("vendor", "bin", "phpunit")).
		Return(nil, nil)
	env.OnFindExecutable("phpunit").
		Return(nil, nil)

	phpUnit := tool.DetectPHPUnit(t.Context(), tool.DetectOptions{Environment: env})
	assert.NotNil(t, phpUnit)
	assert.Equal(t, "phpunit", phpUnit.Id())
	assert.False(t, phpUnit.IsAvailable())
	assert.Nil(t, phpUnit.Executable())

	env.AssertExpectations(t)
}

func TestPHPUnit_DetectPHPUnit_Config(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("phpunit-12").
		Return(env.NewExecutable("/bin/phpunit"), nil)

	phpUnit := tool.DetectPHPUnit(t.Context(), tool.DetectOptions{
		Environment: env,
		ToolsPaths:  map[string]string{"phpunit": "phpunit-12"},
	})
	assert.NotNil(t, phpUnit)
	assert.Equal(t, "phpunit", phpUnit.Id())
	assert.True(t, phpUnit.IsAvailable())

	if executable := phpUnit.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/phpunit", executable.Path)
	}

	env.AssertExpectations(t)
}

func TestPHPUnit_PHPUnit_TestAll(t *testing.T) {
	exec := test.NewExecutable(t)

	phpUnit := tool.NewPHPUnit(exec.LazyExecutable("test"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("test", []string{}).
		Return(nil)

	err := phpUnit.TestAll(t.Context(), internal.ProjectTesterOptions{ProjectInfo: info})
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestPHPUnit_PHPUnit_Test(t *testing.T) {
	exec := test.NewExecutable(t)

	phpUnit := tool.NewPHPUnit(exec.LazyExecutable("test"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("test", []string{"tests/*"}).
		Return(nil)

	err := phpUnit.TestPattern(
		t.Context(),
		internal.ProjectTesterOptions{ProjectInfo: info},
		"tests/*",
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestPHPUnit_TestPattern_UnknownFormat(t *testing.T) {
	exec := test.NewExecutable(t)

	phpUnit := tool.NewPHPUnit(exec.LazyExecutable("test"))

	info := internal.ProjectInfo{Directory: "."}

	err := phpUnit.TestPattern(
		t.Context(),
		internal.ProjectTesterOptions{
			ProjectInfo: info,
			Output:      internal.OutputOptions[internal.TestsReportFormat]{Format: "test"},
		},
		"test1",
	)
	require.EqualError(t, err, "unknown report format: test")

	exec.AssertExpectations(t)
}

func TestPHPUnit_PHPUnit_CollectCoverageAll(t *testing.T) {
	exec := test.NewExecutable(t)

	phpUnit := tool.NewPHPUnit(exec.LazyExecutable("test"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("test", []string{"--coverage-text"}).
		Return(nil)

	err := phpUnit.CollectCoverageAll(
		t.Context(),
		internal.ProjectCoverageCollectorOptions{ProjectInfo: info},
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestPHPUnit_PHPUnit_CollectCoveragePattern(t *testing.T) {
	exec := test.NewExecutable(t)

	phpUnit := tool.NewPHPUnit(exec.LazyExecutable("test"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("test", []string{"--coverage-text", "tests/*"}).
		Return(nil)

	err := phpUnit.CollectCoveragePattern(
		t.Context(),
		internal.ProjectCoverageCollectorOptions{ProjectInfo: info},
		"tests/*",
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}
