package core
 
type Vault interface {
    Load()
    Store()
    AddEntry()
    DeleteEntry() 
    GetEntries(url string) []VaultEntry
    GetEntry(url string, username string) *VaultEntry
    UpdateEntry()
}