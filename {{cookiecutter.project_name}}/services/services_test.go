package services

import (
	"example/skeleton/config"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Write code here to run before tests
	config.RunInit()

	// Run tests
	exitVal := m.Run()

	// Write code here to run after tests

	// Exit with exit value from tests
	os.Exit(exitVal)
}
