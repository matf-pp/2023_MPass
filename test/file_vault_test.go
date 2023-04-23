package test

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"2023_MPass/core"
)

func loadVault() *core.FileVault {
	v := &core.FileVault {
		FilePath : "vault.json",
		VaultKey: "ana",
	}
	v.Load()
	return v
}

func assertEntry(t *testing.T, e *core.VaultEntry, url string, username string, password string) {
	assert := assert.New(t)
	assert.Equal(url, e.GetUrl(), "invalid url")
	assert.Equal(username, e.GetUsername(), "invalid username")
	assert.Equal(password, e.GetPassword(), "invalid password")
}

func TestFileVaultGetEntries(t *testing.T) {
	// v := loadVault()
	// entries := v.GetEntries("google.com")
	// assert := assert.New(t)
	// assert.Equal(1, len(entries), "unexpected amount of entries")
	// testEntry(t, entries[0], "ana", "anaana123", "google.com")
}

func TestFileVaultGetEntry(t *testing.T) {
	v := loadVault()
	entry1 := v.GetEntry("google.com", "ana")
	entry2 := v.GetEntry("google.com", "ana@gmail.com")
	entry3 := v.GetEntry("facebook.com", "ruzica")
	entry4 := v.GetEntry("facebook.com", "ruzica1234")
	entry5 := v.GetEntry("amazon.com", "ruzica")
	assertEntry(t, entry1, "google.com", "ana", "anaana123")
	assertEntry(t, entry2, "google.com", "ana@gmail.com", "lospas")
	assertEntry(t, entry3, "facebook.com", "ruzica", "wewe")
	assert := assert.New(t)
	assert.Nil(entry4)
	assert.Nil(entry5)
}
