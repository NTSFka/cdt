package internal

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnvironment_EnvironmentProviders_PrintTable_Empty(t *testing.T) {
	providers := EnvironmentProviders{}

	output := bytes.Buffer{}
	providers.PrintTable(&output)

	assert.Empty(t, output.String())
}

type testEnvironmentProvider struct {
}

func (t *testEnvironmentProvider) Id() string {
	return "test"
}

func (t *testEnvironmentProvider) Aliases() []string {
	return []string{}
}

func (t *testEnvironmentProvider) ParameterInfo() string {
	return "test"
}

func (t *testEnvironmentProvider) Detect(_ string) *Environment {
	return nil
}

func (t *testEnvironmentProvider) Name() string {
	return "Test"
}

func (t *testEnvironmentProvider) Info() string {
	return "Test env provider"
}

func (t *testEnvironmentProvider) IsAvailable() bool {
	return false
}

func (t *testEnvironmentProvider) CreateEnvironment(_ string, _ string) (Environment, error) {
	return nil, nil // nolint: nilnil
}

func TestEnvironment_EnvironmentProviders_PrintTable(t *testing.T) {
	providers := EnvironmentProviders{
		SystemEnvironmentProvider,
		&testEnvironmentProvider{},
	}

	output := bytes.Buffer{}
	providers.PrintTable(&output)

	assert.NotEmpty(t, output.String())
}

func TestEnvironment_SystemEnvironmentProvider_Data(t *testing.T) {
	assert.Equal(t, "system", SystemEnvironmentProvider.Id())
	assert.Equal(t, "System", SystemEnvironmentProvider.Name())
	assert.Equal(t, "Native OS system environment", SystemEnvironmentProvider.Info())
	assert.Equal(t, []string{"s"}, SystemEnvironmentProvider.Aliases())
	assert.True(t, SystemEnvironmentProvider.IsAvailable())

	env, err := SystemEnvironmentProvider.CreateEnvironment(".", "test")
	require.NoError(t, err)
	assert.Equal(t, SystemEnvironment, env)
}

func TestEnvironment_SystemEnvironmentProvider_Detect(t *testing.T) {
	assert.Nil(t, SystemEnvironmentProvider.Detect("."))
}

func TestEnvironment_SystemEnvironment_Id(t *testing.T) {
	assert.Equal(t, "system", SystemEnvironment.Id())
}

func TestEnvironment_SystemEnvironment_Start(t *testing.T) {
	require.NoError(t, SystemEnvironment.Start(context.Background()))
}

func TestEnvironment_SystemEnvironment_IsRunning(t *testing.T) {
	assert.True(t, SystemEnvironment.IsRunning(context.Background()))
}

func TestEnvironment_SystemEnvironment_Stop(t *testing.T) {
	require.NoError(t, SystemEnvironment.Stop(context.Background()))
}

func TestEnvironment_SystemEnvironment_Cleanup(t *testing.T) {
	require.NoError(t, SystemEnvironment.Cleanup(context.Background()))
}

func TestEnvironment_SystemEnvironment_FindExecutable_NotFound(t *testing.T) {
	executable := SystemEnvironment.FindExecutable(context.Background(), "tool-not-found")

	assert.Nil(t, executable)
}

func TestEnvironment_SystemEnvironment_FindExecutable(t *testing.T) {
	executable := SystemEnvironment.FindExecutable(context.Background(), "echo")

	require.NotNil(t, executable)
	assert.NotNil(t, executable.Runtime)
	assert.Contains(t, executable.Path, "echo")
}

func TestEnvironment_SystemEnvironment_RunExecutable(t *testing.T) {
	buffer := bytes.Buffer{}
	options := RunOptions{
		Directory: ".",
		Output:    &buffer,
		Error:     nil,
	}

	err := SystemEnvironment.RunExecutable(context.Background(), options, "echo", []string{"test"})
	require.NoError(t, err)
	assert.Equal(t, "test\n", buffer.String())
}
