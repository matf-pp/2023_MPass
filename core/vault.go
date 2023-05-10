package core

import "gorm.io/gorm"

type Vault interface {
	DoesFileExist(pathname string) (bool, string)
	OpenVault() *gorm.DB
	Create()
	Load()
	Store()
	Delete(pathname string)
	AddEntry(url string, username string, password string)
	DeleteEntry(url string, username string)
	GetEntries(url string) []VaultEntry
	GetEntry(url string, username string) *VaultEntry
	UpdateEntryUsername(url string, oldUsername string, newUsername string)
	UpdateEntryPassword(url string, username string, newPassword string)
	UpdateVaultKey()
	PrintVault()
}
