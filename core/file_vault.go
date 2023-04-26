package core

import (
	"2023_MPass/encryption"
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/howeyc/gopass"
)

// TODO: function(s) for deleting databases and their respective map entries. lmao
// * also refactoring and error handling

type FileVault struct {
	FilePath, VaultKey string
	entries            map[string]map[string]string
}

func openDb() (map[string]string, int) {
	database := make(map[string]string)
	file, err := os.Open(".databases") //TODO: don't leave this hardcoded either. idc about it now
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(file)
	i := 0
	for scanner.Scan() {
		entry := strings.Split(scanner.Text(), ":")
		database[entry[0]] = entry[1]
		i++
	}
	//fmt.Println(database)
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return database, i

}

func FindKey(pathfile string) (string, string) {
	db, _ := openDb()
	db1, _ := db[pathfile]
	return db1, db[pathfile]
}
func GenerateAndStoreKeyFile(pathfile string) string {
	db, i := openDb()
	db[pathfile] = "key" + fmt.Sprint(i) + ".txt"
	stringline := ""
	for key, val := range db {
		tmpString := key + ":" + val + "\n"
		stringline += tmpString
	}
	err := ioutil.WriteFile(".databases", []byte(stringline), 0664)
	if err != nil {
		log.Fatalf("error writing to .databases...")
	}
	return "key" + fmt.Sprint(i) + ".txt"

}

func (v *FileVault) Create() {
	fmt.Println("Name your new database:")
	reader := bufio.NewReader(os.Stdin)
	filepath, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Failed reading string from input")
	}
	filepath = strings.TrimSuffix(filepath, "\n")

	v.FilePath = filepath
	fmt.Println("Enter new master password: ")
	passwd, err := gopass.GetPasswd()
	if err != nil {
		log.Fatalf("master password error:", err.Error())
	}
	v.VaultKey = string(passwd)
	data := `{ }`
	ciphertext := encryption.Encrypt(v.VaultKey, data)
	encryption.StoreEncryptedData(v.FilePath, ciphertext)

	authKey, _ := encryption.CreateAuthKey(v.VaultKey, ciphertext)

	key := GenerateAndStoreKeyFile(v.FilePath)
	encryption.StoreAuthKey(key, authKey)

}

func (v *FileVault) Load() {
	key, _ := FindKey(v.FilePath)
	file := encryption.RetreiveEncryptedData(v.FilePath)
	if file == "-y" {
		v.Create()
		os.Exit(0)
	} else if file == "-n" {
		os.Exit(0)
	}
	// fmt.Println(v.VaultKey)
	//* should this stay here? -> no :)
	authKey := encryption.RetreiveAuthKey(key)
	//* decrypt and truncate
	//* no need for truncate, idk why i thought i needed that
	//* left as comment in case its needed for some reason
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
		os.Exit(0) //quits program if wrong password was entered
	}
}

func (v *FileVault) Store() {
	jsonBytes, err := json.Marshal(v.entries)
	if err != nil {
		panic(err)
	}

	//encrypt jsonBytes & save authentication key
	key, _ := FindKey(v.FilePath)
	ciphertext := encryption.Encrypt(v.VaultKey, string(jsonBytes))
	encryption.StoreEncryptedData(v.FilePath, ciphertext)
	authKey, _ := encryption.CreateAuthKey(v.VaultKey, ciphertext)
	encryption.StoreAuthKey(key, authKey)

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
	fmt.Println("\t\tDATABASE: " + v.FilePath)
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
