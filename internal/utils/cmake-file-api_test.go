package utils_test

import (
	"cdt/internal/utils"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCmakeFileApi_Directories(t *testing.T) {
	dir := t.TempDir()

	api := utils.MakeCmakeFileApi(dir)

	assert.Equal(t, dir, api.BuildDirectory)
	assert.Equal(t, filepath.Join(dir, ".cmake/api/v1/"), api.GetApiDirectory())
	assert.Equal(t, filepath.Join(dir, ".cmake/api/v1/query"), api.GetQueryDirectory())
}

func TestCmakeFileApi_Query(t *testing.T) {
	dir := t.TempDir()

	api := utils.MakeCmakeFileApi(dir)

	err := api.Query("codemodel", 2)
	require.NoError(t, err)
	assert.FileExists(t, filepath.Join(api.GetQueryDirectory(), "codemodel-v2"))
}

// nolint: funlen
func TestCmakeFileApi_Reply_Executable(t *testing.T) {
	dir := t.TempDir()

	api := utils.MakeCmakeFileApi(dir)

	// Build reply
	require.NoError(t, os.MkdirAll(api.GetReplyDirectory(), 0755))

	indexPath := filepath.Join(api.GetReplyDirectory(), "index-123456.json")
	require.NoError(t, os.WriteFile(indexPath, []byte(`{
		"reply": {
			"codemodel-v2": {
				"jsonFile" : "codemodel-v2-123456.json",
				"kind" : "codemodel"
			}
		}
	}`), 0644))

	codemodelPath := filepath.Join(api.GetReplyDirectory(), "codemodel-v2-123456.json")
	require.NoError(t, os.WriteFile(codemodelPath, []byte(`{
		"configurations": [
			{
				"targets": [
					{
						"jsonFile": "target-test-123456.json",
						"name" : "test"
					}
				]
			}
		],
		"kind" : "codemodel"
	}`), 0644))

	targetPath := filepath.Join(api.GetReplyDirectory(), "target-test-123456.json")
	require.NoError(t, os.WriteFile(targetPath, []byte(`{
		"name" : "test",
		"nameOnDisk" : "test",
		"paths" : 
		{
			"build" : ".",
			"source" : "."
		},
		"sources" : 
		[
			{
				"path" : "main.cpp"
			},
			{
				"isGenerated" : true,
				"path" : "version.hpp"
			}
		],
		"type" : "EXECUTABLE"
	}
	`), 0644))

	reply, err := api.Reply()
	require.NoError(t, err)
	require.NotNil(t, reply)

	require.NotEmpty(t, reply.Targets)
	target1 := reply.Targets[0]
	assert.Equal(t, utils.TargetExecutable, target1.Type)
	assert.Equal(t, "test", target1.Name)
	assert.Equal(t, "test", target1.Path)
	assert.False(t, target1.External)
	assert.Equal(t, []string{"main.cpp"}, target1.Files)
}

// nolint: funlen
func TestCmakeFileApi_Reply_Utility(t *testing.T) {
	dir := t.TempDir()

	api := utils.MakeCmakeFileApi(dir)

	// Build reply
	require.NoError(t, os.MkdirAll(api.GetReplyDirectory(), 0755))

	indexPath := filepath.Join(api.GetReplyDirectory(), "index-123456.json")
	require.NoError(t, os.WriteFile(indexPath, []byte(`{
		"reply": {
			"codemodel-v2": {
				"jsonFile" : "codemodel-v2-123456.json",
				"kind" : "codemodel"
			}
		}
	}`), 0644))

	codemodelPath := filepath.Join(api.GetReplyDirectory(), "codemodel-v2-123456.json")
	require.NoError(t, os.WriteFile(codemodelPath, []byte(`{
		"configurations": [
			{
				"targets": [
					{
						"jsonFile": "target-test-123456.json",
						"name" : "test"
					}
				]
			}
		],
		"kind" : "codemodel"
	}`), 0644))

	targetPath := filepath.Join(api.GetReplyDirectory(), "target-test-123456.json")
	require.NoError(t, os.WriteFile(targetPath, []byte(`{
		"name" : "test",
		"nameOnDisk" : "test",
		"paths" : 
		{
			"build" : ".",
			"source" : "."
		},
		"sources" : 
		[
			{
				"path" : "main.cpp"
			},
			{
				"isGenerated" : true,
				"path" : "version.hpp"
			}
		],
		"type" : "UTILITY"
	}
	`), 0644))

	reply, err := api.Reply()
	require.NoError(t, err)
	require.NotNil(t, reply)

	assert.Empty(t, reply.Targets)
}
