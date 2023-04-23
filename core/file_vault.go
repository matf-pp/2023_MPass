package core

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type FileVault struct {
	FilePath, VaultKey string
	entries map[string]map[string]string
}

func (v *FileVault) Load() {
	file, err := os.Open(v.FilePath) //TODO open or create
	if err != nil {
	   panic(err)
	}

	//TODO decrypt

	jsonBytes, err := ioutil.ReadAll(file)
	if err != nil {
	   panic(err)
	}

	if err := json.Unmarshal(jsonBytes, &(v.entries)); err != nil {
		panic(err)
	}
}

func (v *FileVault) Store() {
}

func (v *FileVault) AddEntry() {
}

func (v *FileVault) DeleteEntry() {
}

func (v *FileVault) GetEntry(url string, username string) *VaultEntry {
	siteEntries, exists := v.entries[url]
	if !exists {
		return nil
	}
	password, exists := siteEntries[username]
	if !exists {
		return nil
	}
	return CreateVaultEntry(url, username, password)
}

func (v *FileVault) GetEntries(url string) []VaultEntry{
    return nil
}

func (v *FileVault) UpdateEntry() {
}