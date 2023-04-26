package core

type Vault interface {
	openDb() (map[string]string, int)
	FindKey(pathfile string) (string, string)
	GenerateAndStoreKeyFile(pathfile string) string
	Create()
	Load()
	Store()
	AddEntry(url string, username string, password string)
	DeleteEntry(url string, username string)
	GetEntries(url string) []VaultEntry
	GetEntry(url string, username string) *VaultEntry
	UpdateEntryUsername(url string, oldUsername string, newUsername string)
	UpdateEntryPassword(url string, username string, newPassword string)
	UpdateVaultKey()
	PrintVault()
}
