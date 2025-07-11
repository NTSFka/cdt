package internal

import (
	"bytes"
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEnvironment_EnvironmentProviders_PrintTable_Empty(t *testing.T) {
	providers := EnvironmentProviders{}

	output := bytes.Buffer{}
	providers.PrintTable(&output)

	assert.Empty(t, output.String())
}

func TestEnvironment_EnvironmentProviders_PrintTable(t *testing.T) {
	providers := EnvironmentProviders{
		SystemEnvironmentProvider,
	}

	output := bytes.Buffer{}
	providers.PrintTable(&output)

	assert.NotEmpty(t, output.String())
}

func TestEnvironment_SystemEnvironmentProvider_Data(t *testing.T) {
	assert.Equal(t, "system", SystemEnvironmentProvider.Id())
	assert.Equal(t, "System", SystemEnvironmentProvider.Name())
	assert.Equal(t, "Native OS system environment", SystemEnvironmentProvider.Info())
	assert.True(t, SystemEnvironmentProvider.IsAvailable())

	env, err := SystemEnvironmentProvider.CreateEnvironment(".", "test")
	assert.NoError(t, err)
	assert.Equal(t, SystemEnvironment, env)
}

func TestEnvironment_SystemEnvironment_Id(t *testing.T) {
	assert.Equal(t, "system", SystemEnvironment.Id())
}

func TestEnvironment_SystemEnvironment_Start(t *testing.T) {
	assert.NoError(t, SystemEnvironment.Start(context.Background()))
}

func TestEnvironment_SystemEnvironment_IsRunning(t *testing.T) {
	assert.True(t, SystemEnvironment.IsRunning(context.Background()))
}

func TestEnvironment_SystemEnvironment_Stop(t *testing.T) {
	assert.NoError(t, SystemEnvironment.Stop(context.Background()))
}

func TestEnvironment_SystemEnvironment_Cleanup(t *testing.T) {
	assert.NoError(t, SystemEnvironment.Cleanup(context.Background()))
}

func TestEnvironment_SystemEnvironment_FindExecutable_NotFound(t *testing.T) {
	executable := SystemEnvironment.FindExecutable("tool-not-found")

	assert.Nil(t, executable)
}

func TestEnvironment_SystemEnvironment_FindExecutable(t *testing.T) {
	executable := SystemEnvironment.FindExecutable("echo")

	if assert.NotNil(t, executable) {
		assert.NotNil(t, executable.RunFunc)
		assert.Contains(t, executable.Path, "echo")
	}
}

func TestEnvironment_SystemEnvironment_RunExecutable(t *testing.T) {
	buffer := bytes.Buffer{}
	options := RunOptions{
		Directory: ".",
		Output:    &buffer,
		Error:     nil,
	}

	err := SystemEnvironment.RunExecutable(context.Background(), options, "echo", []string{"test"})
	assert.NoError(t, err)
	assert.Equal(t, "test\n", buffer.String())
}
