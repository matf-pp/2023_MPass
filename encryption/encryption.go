// source : https://gist.github.com/tscholl2/dc7dc15dc132ea70a98e8542fefffa28
// TODO: fix error handling, hardcoded paths
// ...probably a lot more than that but w/e...
package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"
	"log"
	"strings"

	"golang.org/x/crypto/pbkdf2"
)

func deriveKey(passphrase string, salt []byte) ([]byte, []byte) {
	if salt == nil {
		salt = make([]byte, 8)
		rand.Read(salt)
		// salt = ""
	}
	return pbkdf2.Key([]byte(passphrase), salt, 1000, 32, sha256.New), salt
}

func encrypt(passphrase, plaintext string) string {
	key, salt := deriveKey(passphrase, nil)
	iv := make([]byte, 12)
	rand.Read(iv)
	b, _ := aes.NewCipher(key)
	aesgcm, _ := cipher.NewGCM(b)
	data := aesgcm.Seal(nil, iv, []byte(plaintext), nil)
	return hex.EncodeToString(salt) + "-" + hex.EncodeToString(iv) + "-" + hex.EncodeToString(data)
}

func decrypt(passphrase, ciphertext string) string {
	arr := strings.Split(ciphertext, "-")
	salt, _ := hex.DecodeString(arr[0])
	iv, _ := hex.DecodeString(arr[1])
	data, _ := hex.DecodeString(arr[2])
	key, _ := deriveKey(passphrase, salt)
	b, _ := aes.NewCipher(key)
	aesgcm, _ := cipher.NewGCM(b)
	data, _ = aesgcm.Open(nil, iv, data, nil)
	return string(data)
}

func createAuthKey(passphrase, cipthertext string) ([]byte, []byte) {
	arr := strings.Split(cipthertext, "-")
	salt, _ := hex.DecodeString(arr[0])
	key, _ := deriveKey(passphrase, salt)
	// fmt.Println(deriveKey(string(key), []byte(passphrase)))
	return deriveKey(string(key), []byte(passphrase))
}

func storeAuthKey(pathname string, authKey []byte) {
	authKeyHex := hex.EncodeToString(authKey)
	err := ioutil.WriteFile(pathname, []byte(authKeyHex), 0777)
	if err != nil {
		log.Fatal("key storage error")
	}
}
func retreiveAuthKey(pathname string) []byte {
	authKeyHex, err := ioutil.ReadFile(pathname)
	if err != nil {
		log.Fatalf("can't retreive the key..?")
	}
	authKey, _ := hex.DecodeString(string(authKeyHex))
	return authKey
}
func storeEncryptedData(pathname, encryptedData string) {
	err := ioutil.WriteFile(pathname, []byte(encryptedData), 0777)
	if err != nil {
		log.Fatalf("....")
	}
}
func retreiveEncryptedData(pathname string) string {
	data, err := ioutil.ReadFile(pathname)
	if err != nil {
		log.Fatalf("....")
	}
	return string(data)
}
func validatePassword(passphrase, cipthertext string, authKey []byte) bool {
	key, _ := createAuthKey(passphrase, cipthertext)
	keyHex := hex.EncodeToString(key)
	if keyHex == (hex.EncodeToString(authKey)) {
		return true
	}
	return false

}

// func main() {
// 	// passwords, err := ioutil.ReadFile("file.txt")
// 	// if err != nil {
// 	// 	log.Fatalf("file...: %v", err.Error())
// 	// }
// 	// key, err := ioutil.ReadFile("key.txt")
// 	// if err != nil {
// 	// 	log.Fatalf("file...: %v", err.Error())
// 	// }
// 	key := "thisISaKEy!!"
// 	// c := encrypt(string(key), string(passwords))
// 	// fmt.Println(decrypt(key, c))
// 	// storeEncryptedData("file.txt", c)
// 	// ioutil.WriteFile(hex.EncodeToString())
// 	c := retreiveEncryptedData("../file.txt")
// 	//* just a test; the only case the authentication key is going to be created is when we create a new database/change the master password
// 	newAuthenticationKey, _ := createAuthKey(string(key), c)
// 	storeAuthKey("../key.txt", newAuthenticationKey)
// 	keyFromFile := retreiveAuthKey("../key.txt")
// 	userInputPassword := "somewrongpassword"
// 	fmt.Println(validatePassword(userInputPassword, c, keyFromFile))
// 	userInputPassword = string(key)
// 	fmt.Println(validatePassword(userInputPassword, c, keyFromFile))
// 	plaintext := decrypt(key, c)
// 	fmt.Println(plaintext)

// }
