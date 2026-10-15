package test

import (
	"cdt/internal"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ToolDetector is a function that tries to detect a tool.
type ToolDetector func(context.Context, internal.DetectOptions) internal.Tool

// RunDetectTool run a test that verifies a function for detecting a tool.
//
// The detector function is used to detect a tool = tested function.
// The envMockSetup function is used to set up the environment mock.
// The resultCheck function is used to check the result of the detection.
// The toolPaths are user mapping of names to paths.
func RunDetectTool(
	t *testing.T,
	detector ToolDetector,
	envMockSetup func(env *Environment),
	resultCheck func(tool internal.Tool),
	toolPaths map[string]string,
) internal.Tool {
	env := NewEnvironment(t)

	envMockSetup(env)

	detectedTool := detector(t.Context(), internal.DetectOptions{
		Environment: env,
		ToolsPaths:  toolPaths,
	})
	resultCheck(detectedTool)

	env.AssertExpectations(t)

	return detectedTool
}

// RunDetectToolLastFound run a test that verifies function for detecting a tool by checking testPaths.
// It should succeed for the last path in testPaths.
//
// The detector function is used to detect a tool = tested function.
// The testPaths are the paths that the function should check. The last path is the one that should be found.
func RunDetectToolLastFound(
	t *testing.T,
	detector ToolDetector,
	testPaths []string,
) internal.Tool {
	return RunDetectTool(
		t,
		detector,
		func(env *Environment) {
			for idx, path := range testPaths {
				exe := env.OnFindExecutable(path)

				// Last path is the one that should be found.
				if idx == len(testPaths)-1 {
					exe.Return(env.NewExecutable("/bin/tool"), nil)
				} else {
					exe.Return(nil, nil)
				}
			}
		},
		func(tool internal.Tool) {
			require.NotNil(t, tool)
			assert.True(t, tool.IsAvailable())

			executable := tool.Executable()
			require.NotNil(t, executable)
			assert.Equal(t, "/bin/tool", executable.Path)
		},
		nil,
	)
}

// RunDetectToolNotFound run a test that verifies function for detecting a tool by checking testPaths.
// It should fail for all testPaths.
//
// The detector function is used to detect a tool = tested function.
// The testPaths are the paths that the function should check.
func RunDetectToolNotFound(
	t *testing.T,
	detector ToolDetector,
	testPaths []string,
) internal.Tool {
	return RunDetectTool(
		t,
		detector,
		func(env *Environment) {
			for _, path := range testPaths {
				env.OnFindExecutable(path).Return(nil, nil)
			}
		},
		func(tool internal.Tool) {
			require.NotNil(t, tool)
			assert.False(t, tool.IsAvailable())
			assert.Nil(t, tool.Executable())
		},
		nil,
	)
}

// RunDetectToolConfig run a test that verifies function for detecting a tool with custom tool mapping.
//
// The detector function is used to detect a tool = tested function.
// The searchName is the name of the tool in the tool mapping.
func RunDetectToolConfig(t *testing.T, detector ToolDetector, searchName string) internal.Tool {
	const toolPath = "tool-path"

	return RunDetectTool(
		t,
		detector,
		func(env *Environment) {
			env.OnFindExecutable(toolPath).
				Return(env.NewExecutable("/bin/tool"), nil)
		},
		func(tool internal.Tool) {
			assert.True(t, tool.IsAvailable())
			executable := tool.Executable()

			require.NotNil(t, executable)
			assert.Equal(t, "/bin/tool", executable.Path)
		},
		map[string]string{searchName: toolPath},
	)
}
