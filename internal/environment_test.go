package internal_test

import (
	"bytes"
	"testing"

	"cdt/internal"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnvironment_EnvironmentProviders_PrintTable_Empty(t *testing.T) {
	providers := internal.EnvironmentProviders{}

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

func (t *testEnvironmentProvider) Detect(_ string) *internal.Environment {
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

func (t *testEnvironmentProvider) CreateEnvironment(
	_ string,
	_ string,
) (internal.Environment, error) {
	return nil, nil // nolint: nilnil
}

func TestEnvironment_EnvironmentProviders_PrintTable(t *testing.T) {
	providers := internal.EnvironmentProviders{
		internal.SystemEnvironmentProvider,
		&testEnvironmentProvider{},
	}

	output := bytes.Buffer{}
	providers.PrintTable(&output)

	assert.NotEmpty(t, output.String())
}

func TestEnvironment_SystemEnvironmentProvider_Data(t *testing.T) {
	assert.Equal(t, "system", internal.SystemEnvironmentProvider.Id())
	assert.Equal(t, "System", internal.SystemEnvironmentProvider.Name())
	assert.Equal(t, "Native OS system environment", internal.SystemEnvironmentProvider.Info())
	assert.Equal(t, []string{"s"}, internal.SystemEnvironmentProvider.Aliases())
	assert.True(t, internal.SystemEnvironmentProvider.IsAvailable())

	env, err := internal.SystemEnvironmentProvider.CreateEnvironment(".", "test")
	require.NoError(t, err)
	assert.Equal(t, internal.SystemEnvironment, env)
}

func TestEnvironment_SystemEnvironmentProvider_Detect(t *testing.T) {
	assert.Nil(t, internal.SystemEnvironmentProvider.Detect("."))
}

func TestEnvironment_SystemEnvironment_Id(t *testing.T) {
	assert.Equal(t, "system", internal.SystemEnvironment.Id())
}

func TestEnvironment_SystemEnvironment_Start(t *testing.T) {
	require.NoError(t, internal.SystemEnvironment.Start(t.Context()))
}

func TestEnvironment_SystemEnvironment_IsRunning(t *testing.T) {
	assert.True(t, internal.SystemEnvironment.IsRunning(t.Context()))
}

func TestEnvironment_SystemEnvironment_Stop(t *testing.T) {
	require.NoError(t, internal.SystemEnvironment.Stop(t.Context()))
}

func TestEnvironment_SystemEnvironment_Cleanup(t *testing.T) {
	require.NoError(t, internal.SystemEnvironment.Cleanup(t.Context()))
}

func TestEnvironment_SystemEnvironment_FindExecutable_NotFound(t *testing.T) {
	executable, err := internal.SystemEnvironment.FindExecutable(t.Context(), "tool-not-found")

	assert.Nil(t, executable)
	assert.Error(t, err)
}

func TestEnvironment_SystemEnvironment_FindExecutable(t *testing.T) {
	executable, err := internal.SystemEnvironment.FindExecutable(t.Context(), "echo")

	require.NotNil(t, executable)
	require.NoError(t, err)
	assert.NotNil(t, executable.Runtime)
	assert.Contains(t, executable.Path, "echo")
}

func TestEnvironment_SystemEnvironment_RunExecutable(t *testing.T) {
	buffer := bytes.Buffer{}
	options := internal.RunOptions{
		Directory: ".",
		Output:    &buffer,
		Error:     nil,
	}

	err := internal.SystemEnvironment.RunExecutable(t.Context(), options, "echo", []string{"test"})
	require.NoError(t, err)
	assert.Equal(t, "test\n", buffer.String())
}
