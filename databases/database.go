package dat

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DatabaseInfo struct {
	VaultName string `gorm:"primaryKey"`
	VaultHash string
}

func OpenDb() *gorm.DB {
	database, err := gorm.Open(sqlite.Open("../main/databases.db"), &gorm.Config{})
	if err != nil {
		log.Fatalln(err.Error())
	}
	if !database.Migrator().HasTable(&DatabaseInfo{}) {
		err = database.Migrator().CreateTable(&DatabaseInfo{})
		if err != nil {
			log.Fatalln(err.Error())
		}
	}
	return database
}

func FindKey(pathname string) string {
	// db.OpenDb()
	db := OpenDb()
	dbInfo := DatabaseInfo{}
	db.Table("database_infos").Where("vault_name = ?", pathname).Find(&dbInfo)

	return dbInfo.VaultHash
}

func UpdateAndStoreKeyHashes(pathname, keyHash string) {

	db := OpenDb()
	dbInfo := DatabaseInfo{}
	db.Table("database_infos").Where("vault_name == ?", pathname).Find(&dbInfo)
	if len(dbInfo.VaultName) == 0 {
		dbInfo.VaultName = pathname
	}
	dbInfo.VaultHash = keyHash
	db.Table("database_infos").Save(&dbInfo)
}
func DeleteDatabaseEntry(pathname string) {
	db := OpenDb()

	dbInfo := DatabaseInfo{}

	db.Table("database_infos").Where("vault_name == ?", pathname).Delete(&dbInfo) //...??? 🤔
}
