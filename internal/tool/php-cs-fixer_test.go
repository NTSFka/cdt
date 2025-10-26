package tool_test

import (
	"testing"

	"cdt/internal"
	"cdt/internal/test"
	"cdt/internal/tool"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPHPCSFixer_DetectPHPCSFixer_Composer(t *testing.T) {
	env := test.NewEnvironment(t)

	// Composer installation
	env.OnFindExecutable("vendor/bin/php-cs-fixer").
		Return(env.NewExecutable("/bin/php-cs-fixer"), nil)

	phpcsFixer := tool.DetectPHPCSFixer(t.Context(), internal.ConfigTools{}, env)
	assert.NotNil(t, phpcsFixer)
	assert.Equal(t, "php-cs-fixer", phpcsFixer.Id())
	assert.True(t, phpcsFixer.IsAvailable())

	if executable := phpcsFixer.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/php-cs-fixer", executable.Path)
	}

	env.AssertExpectations(t)
}

func TestPHPCSFixer_DetectPHPCSFixer_System(t *testing.T) {
	env := test.NewEnvironment(t)

	// Composer installation
	env.OnFindExecutable("vendor/bin/php-cs-fixer").
		Return(nil, nil)

	// System installation
	env.OnFindExecutable("php-cs-fixer").
		Return(env.NewExecutable("/bin/php-cs-fixer"), nil)

	phpcsFixer := tool.DetectPHPCSFixer(t.Context(), internal.ConfigTools{}, env)
	assert.NotNil(t, phpcsFixer)
	assert.Equal(t, "php-cs-fixer", phpcsFixer.Id())
	assert.True(t, phpcsFixer.IsAvailable())

	if executable := phpcsFixer.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/php-cs-fixer", executable.Path)
	}

	env.AssertExpectations(t)
}

func TestPHPCSFixer_DetectPHPCSFixer_NotFound(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("vendor/bin/php-cs-fixer").
		Return(nil, nil)
	env.OnFindExecutable("php-cs-fixer").
		Return(nil, nil)

	phpcsFixer := tool.DetectPHPCSFixer(t.Context(), internal.ConfigTools{}, env)
	assert.NotNil(t, phpcsFixer)
	assert.Equal(t, "php-cs-fixer", phpcsFixer.Id())
	assert.False(t, phpcsFixer.IsAvailable())
	assert.Nil(t, phpcsFixer.Executable())

	env.AssertExpectations(t)
}

func TestPHPCSFixer_DetectPHPCSFixer_Config(t *testing.T) {
	env := test.NewEnvironment(t)

	// Composer installation
	env.OnFindExecutable("php-cs-fixer-3").
		Return(env.NewExecutable("/bin/php-cs-fixer"), nil)

	phpcsFixer := tool.DetectPHPCSFixer(t.Context(), internal.ConfigTools{
		"php-cs-fixer": "php-cs-fixer-3",
	}, env)
	assert.NotNil(t, phpcsFixer)
	assert.Equal(t, "php-cs-fixer", phpcsFixer.Id())
	assert.True(t, phpcsFixer.IsAvailable())

	if executable := phpcsFixer.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/php-cs-fixer", executable.Path)
	}

	env.AssertExpectations(t)
}

func TestPHPCSFixer_PHPCSFixer_FormatAll(t *testing.T) {
	exec := test.NewExecutable(t)

	phpcsFixer := tool.NewPHPCSFixer(exec.LazyExecutable("format"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("format", []string{"fix"}).
		Return(nil)

	err := phpcsFixer.FormatAll(t.Context(), internal.ProjectFormatterOptions{ProjectInfo: info})
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestPHPCSFixer_PHPCSFixer_FormatFiles(t *testing.T) {
	exec := test.NewExecutable(t)

	phpcsFixer := tool.NewPHPCSFixer(exec.LazyExecutable("format"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("format", []string{"fix", "tests/*"}).
		Return(nil)

	err := phpcsFixer.FormatFiles(
		t.Context(),
		internal.ProjectFormatterOptions{ProjectInfo: info},
		[]string{"tests/*"},
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestPHPCSFixer_PHPCSFixer_FormatCheckAll(t *testing.T) {
	exec := test.NewExecutable(t)

	phpcsFixer := tool.NewPHPCSFixer(exec.LazyExecutable("format"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("format", []string{"fix", "--dry-run"}).
		Return(nil)

	err := phpcsFixer.FormatCheckAll(
		t.Context(),
		internal.ProjectFormatterOptions{ProjectInfo: info},
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestPHPCSFixer_PHPCSFixer_FormatCheckFiles(t *testing.T) {
	exec := test.NewExecutable(t)

	phpcsFixer := tool.NewPHPCSFixer(exec.LazyExecutable("format"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("format", []string{"fix", "--dry-run", "tests/*", "/path/to/file.php"}).
		Return(nil)

	err := phpcsFixer.FormatCheckFiles(
		t.Context(),
		internal.ProjectFormatterOptions{ProjectInfo: info},
		[]string{"tests/*", "/path/to/file.php"},
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}
