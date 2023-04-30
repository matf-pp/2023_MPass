package databases

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type DatabaseInfo struct {
	info map[string]string
}

func (db *DatabaseInfo) OpenDb() {
	db.info = make(map[string]string)
	var file *os.File
	_, err := os.Stat(".databases")

	if os.IsNotExist(err) {
		file, err = os.Create(".databases")
		if err != nil {
			log.Fatalln("failed to create a file (probably)", err.Error())
		}
	} else {
		file, err = os.Open(".databases") //TODO: don't leave this hardcoded either. idc about it now
		if err != nil {
			log.Fatalln("failed to open .databases", err.Error())
		}
	}
	scanner := bufio.NewScanner(file)
	i := 0
	for scanner.Scan() {
		entry := strings.Split(scanner.Text(), ":")
		db.info[entry[0]] = entry[1]
		i++
	}
	// fmt.Println(db.info)
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	defer file.Close()
}

func (db *DatabaseInfo) FindKey(pathname string) string {
	// db.OpenDb()
	return db.info[pathname]
}

func (db *DatabaseInfo) UpdateAndStoreKeyHashes(pathname, keyHash string) {

	db.info[pathname] = keyHash

	stringline := ""
	for key, val := range db.info {
		tmpString := key + ":" + val + "\n"
		stringline += tmpString
	}
	// fmt.Println(stringline)
	err := ioutil.WriteFile(".databases", []byte(stringline), 0644)
	if err != nil {
		log.Fatalln("error writing to .databases...", err.Error())
	}
	// fmt.Println(pathname, hex.EncodeToString([]byte(keyHash)))
}
func (db *DatabaseInfo) DeleteDatabaseEntry(pathname string) {
	delete(db.info, pathname)
	stringline := ""
	for key, val := range db.info {
		tmpString := key + ":" + val + "\n"
		stringline += tmpString
	}
	// fmt.Println(stringline)
	err := ioutil.WriteFile(".databases", []byte(stringline), 0644)
	if err != nil {
		log.Fatalln("error writing to .databases...", err.Error())
	}
}
