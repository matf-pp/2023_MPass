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
	jsonBytes, err := json.Marshal(v.entries)
	if err != nil {
		panic(err)
	 }

	 //encrypt jsonBytes
 
	 err = os.WriteFile(v.FilePath, jsonBytes, 0644)
	 if err != nil {
		panic(err)
	 }

}

func (v *FileVault) AddEntry(url string, username string, password string) {
	siteEntries,_ := v.entries[url]
	siteEntries[username] = password
}

func (v *FileVault) DeleteEntry(url string, username string) {
	siteEntries,_ := v.entries[url]
	delete(siteEntries, username)
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
	siteEntries, exists := v.entries[url]
	if !exists {
		return nil
	}

	var entryArray []VaultEntry
	for username, password := range siteEntries {
        entryArray = append(entryArray, *(CreateVaultEntry(url, username, password)))
    }
	return entryArray
}

func (v *FileVault) UpdateEntryUsername(url string, oldUsername string, newUsername string) {
	siteEntries,exists := v.entries[url]
	if !exists {
		return 
	}
	password, exists := siteEntries[oldUsername]
	if !exists {
		return 
	}
	siteEntries[newUsername] = password
	delete(siteEntries, oldUsername)
}

func (v *FileVault) UpdateEntryPassword(url string, username string, newPassword string) {
	siteEntries,exists := v.entries[url]
	if !exists {
		return 
	}
	_, exists = siteEntries[username]
	if !exists {
		return 
	}
	siteEntries[username] = newPassword
}