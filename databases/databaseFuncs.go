package databases

type Database interface {
	openDb() (map[string]string, int)
	FindKey(pathfile string) string
	UpdateAndStoreKeyHashes(pathname, keyHash string)
	DeleteDatabaseEntry(pathname string)
}
