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
	v := loadVault()
	entries := v.GetEntries("google.com")
	assertEntry(t, &(entries[0]), "google.com", "ana", "anaana123")
	assertEntry(t, &(entries[1]), "google.com", "ana@gmail.com", "lospas")
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

func TestFileEntryAddEntry(t *testing.T){
	v := loadVault()
	v.AddEntry("google.com", "ana111", "nekipassword")
	entry := v.GetEntry("google.com", "ana111")
	assertEntry(t, entry, "google.com", "ana111", "nekipassword")
}

func TestFileEntryDeleteEntry(t *testing.T){
	v := loadVault()
	v.DeleteEntry("google.com", "ana111")
	entry := v.GetEntry("google.com", "ana111")
	assert := assert.New(t)
	assert.Nil(entry)
}

func TestFileEntryUpdateEntryUsername(t *testing.T){
	v := loadVault()
	v.UpdateEntryUsername("google.com", "ana@gmail.com", "ana444")
	entry := v.GetEntry("google.com", "ana444")
	assertEntry(t, entry, "google.com", "ana444", "lospas")
}

func TestFileEntryUpdateEntryPassword(t *testing.T){
	v := loadVault()
	v.UpdateEntryPassword("google.com", "ana@gmail.com", "novasifra")
	entry := v.GetEntry("google.com", "ana@gmail.com")
	assertEntry(t, entry, "google.com", "ana@gmail.com", "novasifra")
}

//commented cause this changes the json test file
// func TestFileEntryStore(t *testing.T){
// 	v := loadVault()
// 	v.UpdateEntryPassword("google.com", "ana@gmail.com", "novasifra")
// 	v.UpdateEntryUsername("google.com", "ana@gmail.com", "ana444")
// 	v.Store()
// }