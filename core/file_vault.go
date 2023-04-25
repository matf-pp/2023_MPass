package core

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"fmt"
	"2023_MPass/encryption"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

type FileVault struct {
	FilePath, VaultKey string
	entries map[string]map[string]string
}

func (v *FileVault) Load() {
	file, err := os.Open(v.FilePath) //TODO open or create
	check(err)

	encryptedBytes, err := ioutil.ReadAll(file)
	// fmt.Println(encryptedBytes)
	check(err)
	encryptedString := string(encryptedBytes)
	// fmt.Println(encryptedString)
	decryptedString := encryption.Decrypt(v.VaultKey, encryptedString)
	// fmt.Println(decryptedString)
	decryptedBytes := []byte(decryptedString)
	// fmt.Println(decryptedBytes)

	if err := json.Unmarshal(decryptedBytes, &(v.entries)); err != nil {
		panic(err)
	}
}

func (v *FileVault) Store() {
	jsonBytes, err := json.Marshal(v.entries)
	check(err)

	// fmt.Println(string(jsonBytes))
	jsonBytesString := string(jsonBytes)
	vaultCiphered := encryption.Encrypt(v.VaultKey, jsonBytesString)
	// fmt.Println(vaultCiphered
	
	f, err := os.Create(v.FilePath)
	check(err)
	_, err = f.WriteString(vaultCiphered)
	check(err)

}

func (v *FileVault) AddEntry(url string, username string, password string) {
	siteEntries, exists := v.entries[url]
	if !exists {
		v.entries[url] = make(map[string]string)
		siteEntries,_ = v.entries[url]
	}
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

func (v *FileVault) PrintVault() {
	fmt.Println("+++++++++++++++++++++++++++++++++++++++++++++++++")
    for url, usernameMap := range v.entries {
		for username, _ := range usernameMap {
			fmt.Println(url, " : ", username)
		}
    }
	fmt.Println("+++++++++++++++++++++++++++++++++++++++++++++++++")
}

func (v *FileVault) UpdateVaultKey(){

}