package encryption

// source : https://gist.github.com/tscholl2/dc7dc15dc132ea70a98e8542fefffa28
// TODO: fix error handling, hardcoded paths
// ...probably a lot more than that but w/e...

import (
	"bytes"
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

func DeriveKey(passphrase string, salt []byte) ([]byte, []byte) {
	if salt == nil {
		salt = make([]byte, 8)
		rand.Read(salt)
		// salt = ""
	}
	return pbkdf2.Key([]byte(passphrase), salt, 1000, 32, sha256.New), salt
}

func Encrypt(passphrase, plaintext string) string {
	key, salt := DeriveKey(passphrase, nil)
	iv := make([]byte, 12)
	rand.Read(iv)
	b, _ := aes.NewCipher(key)
	aesgcm, _ := cipher.NewGCM(b)
	data := aesgcm.Seal(nil, iv, []byte(plaintext), nil)
	return hex.EncodeToString(salt) + "-" + hex.EncodeToString(iv) + "-" + hex.EncodeToString(data)
}

func Decrypt(passphrase, ciphertext string) string {
	arr := strings.Split(ciphertext, "-")
	salt, _ := hex.DecodeString(arr[0])
	iv, _ := hex.DecodeString(arr[1])
	data, _ := hex.DecodeString(arr[2])
	key, _ := DeriveKey(passphrase, salt)
	b, _ := aes.NewCipher(key)
	aesgcm, _ := cipher.NewGCM(b)
	data, _ = aesgcm.Open(nil, iv, data, nil)

	return string(data)
}

func CreateAuthKey(passphrase, cipthertext string) ([]byte, []byte) {
	arr := strings.Split(cipthertext, "-")
	salt, _ := hex.DecodeString(arr[0])
	key, _ := DeriveKey(passphrase, salt)
	// fmt.Println(deriveKey(string(key), []byte(passphrase)))
	return DeriveKey(string(key), []byte(passphrase))
}

// * obsolete
func StoreAuthKey(pathname string, authKey []byte) {
	authKeyHex := hex.EncodeToString(authKey)
	err := ioutil.WriteFile(pathname, []byte(authKeyHex), 0644)
	if err != nil {
		log.Fatal("key storage error")
	}
}
func RetreiveAuthKey(pathname string) []byte {
	authKeyHex, err := ioutil.ReadFile(pathname)
	if err != nil {
		log.Fatalln("can't retreive the key..")
	}
	authKey, _ := hex.DecodeString(string(authKeyHex))
	return authKey
}
func StoreEncryptedData(pathname, encryptedData string) {
	err := ioutil.WriteFile(pathname, []byte(encryptedData), 0644)
	if err != nil {
		log.Fatalln("can't store data...")
	}
}
func RetreiveEncryptedData(pathname string) string {

	data, err := ioutil.ReadFile(pathname)

	if err != nil {
		log.Fatalln("Can't retreive data..")
	}
	return string(data)
}
func ValidatePassword(passphrase, cipthertext string, authKey string) bool {
	//* NIKAD ne proveravati hash sa []byte konverzijom nego hex.Decode string!!!!!
	key, _ := CreateAuthKey(passphrase, cipthertext)
	// fmt.Println(hex.EncodeToString(key))
	// fmt.Println(hex.EncodeToString(authKey))
	authKeyByte, err := hex.DecodeString(authKey)

	//* nil byte pri proveri nije smetao al mogu cisto reda radi da ga sklonim da se osiguram
	authKeyByte = bytes.Trim(authKeyByte, "\x00")
	key = bytes.Trim(key, "\x00")
	if err != nil {
		log.Fatalln(err.Error())
	}
	if bytes.Equal(key, authKeyByte) {
		return true
	}
	return false

}
