package test

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"2023_MPass/core"
)

func TestCreateVaultEntry(t *testing.T) {
	v := core.CreateVaultEntry("google.com", "ana", "anaana123")
	assert := assert.New(t)
	assert.Equal("ana", v.GetUsername(), "invalid username")
	assert.Equal("google.com", v.GetUrl(), "invalid url")
	assert.Equal("anaana123", v.GetPassword(), "invalid password")
}