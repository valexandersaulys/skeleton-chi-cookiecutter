package models

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"{{cookiecutter.project_name}}/config"
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
	RunInit()
	db, err := GetDbWithNoContext()
	if err != nil {
		panic(err)
	}
	ClearAll(db)
	CreateDummyPosts(db)
}

func TestSmoke(t *testing.T) {
	assert.Equal(t, 2, 2)
}
