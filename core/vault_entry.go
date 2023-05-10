package core

type VaultEntry struct {
	url, username, password string
}

func CreateVaultEntry(url string, username string, password string) *VaultEntry {
	//TODO perform checks ?
	v := &VaultEntry{
		url:      url,
		username: username,
		password: password,
	}
	return v
}

func (v VaultEntry) GetUsername() string {
	return v.username
}

func (v VaultEntry) GetPassword() string {
	return v.password
}

func (v VaultEntry) GetUrl() string {
	return v.url
}
