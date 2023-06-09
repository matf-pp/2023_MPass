package main

import (
	"bufio"
	vault "db/core"
	"db/encryption"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/akamensky/argparse"
	"github.com/atotto/clipboard"
	"github.com/howeyc/gopass"
)

// TODO (read one day maybe): better error handling. right now it does the job

func main() {

	// var v vault.FileVault
	// // fmt.Println(v)
	// v.FilePath = "test1"
	// v.VaultKey = "key1"
	// // v.Create()
	// v.Load()
	// v.AddEntry("google.com", "username1", "passwd1")
	// v.Load()
	// v.PrintVault()
	parser := argparse.NewParser("MPass", "Password manager program")
	masterPassOption := parser.String("m", "masterpass", &argparse.Options{Required: false, Help: "masterpass"})
	vaultNameOption := parser.String("v", "vault", &argparse.Options{Required: false, Help: "vault name"})
	//commands

	createCmd := parser.NewCommand("create", "creates a new database")

	generateCmd := parser.NewCommand("generate", "generates new password")
	lenOption := generateCmd.Int("l", "length", &argparse.Options{Required: true, Help: "length of an argument"})

	//list
	listCmd := parser.NewCommand("list", "lists vault")

	//copy --url --username
	copyCmd := parser.NewCommand("copy", "copies item i from vault to clipboard")
	copyUrlOption := copyCmd.String("u", "url", &argparse.Options{Required: true, Help: "url of entry we wish to copy"})
	copyUsernameOption := copyCmd.String("n", "username", &argparse.Options{Required: true, Help: "username of entry we wish to copy"})

	//add --url --username --password
	addCmd := parser.NewCommand("add", "adds new entry to vault")
	addUrlOption := addCmd.String("u", "url", &argparse.Options{Required: true, Help: "adds url to entry.."})
	addUsernameOption := addCmd.String("n", "username", &argparse.Options{Required: true, Help: "username we want to add for the new entry"})
	addPasswordOption := addCmd.String("p", "password", &argparse.Options{Required: false, Help: "password we want to add to new entry"})

	//change
	changeCmd := parser.NewCommand("change", "changes password entry or masterpass")
	changeMasterPassCmd := changeCmd.NewCommand("masterpass", "changes masterpass")
	newMasterPassOption := changeMasterPassCmd.String("n", "new-masterpass", &argparse.Options{Required: false, Help: "takes in new masterpass"})

	//delete --database
	deleteDbCmd := parser.NewCommand("deletedb", "deletes an existing database")
	deleteDatabaseOption := deleteDbCmd.String("b", "database", &argparse.Options{Required: true, Help: "name of database we want to delete"})

	//change username --url --username --new-username
	changeUsernameCmd := changeCmd.NewCommand("username", "changes username of entry")
	changeUsernameUrlOption := changeUsernameCmd.String("u", "url", &argparse.Options{Required: true, Help: "url of entry we wish to update"})
	changeUsernameUsernameOption := changeUsernameCmd.String("n", "username", &argparse.Options{Required: true, Help: "username of entry we wish to update"})
	changeUsernameNewUsernameOption := changeUsernameCmd.String("w", "new-username", &argparse.Options{Required: true, Help: "new username"})

	//change password --url --username --new-password
	changePasswordCmd := changeCmd.NewCommand("password", "changes password of entry")
	changePasswordUrlOption := changePasswordCmd.String("u", "url", &argparse.Options{Required: true, Help: "url of entry we wish to update"})
	changePasswordUsernameOption := changePasswordCmd.String("n", "username", &argparse.Options{Required: true, Help: "username of entry we wish to update"})
	changePasswordNewPasswordOption := changePasswordCmd.String("p", "new-password", &argparse.Options{Required: false, Help: "new password"})

	//delete --url --username
	deleteCmd := parser.NewCommand("delete", "deletes entry")
	deleteUrlOption := deleteCmd.String("u", "url", &argparse.Options{Required: true, Help: "url of entry we wish to delete"})
	deleteUsernameOption := deleteCmd.String("n", "username", &argparse.Options{Required: true, Help: "username of entry we wish to delete"})

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Println(parser.Usage(err))
		return
	}

	if generateCmd.Happened() {
		randomString, err := encryption.GenerateRandomString(*lenOption)
		// check(err)
		if err != nil {
			log.Fatalln(err.Error())
		}
		clipboard.WriteAll(randomString)
		fmt.Println("# Copied to clipboard!")
	} else {
		v := &vault.FileVault{
			FilePath: "",
			VaultKey: "",
		}
		if *vaultNameOption == "" {
			fmt.Println("Open database: ")
			reader := bufio.NewReader(os.Stdin)
			filepath, err := reader.ReadString('\n')
			if err != nil {
				log.Fatalln("Failed reading string from input")
			}
			v.FilePath = strings.TrimSuffix(filepath, "\n")
		} else {
			v.FilePath = *vaultNameOption
		}

		if *masterPassOption == "" {
			fmt.Println("Enter master password: ")
			passwd, err := gopass.GetPasswdMasked()
			if err != nil {
				log.Fatalln("master password error:", err.Error())
			}
			v.VaultKey = string(passwd)
		} else {
			v.VaultKey = *masterPassOption
		}

		if createCmd.Happened() {
			v.Create()
			// v.Store()
		} else {
			v.Load()
		}
		if listCmd.Happened() {
			v.PrintVault()
		} else if copyCmd.Happened() {
			entry := v.GetEntry(*copyUrlOption, *copyUsernameOption)
			if entry == nil {
				err := errors.New("Entry doesn't exist. Try again?")
				fmt.Println(err)
				os.Exit(0)
			}
			password := entry.GetPassword()
			clipboard.WriteAll(password)
			fmt.Println("# Copied to clipboard!")
		} else if changeMasterPassCmd.Happened() {
			passwd := *newMasterPassOption
			if passwd == "" {
				fmt.Println("Enter new master password: ")
				newPass, err := gopass.GetPasswdMasked()
				if err != nil {
					log.Fatalln("master password error:", err.Error())
				}
				passwd = string(newPass)
			}
			v.UpdateVaultKey(passwd)
		} else if changeUsernameCmd.Happened() { // change to non existing entry doesnt affect anything
			v.UpdateEntryUsername(*changeUsernameUrlOption, *changeUsernameUsernameOption, *changeUsernameNewUsernameOption)
			v.Store()
		} else if changePasswordCmd.Happened() {
			passwd := *changePasswordNewPasswordOption
			if passwd == "" {
				fmt.Println("Enter new password : ")
				newPass, err := gopass.GetPasswdMasked()
				if err != nil {
					log.Fatalln("password error:", err.Error())
				}
				passwd = string(newPass)
			}
			v.UpdateEntryPassword(*changePasswordUrlOption, *changePasswordUsernameOption, passwd)
			v.Store()
		} else if deleteCmd.Happened() { //deleting entry that doesnt exist doesnt affect the db
			v.DeleteEntry(*deleteUrlOption, *deleteUsernameOption)
			v.Store()
		} else if addCmd.Happened() {
			passwd := *addPasswordOption
			if passwd == "" {
				fmt.Println("Enter password : ")
				newPass, err := gopass.GetPasswdMasked()
				if err != nil {
					log.Fatalln("password error:", err.Error())
				}
				passwd = string(newPass)
			}
			v.AddEntry(*addUrlOption, *addUsernameOption, passwd)
			v.Store()
			fmt.Println("# Added to database!")
		} else if deleteDbCmd.Happened() {
			if v.FilePath == *deleteDatabaseOption {
				v.Delete(*deleteDatabaseOption)
			} else {
				fmt.Println("Database with that name doesn't exist. Try again?")
			}
		}
	}
}
