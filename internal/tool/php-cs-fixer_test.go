package tool

import (
	"cdt/internal"
	"cdt/internal/test"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPHPCSFixer_DetectPHPCSFixer_Composer(t *testing.T) {
	env := test.NewEnvironment(t)

	// Composer installation
	env.OnFindExecutable("vendor/bin/php-cs-fixer").
		Return(env.NewExecutable("/bin/tool"))

	tool := DetectPHPCSFixer(env)
	assert.NotNil(t, tool)
	assert.Equal(t, "php-cs-fixer", tool.Id())
	assert.True(t, tool.IsAvailable())

	if executable := tool.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/tool", executable.Path)
	}

	env.AssertExpectations(t)
}

func TestPHPCSFixer_DetectPHPCSFixer_System(t *testing.T) {
	env := test.NewEnvironment(t)

	// Composer installation
	env.OnFindExecutable("vendor/bin/php-cs-fixer").
		Return(nil)

	// System installation
	env.OnFindExecutable("php-cs-fixer").
		Return(env.NewExecutable("/bin/tool"))

	tool := DetectPHPCSFixer(env)
	assert.NotNil(t, tool)
	assert.Equal(t, "php-cs-fixer", tool.Id())
	assert.True(t, tool.IsAvailable())

	if executable := tool.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/tool", executable.Path)
	}

	env.AssertExpectations(t)
}

func TestPHPCSFixer_DetectPHPCSFixer_NotFound(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("vendor/bin/php-cs-fixer").
		Return(nil)
	env.OnFindExecutable("php-cs-fixer").
		Return(nil)

	tool := DetectPHPCSFixer(env)
	assert.NotNil(t, tool)
	assert.Equal(t, "php-cs-fixer", tool.Id())
	assert.False(t, tool.IsAvailable())
	assert.Nil(t, tool.Executable())

	env.AssertExpectations(t)
}

func TestPHPCSFixer_PHPCSFixer_FormatAll(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewPHPCSFixer(exec.LazyExecutable("format"))

	desc := internal.Project{Directory: "."}

	exec.OnRun("format", []string{"fix"}).
		Return(nil)

	err := tool.FormatAll(desc, []string{})
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestPHPCSFixer_PHPCSFixer_FormatFiles(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewPHPCSFixer(exec.LazyExecutable("format"))

	desc := internal.Project{Directory: "."}

	exec.OnRun("format", []string{"fix", "tests/*"}).
		Return(nil)

	err := tool.FormatFiles(desc, []string{"tests/*"}, []string{})
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestPHPCSFixer_PHPCSFixer_FormatCheckAll(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewPHPCSFixer(exec.LazyExecutable("format"))

	desc := internal.Project{Directory: "."}

	exec.OnRun("format", []string{"fix", "--dry-run"}).
		Return(nil)

	err := tool.FormatCheckAll(desc, []string{})
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestPHPCSFixer_PHPCSFixer_FormatCheckFiles(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewPHPCSFixer(exec.LazyExecutable("format"))

	desc := internal.Project{Directory: "."}

	exec.OnRun("format", []string{"fix", "--dry-run", "tests/*", "/path/to/file.php"}).
		Return(nil)

	err := tool.FormatCheckFiles(desc, []string{"tests/*", "/path/to/file.php"}, []string{})
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}
