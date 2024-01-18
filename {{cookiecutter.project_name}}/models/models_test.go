package models

import (
	"{{cookiecutter.project_name}}/config"
	"github.com/stretchr/testify/assert"
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

func TestRunInit(t *testing.T) {
	assert.Nil(t, Db, "Database should be Nil before RunInit()")
	RunInit()
	assert.NotNil(t, Db, "Database should be instantiated after RunInit()")
	ClearAll()
	CreateDummyPosts()
}

func TestSmoke(t *testing.T) {
	assert.Equal(t, 2, 2)
}
