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
	"fmt"
	"os"

	"github.com/akamensky/argparse"
	"golang.design/x/clipboard"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func loadVault() *core.FileVault {
	v := &core.FileVault{
		FilePath: "encoded.json",
		VaultKey: "newPassphrase",
	}
	v.Load()
	return v
}

func main() {
	parser := argparse.NewParser("MPass", "password manager program")
	//commands
	generateCmd := parser.NewCommand("generate", "generates new password")
	lenOption := generateCmd.Int("l", "length", &argparse.Options{Required: true, Help: "length of an argument"})

	// //save --password s --url u --username n
	saveCmd := parser.NewCommand("save", "saves password")
	passOption := saveCmd.String("p", "password", &argparse.Options{Required: true, Help: "password"})
	urlOption := saveCmd.String("u", "url", &argparse.Options{Required: false, Help: "website url"})
	usernameOption := saveCmd.String("n", "username", &argparse.Options{Required: false, Help: "username"})

	//list
	//list -i x
	//TODO split this into two subcommands
	listCmd := parser.NewCommand("list", "lists vault")

	//copy -i x
	copyCmd := parser.NewCommand("copy", "copies item i from vault to clipboard")
	copyUrlOption := copyCmd.String("u", "url", &argparse.Options{Required: true, Help: "url of entry we wish to copy"})
	copyUsernameOption := copyCmd.String("n", "username", &argparse.Options{Required: true, Help: "username of entry we wish to copy"})

	changeCmd := parser.NewCommand("change", "changes password entry or masterpass")

	//change --masterpass
	changeMasterPassCmd := changeCmd.NewCommand("masterpass", "changes masterpass")
	newMasterPassOption := changeMasterPassCmd.String("n", "newPass", &argparse.Options{Required: true, Help: "takes in new masterpass"})

	//change username --url --username --newusername
	changeUsernameCmd := changeCmd.NewCommand("username", "changes username of entry")
	changeUsernameUrlOption := changeUsernameCmd.String("u", "url", &argparse.Options{Required: true, Help: "url of entry we wish to update"})
	changeUsernameUsernameOption := changeUsernameCmd.String("n", "username", &argparse.Options{Required: true, Help: "username of entry we wish to update"})
	changeUsernameNewUsernameOption := changeUsernameCmd.String("w", "newusername", &argparse.Options{Required: true, Help: "new username"})

	//change password --url --username --newpassword
	changePasswordCmd := changeCmd.NewCommand("password", "changes password of entry")
	changePasswordUrlOption := changePasswordCmd.String("u", "url", &argparse.Options{Required: true, Help: "url of entry we wish to update"})
	changePasswordUsernameOption := changePasswordCmd.String("n", "username", &argparse.Options{Required: true, Help: "username of entry we wish to update"})
	changePasswordNewPasswordOption := changePasswordCmd.String("p", "newpassword", &argparse.Options{Required: true, Help: "new password"})

	//delete --url --username
	deleteCmd := parser.NewCommand("delete", "deletes entry")
	deleteUrlOption := deleteCmd.String("u", "url", &argparse.Options{Required: true, Help: "url of entry we wish to delete"})
	deleteUsernameOption := deleteCmd.String("n", "username", &argparse.Options{Required: true, Help: "username of entry we wish to delete"})

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Println(parser.Usage(err))
		return
	}

	err = clipboard.Init()
	check(err)
	v := loadVault()
	//TODO substitute prints with actual pass manager stuff duh
	if generateCmd.Happened() {
		fmt.Println(*lenOption)
		password, err := encryption.GenerateRandomString(*lenOption)
		check(err)
		clipboard.Write(clipboard.FmtText, []byte(password))
		// fmt.Println(clipboard.Read(clipboard.FmtText))
	} else if saveCmd.Happened() {
		fmt.Println(*passOption)
		fmt.Println(*urlOption)
		fmt.Println(*usernameOption)
		v.AddEntry(*urlOption, *usernameOption, *passOption)
		v.Store()
	} else if listCmd.Happened() {
		v.PrintVault()
	} else if copyCmd.Happened() {
		fmt.Println(*copyUrlOption, *copyUsernameOption)
		//copies to clipboard
	} else if changeMasterPassCmd.Happened() {
		fmt.Println(*newMasterPassOption)
		v.UpdateVaultKey() //TODO
	} else if changeUsernameCmd.Happened() {
		fmt.Println(*changeUsernameUrlOption)
		fmt.Println(*changeUsernameUsernameOption)
		fmt.Println(*changeUsernameNewUsernameOption)
		v.UpdateEntryUsername(*changeUsernameUrlOption, *changeUsernameUsernameOption, *changeUsernameNewUsernameOption)
		v.Store()
	} else if changePasswordCmd.Happened() {
		fmt.Println(*changePasswordUrlOption)
		fmt.Println(*changePasswordUsernameOption)
		fmt.Println(*changePasswordNewPasswordOption)
		v.UpdateEntryPassword(*changePasswordUrlOption, *changePasswordUsernameOption, *changePasswordNewPasswordOption)
		v.Store()
	} else if deleteCmd.Happened() {
		fmt.Println(*deleteUrlOption)
		fmt.Println(*deleteUsernameOption)
		v.DeleteEntry(*deleteUrlOption, *deleteUsernameOption)
		v.Store()
	}
	// var v core.FileVault
	// v.FilePath = "encoded.json"
	// v.VaultKey = "newPassphrase"
	// v.Store()
}
