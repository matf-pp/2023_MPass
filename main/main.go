package main

import (
	"2023_MPass/core"
	"2023_MPass/encryption"
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/akamensky/argparse"
	"github.com/atotto/clipboard"
	"github.com/howeyc/gopass"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}
func loadDatabase() (string, string) {
	fmt.Println("Enter path to database: ")
	reader := bufio.NewReader(os.Stdin)
	filepath, err := reader.ReadString('\n')
	check(err)
	filepath_trim := strings.TrimSuffix(filepath, "\n")
	fmt.Println("Enter passcode: ")
	passwd, err := gopass.GetPasswd()
	check(err)

	return filepath_trim, string(passwd)

}
func loadVault() *core.FileVault {
	filepath, passwd := loadDatabase()
	v := &core.FileVault{
		// FilePath: "encoded.json",
		// VaultKey: "newPassphrase1",
		FilePath: filepath,
		VaultKey: passwd,
	}
	v.Load()
	return v
}

func main() {
	// var vv core.FileVault
	// vv.FilePath = "db.json"
	// vv.VaultKey = "key1"
	// vv.Create()
	//* note: in order for clipboard to work users need to install xclip or xsel

	parser := argparse.NewParser("MPass", "password manager program")
	//commands
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
	addPasswordOption := addCmd.String("p", "password", &argparse.Options{Required: true, Help: "password we want to add to new entry"})

	//change --masterpass
	changeCmd := parser.NewCommand("change", "changes password entry or masterpass")
	changeMasterPassCmd := changeCmd.NewCommand("masterpass", "changes masterpass")
	// newMasterPassOption := changeMasterPassCmd.String("n", "newPass", &argparse.Options{Required: true, Help: "takes in new masterpass"})

	//delete --database
	deleteDbCmd := parser.NewCommand("deletedb", "deletes an existing database")
	deleteDatabaseOption := deleteDbCmd.String("b", "database", &argparse.Options{Required: true, Help: "name of database we want to delete"})

	//change username --url --username --newusername
	changeUsernameCmd := changeCmd.NewCommand("username", "changes username of entry")
	changeUsernameUrlOption := changeUsernameCmd.String("u", "url", &argparse.Options{Required: true, Help: "url of entry we wish to update"})
	changeUsernameUsernameOption := changeUsernameCmd.String("n", "username", &argparse.Options{Required: true, Help: "username of entry we wish to update"})
	changeUsernameNewUsernameOption := changeUsernameCmd.String("w", "newusername", &argparse.Options{Required: true, Help: "new username"})

	//change password --url --username --newpassword
	changePasswordCmd := changeCmd.NewCommand("password", "changes password of entry")
	changePasswordUrlOption := changePasswordCmd.String("u", "url", &argparse.Options{Required: true, Help: "url of entry we wish to update"})
	changePasswordUsernameOption := changePasswordCmd.String("n", "username", &argparse.Options{Required: true, Help: "username of entry we wish to update"})
	// changePasswordNewPasswordOption := changePasswordCmd.String("p", "newpassword", &argparse.Options{Required: true, Help: "new password"})

	//delete --url --username
	deleteCmd := parser.NewCommand("delete", "deletes entry")
	deleteUrlOption := deleteCmd.String("u", "url", &argparse.Options{Required: true, Help: "url of entry we wish to delete"})
	deleteUsernameOption := deleteCmd.String("n", "username", &argparse.Options{Required: true, Help: "username of entry we wish to delete"})

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Println(parser.Usage(err))
		return
	}

	//TODO substitute prints with actual pass manager stuff duh
	if generateCmd.Happened() {
		randomString, err := encryption.GenerateRandomString(*lenOption)
		check(err)
		clipboard.WriteAll(randomString)
	} else {
		v := loadVault()
		if listCmd.Happened() {
			v.PrintVault()
		} else if copyCmd.Happened() {
			entry := v.GetEntry(*copyUrlOption, *copyUsernameOption)
			password := entry.GetPassword()
			clipboard.WriteAll(password)
			// time.Sleep(5)
			// exec.Command("xsel", "-z")

			// fmt.Println(string1)
		} else if changeMasterPassCmd.Happened() {
			fmt.Println("Enter new master password: ")
			passwd, err := gopass.GetPasswd()
			check(err)
			v.UpdateVaultKey(string(passwd)) //eg: newPassphrase1
		} else if changeUsernameCmd.Happened() {
			v.UpdateEntryUsername(*changeUsernameUrlOption, *changeUsernameUsernameOption, *changeUsernameNewUsernameOption)
			v.Store()
		} else if changePasswordCmd.Happened() {
			fmt.Println("Enter new entry password: ")
			passwd, err := gopass.GetPasswd()
			check(err)
			v.UpdateEntryPassword(*changePasswordUrlOption, *changePasswordUsernameOption, string(passwd))
			v.Store()
		} else if deleteCmd.Happened() {
			v.DeleteEntry(*deleteUrlOption, *deleteUsernameOption)
			v.Store()
		} else if addCmd.Happened() {
			v.AddEntry(*addUrlOption, *addUsernameOption, *addPasswordOption)
			v.Store()
		} else if deleteDbCmd.Happened() {
			v.Delete(*deleteDatabaseOption)
		}
	}
}
