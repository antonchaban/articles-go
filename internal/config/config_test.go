package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfigSuccessfully(t *testing.T) {
	cfg, err := Load()

	require.NoError(t, err)
	assert.NotNil(t, cfg)
}

func TestLoadConfigWithEnvironmentVariableOverride(t *testing.T) {
	_ = os.Setenv("HTTP_PORT", "9090")
	_ = os.Setenv("APP_ENV", "testing")
	defer func() {
		_ = os.Unsetenv("HTTP_PORT")
		_ = os.Unsetenv("APP_ENV")
	}()

	cfg, err := Load()

	require.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "9090", cfg.HTTPPort)
	assert.Equal(t, "testing", cfg.AppEnv)
}

func TestLoadConfigUsesDefaultsWhenNoConfigFile(t *testing.T) {
	_ = os.Setenv("APP_ENV", "production")
	_ = os.Setenv("HTTP_PORT", "3000")
	defer func() {
		_ = os.Unsetenv("APP_ENV")
		_ = os.Unsetenv("HTTP_PORT")
	}()

	cfg, err := Load()

	require.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "production", cfg.AppEnv)
	assert.Equal(t, "3000", cfg.HTTPPort)
}
