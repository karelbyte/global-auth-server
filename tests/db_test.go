package test

import (
	"global-auth-server/libs"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDB_InvalidConfig(t *testing.T) {
	cfg := libs.DBConfig{
		Host:     "invalid-host",
		Port:     "1433",
		User:     "invalid",
		Password: "invalid",
		Name:     "nonexistent",
	}

	db, err := libs.NewDB(cfg)

	assert.Nil(t, db)
	assert.Error(t, err)
}
