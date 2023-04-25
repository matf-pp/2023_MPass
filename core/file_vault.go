package core

import (
	"2023_MPass/encryption"
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/howeyc/gopass"
)

type FileVault struct {
	FilePath, VaultKey string
	entries            map[string]map[string]string
}

func (v *FileVault) Create() {
	fmt.Println("Enter the path where you want your vault to be stored:")
	reader := bufio.NewReader(os.Stdin)
	filepath, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("failed string read")
	}
	filepath = strings.TrimSuffix(filepath, "\n")

	v.FilePath = filepath
	// data := ""
	fmt.Println("Enter new master password: ")
	passwd, err := gopass.GetPasswd()
	if err != nil {
		log.Fatalf("password err")
	}
	v.VaultKey = string(passwd)
	data := `{ }`
	// data := ":"
	// _ = ioutil.WriteFile(v.FilePath, []byte(data), 0644)
	ciphertext := encryption.Encrypt(v.VaultKey, data)
	fmt.Println(ciphertext)
	encryption.StoreEncryptedData(v.FilePath, ciphertext)
	authKey, _ := encryption.CreateAuthKey(v.VaultKey, ciphertext)
	encryption.StoreAuthKey("key.txt", authKey)

}

func (v *FileVault) Load() {
	// file, err := os.Open(v.FilePath) //TODO open or create
	// if err != nil {
	//    panic(err)
	// }
	file := encryption.RetreiveEncryptedData(v.FilePath)
	// fmt.Println(v.VaultKey)
	//* should this stay here? -> no :)
	authKey := encryption.RetreiveAuthKey("key.txt") //TODO: fix hardcoded path for keys
	//* decrypt and truncate
	if encryption.ValidatePassword(v.VaultKey, file, authKey) {
		jsonBytes := encryption.Decrypt(v.VaultKey, file)
		if err := json.Unmarshal([]byte(jsonBytes), &(v.entries)); err != nil {
			panic(err)
		}
		// if err := os.Truncate(v.FilePath, 0); err != nil {
		// 	log.Printf("Failed to truncate db: %v", err)
		// }
		// if err := os.Truncate("key1.txt", 0); err != nil {
		// 	log.Printf("Failed to truncate key: %v", err)
		// }

	} else {
		fmt.Println("Wrong password!")
	}
	// jsonBytes, err := ioutil.ReadAll(file)
	// if err != nil {
	//    panic(err)
	// }
	// file, err := os.Open(v.FilePath) //TODO open or create
	// if err != nil {
	// 	panic(err)
	// }

	// //TODO decrypt

	// jsonBytes, err := ioutil.ReadAll(file)
	// if err != nil {
	// 	panic(err)
	// }

	// if err := json.Unmarshal(jsonBytes, &(v.entries)); err != nil {
	// 	panic(err)
	// }
}

func (v *FileVault) Store() {
	jsonBytes, err := json.Marshal(v.entries)
	if err != nil {
		panic(err)
	}

	//encrypt jsonBytes & save authentication key

	ciphertext := encryption.Encrypt(v.VaultKey, string(jsonBytes))
	encryption.StoreEncryptedData(v.FilePath, ciphertext)
	authKey, _ := encryption.CreateAuthKey(v.VaultKey, ciphertext)
	encryption.StoreAuthKey("key.txt", authKey)

	// err = os.WriteFile(v.FilePath, jsonBytes, 0644)
	// if err != nil {
	// 	panic(err)
	// }

}

func (v *FileVault) AddEntry(url string, username string, password string) {
	siteEntries, exists := v.entries[url]
	if !exists {
		v.entries[url] = make(map[string]string)
		siteEntries, _ = v.entries[url]
	}
	siteEntries[username] = password
}

func (v *FileVault) DeleteEntry(url string, username string) {
	siteEntries, _ := v.entries[url]
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

func (v *FileVault) GetEntries(url string) []VaultEntry {
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
	siteEntries, exists := v.entries[url]
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
	siteEntries, exists := v.entries[url]
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

func (v *FileVault) UpdateVaultKey(newPassphrase string) {
	v.VaultKey = newPassphrase
	v.Store()
}
