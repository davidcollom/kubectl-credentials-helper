package cmd

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestRootCommand(t *testing.T) {
	// Test that the root command is properly initialized
	assert.NotNil(t, rootCmd)
	assert.Equal(t, "kubectl-credentials-keychain", rootCmd.Use)
	assert.Contains(t, rootCmd.Short, "Kubernetes credentials helper")

	logrus.Info("Testing root command")
}

func TestExecuteFunction(t *testing.T) {
	// Test that Execute function exists
	assert.NotNil(t, Execute)
}

func TestInitConfig(t *testing.T) {
	// Test that initConfig doesn't panic
	initConfig()

	// Verify that logger was initialized
	assert.NotNil(t, logger)
}
