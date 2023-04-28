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
	file, err := os.Open(".databases") //TODO: don't leave this hardcoded either. idc about it now
	if err != nil {
		log.Fatal(err)
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
}

func (db *DatabaseInfo) FindKey(pathname string) string {
	// db.OpenDb()
	return db.info[pathname]
}

// *separate these two functions...
// * i literally have duplicate code with 1 line different
func (db *DatabaseInfo) UpdateAndStoreKeyHashes(pathname, keyHash string) {

	db.info[pathname] = keyHash

	stringline := ""
	for key, val := range db.info {
		tmpString := key + ":" + val + "\n"
		stringline += tmpString
	}
	// fmt.Println(stringline)
	err := ioutil.WriteFile(".databases", []byte(stringline), 0664)
	if err != nil {
		log.Fatalf("error writing to .databases...", err.Error())
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
	err := ioutil.WriteFile(".databases", []byte(stringline), 0664)
	if err != nil {
		log.Fatalf("error writing to .databases...", err.Error())
	}
}
