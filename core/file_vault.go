package core

import (
	dat "db/databases"

	"db/encryption"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type FileVault struct {
	FilePath, VaultKey string
	entries            map[string]map[string]string
}
type VaultData struct {
	IdName         string `gorm:"primaryKey"`
	EncryptedEntry string
}

func DoesFileExist(pathname string) bool {
	_, err := os.Stat(pathname)
	if os.IsNotExist(err) {
		fmt.Println("File not found. Try again or create a new database....")
		return false
	}
	return true
}
func OpenVault() *gorm.DB {
	database, err := gorm.Open(sqlite.Open("../database/databases.db"), &gorm.Config{})

	if err != nil {
		log.Fatalln(err.Error())
	}
	return database
}
func (v *FileVault) Create() {
	db := OpenVault()
	data := "{ }"

	if !db.Migrator().HasTable(&VaultData{}) {
		err := db.Migrator().CreateTable(&VaultData{})
		if err != nil {
			log.Fatalln(err.Error())
		}
	}

	ciphertext := encryption.Encrypt(v.VaultKey, data)
	vault := VaultData{
		IdName:         v.FilePath,
		EncryptedEntry: ciphertext,
	}
	db.Where("id_name == ?", v.FilePath).Save(&vault)

	authKey, _ := encryption.CreateAuthKey(v.VaultKey, ciphertext)

	dat.UpdateAndStoreKeyHashes(v.FilePath, hex.EncodeToString(authKey))

}

func (v *FileVault) Load() {
	db := OpenVault()
	authKey := dat.FindKey(v.FilePath)

	encryptedInfo := VaultData{}
	db.Where("id_name == ?", v.FilePath).Find(&encryptedInfo)
	if len(encryptedInfo.IdName) == 0 {
		fmt.Errorf("No vault found...")
		os.Exit(0)
	}
	ciphertext := encryptedInfo.EncryptedEntry
	if encryption.ValidatePassword(v.VaultKey, ciphertext, authKey) {
		jsonBytes := encryption.Decrypt(v.VaultKey, ciphertext)

		if jsonBytes == "null" {
			jsonBytes = "{ }"
		}
		if err := json.Unmarshal([]byte(jsonBytes), &(v.entries)); err != nil {
			panic(err)
		}
	} else {
		fmt.Println("Wrong password!")
		os.Exit(0) //quits program if wrong password has been entered
	}
}

func (v *FileVault) Store() {
	db := OpenVault()
	jsonBytes, err := json.Marshal(v.entries)
	if err != nil {
		panic(err)
	}

	enc := VaultData{}
	ciphertext := encryption.Encrypt(v.VaultKey, string(jsonBytes))

	db.Where("id_name == ?", v.FilePath).Find(&enc)
	enc.EncryptedEntry = ciphertext
	db.Save(&enc)
	authKey, _ := encryption.CreateAuthKey(v.VaultKey, ciphertext)
	dat.UpdateAndStoreKeyHashes(v.FilePath, hex.EncodeToString(authKey))

}

func (v *FileVault) Delete(pathfile string) {
	db := OpenVault()
	vaultInfo := VaultData{}
	db.Where("id_name == ?", pathfile).Delete(&vaultInfo)
	dat.DeleteDatabaseEntry(pathfile)

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
	fmt.Println("\t\tVAULT: " + v.FilePath)
	fmt.Println("+++++++++++++++++++++++++++++++++++++++++++++++")
	for url, usernameMap := range v.entries {
		for username := range usernameMap {
			fmt.Printf("   %s : %-22s		\n", url, username)
		}
	}
	fmt.Println("+++++++++++++++++++++++++++++++++++++++++++++++")
}

func (v *FileVault) UpdateVaultKey(newPassphrase string) {
	v.VaultKey = newPassphrase
	v.Store()
}
