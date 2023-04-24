// package main

// //TODO add help, quit

// import (
// 	"fmt"
// 	"os"

// 	"github.com/akamensky/argparse"
// )

// func main() {
// 	parser := argparse.NewParser("MPass", "password manager program")
// 	//commands
// 	generateCmd := parser.NewCommand("generate", "generates new password")
// 	lenOption := generateCmd.Int("l", "length", &argparse.Options{Required: true, Help: "length of an argument"})

// 	// //save --password s --url u --username n
// 	saveCmd := parser.NewCommand("save", "saves password")
// 	passOption := saveCmd.String("p", "password", &argparse.Options{Required: true, Help: "password"})
// 	urlOption := saveCmd.String("u", "url", &argparse.Options{Required: false, Help: "website url"})
// 	usernameOption := saveCmd.String("n", "username", &argparse.Options{Required: false, Help: "username"})

// 	//list
// 	//list -i x
// 	//TODO split this into two subcommands
// 	listCmd := parser.NewCommand("list", "lists vault")
// 	itemOption1 := listCmd.Int("i", "item", &argparse.Options{Required: false, Help: "returns element i in vault"})

// 	//copy -i x
// 	copyCmd := parser.NewCommand("copy", "copies item i from vault to clipboard")
// 	itemOption2 := copyCmd.Int("i", "item", &argparse.Options{Required: true, Help: "returns element i in vault"})

// 	//change --masterpass
// 	//change -i x -newpass s
// 	//TODO split this into two subcommands
// 	changeCmd := parser.NewCommand("change", "changes password entry or masterpass")
// 	masterpassOption := changeCmd.String("m", "masterpass", &argparse.Options{Required: false, Help: "changes masterpass"})
// 	itemOption3 := changeCmd.Int("i", "item", &argparse.Options{Required: false, Help: "returns element i in vault"})
// 	newPassOption := changeCmd.String("p", "newPass", &argparse.Options{Required: false, Help: "new password for item i in vault"})

// 	err := parser.Parse(os.Args)
// 	if err != nil {
// 		fmt.Println(parser.Usage(err))
// 		return
// 	}

// 	//TODO substitute prints with actual pass manager stuff duh
// 	if generateCmd.Happened() {
// 		fmt.Println(*lenOption)
// 	} else if saveCmd.Happened() {
// 		fmt.Println(*passOption)
// 		fmt.Println(*urlOption)
// 		fmt.Println(*usernameOption)
// 	} else if listCmd.Happened() {
// 		fmt.Println(*itemOption1)
// 	} else if copyCmd.Happened() {
// 		fmt.Println(*itemOption2)
// 	} else if changeCmd.Happened() {
// 		fmt.Println(*masterpassOption)
// 		fmt.Println(*itemOption3)
// 		fmt.Println(*newPassOption)
// 	}

// }
package main

import (
	"2023_MPass/core"
	"2023_MPass/encryption"
	"encoding/json"
	"fmt"

	"github.com/howeyc/gopass"
)

func loadVault() *core.FileVault {
	v := &core.FileVault{
		FilePath: "../test/vault.json",
		VaultKey: "ana",
	}
	v.Load()
	return v
}

func main() {
	var entries map[string]map[string]string

	// inputReader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter vault password:")
	key, _ := gopass.GetPasswd()

	// key = strings.TrimSuffix(key, "\n")
	result := encryption.RetreiveEncryptedData("encoded.json")
	// authKeystring, _ := encryption.CreateAuthKey(key, result)
	// encryption.StoreAuthKey("key.txt", authKey)

	// resString := encryption.Encrypt("newPassphrase", string(file))
	// fmt.Println(result)
	// _ = os.WriteFile("encoded.json", []byte(resString), 0644)
	authKey := encryption.RetreiveAuthKey("key.txt")
	if encryption.ValidatePassword(string(key), result, authKey) {
		decrypted := encryption.Decrypt(string(key), result)
		if err := json.Unmarshal([]byte(decrypted), &entries); err != nil {
			panic(err)
		}
		fmt.Println(entries)
	} else {
		fmt.Println("Wrong password!")
	}
}