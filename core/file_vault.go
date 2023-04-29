package core

import (
	"2023_MPass/databases"
	"2023_MPass/encryption"
	"bufio"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/howeyc/gopass"
)

// TODO: function(s) for deleting databases and their respective map entries. lmao
// * also refactoring (slayed) and error handling (did not slay)

type FileVault struct {
	FilePath, VaultKey string
	entries            map[string]map[string]string
	db                 databases.DatabaseInfo
}

// * leave these functions here or move them to another file? package?
// func openDb() map[string]string {
// 	database := make(map[string]string)
// 	file, err := os.Open(".databases") //TODO: don't leave this hardcoded either. idc about it now
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	scanner := bufio.NewScanner(file)
// 	i := 0
// 	for scanner.Scan() {
// 		entry := strings.Split(scanner.Text(), ":")
// 		database[entry[0]] = entry[1]
// 		i++
// 	}
// 	// fmt.Println(database)
// 	if err := scanner.Err(); err != nil {
// 		log.Fatal(err)
// 	}
// 	return database

// }

//	func FindKey(pathfile string) string {
//		db := openDb()
//		return db[pathfile]
//	}
//
//	func UpdateAndStoreKeyHashes(pathname, keyHash string) {
//		db := openDb()
//		db[pathname] = keyHash
//		stringline := ""
//		for key, val := range db {
//			tmpString := key + ":" + val + "\n"
//			stringline += tmpString
//		}
//		// fmt.Println(stringline)
//		err := ioutil.WriteFile(".databases", []byte(stringline), 0664)
//		if err != nil {
//			log.Fatalf("error writing to .databases...", err.Error())
//		}
//		// fmt.Println(pathname, hex.EncodeToString([]byte(keyHash)))
//	}
func DoesFileExist(pathname string) (bool, string) {
	_, err := os.Stat(pathname)
	if os.IsNotExist(err) {
		fmt.Println("File not found. Try again or create a new database? -y -n")
		reader := bufio.NewReader(os.Stdin)
		readString, error := reader.ReadString('\n')
		if error != nil {
			log.Fatalf("Error while reading input -y -n")
		}
		readString = strings.TrimSuffix(readString, "\n")
		if readString == "-y" {
			// var v core.FileVault
			// v.Create()
			return false, "-y"
		} else if readString == "-n" {
			return false, "-n"
		} else {
			log.Fatalf("Invalid input..")
		}
	}
	return true, "..."
}

func (v *FileVault) Create() {
	// var db databases.DatabaseInfo
	v.db.OpenDb()
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

	v.db.UpdateAndStoreKeyHashes(v.FilePath, hex.EncodeToString(authKey))
	// StoreKeyHashes()
	// key := GenerateAndStoreKeyFile(v.FilePath)
	// encryption.StoreAuthKey(key, authKey)

}

func (v *FileVault) Load() {
	v.db.OpenDb()
	authKey := v.db.FindKey(v.FilePath)
	// fmt.Println(authKey)
	exists, file := DoesFileExist(v.FilePath)
	if exists == false {
		if file == "-y" {
			v.Create()
			os.Exit(0)
		} else if file == "-n" {
			os.Exit(0)
		}
	}
	ciphertext := encryption.RetreiveEncryptedData(v.FilePath)
	if encryption.ValidatePassword(v.VaultKey, ciphertext, authKey) {
		jsonBytes := encryption.Decrypt(v.VaultKey, ciphertext)
		if err := json.Unmarshal([]byte(jsonBytes), &(v.entries)); err != nil {
			panic(err)
		}
	} else {
		fmt.Println("Wrong password!")
		os.Exit(0) //quits program if wrong password has been entered
	}
}

func (v *FileVault) Store() {
	v.db.OpenDb()
	jsonBytes, err := json.Marshal(v.entries)
	if err != nil {
		panic(err)
	}

	//encrypt jsonBytes & save authentication key
	// key, _ := FindKey(v.FilePath)

	ciphertext := encryption.Encrypt(v.VaultKey, string(jsonBytes))
	encryption.StoreEncryptedData(v.FilePath, ciphertext)

	authKey, _ := encryption.CreateAuthKey(v.VaultKey, ciphertext)

	v.db.UpdateAndStoreKeyHashes(v.FilePath, hex.EncodeToString(authKey))
	// StoreKeyHashes()
	// encryption.StoreAuthKey(key, authKey)

}

func (v *FileVault) Delete(pathfile string) {
	v.db.OpenDb()
	err := os.Remove(v.FilePath)
	if err != nil {
		log.Fatalf("failed removing the file..", err.Error())
	}
	v.db.DeleteDatabaseEntry(pathfile)

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
