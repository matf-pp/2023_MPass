package main

import (
	"github.com/gdamore/tcell/v2"
    "github.com/rivo/tview"
	"os/exec"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func getInputFieldText(inputForm *tview.Form, i int) string {
	inputField := inputForm.GetFormItem(i).(*tview.InputField)
	text := inputField.GetText()
	return text
}

var app = tview.NewApplication()

func main() {
	var list = tview.NewList().
		AddItem("List", "List all vault entries", 'l', nil).
		AddItem("Generate", "Generate new password", 'g', nil).
		AddItem("Copy", "Copy password of an entry", 'c', nil).
		AddItem("Add", "Add entry", 'a', nil).
		AddItem("Delete", "Delete entry", 'd', nil).
		AddItem("Modify", "Modify entry", 'm', nil).
		AddItem("Change", "Change masterpassword", 'p', nil).
		AddItem("Create", "Create new vault", 'v', nil).
		AddItem("Quit", "Press to exit", 'q', func() {
			app.Stop()
		})

	var textView = tview.NewTextView().SetTextColor(tcell.ColorRed)

	var masterPassForm = tview.NewForm().
		AddInputField("Vault name", "", 20, nil, nil).
		AddPasswordField("Old masterpass", "", 20, '*', nil).
		AddPasswordField("New masterpass", "", 20, '*', nil)
	masterPassForm.AddButton("Save", func() {
		vault := getInputFieldText(masterPassForm, 0)
		oldMasterPass := getInputFieldText(masterPassForm, 1)
		newMasterPass := getInputFieldText(masterPassForm, 2)
		_, err := exec.Command("../main/main","change", "masterpass", "--vault", vault, 
		"--masterpass", oldMasterPass, "--new-masterpass", newMasterPass).Output()
		check(err)
	}).
	AddButton("Quit", func() {
		app.Stop()
	})

	var createForm = tview.NewForm().
	AddInputField("Vault name", "", 20, nil, nil).
	AddPasswordField("Masterpass", "", 20, '*', nil)
	createForm.AddButton("Create", func() {
		vault := getInputFieldText(createForm, 0)
		masterPass := getInputFieldText(createForm, 1)
		_, err := exec.Command("../main/main","create", "--vault", vault, 
		"--masterpass", masterPass).Output()
		check(err)
	}).
	AddButton("Quit", func() {
		app.Stop()
	})

	var cdForm = tview.NewForm().
		AddInputField("Url", "", 20, nil, nil).
		AddInputField("Username", "", 20, nil, nil).
		AddInputField("Vault name", "", 20, nil, nil).
		AddPasswordField("Masterpass", "", 20, '*', nil)
	cdForm.AddButton("Copy", func() {
		url := getInputFieldText(cdForm, 0)
		username := getInputFieldText(cdForm, 1)
		vault := getInputFieldText(cdForm, 2)
		masterPass := getInputFieldText(cdForm, 3)
		_, err := exec.Command("../main/main","copy", "--vault", vault, 
		"--masterpass", masterPass, "--url", url, "--username", username).Output()
		check(err)
	}).
	AddButton("Delete", func() {
		url := getInputFieldText(cdForm, 0)
		username := getInputFieldText(cdForm, 1)
		vault := getInputFieldText(cdForm, 2)
		masterPass := getInputFieldText(cdForm, 3)
		_, err := exec.Command("../main/main","delete", "--vault", vault, 
		"--masterpass", masterPass, "--url", url, "--username", username).Output()
		check(err)
	}).
	AddButton("Quit", func() {
		app.Stop()
	})

	var addForm = tview.NewForm().
		AddInputField("Url", "", 20, nil, nil).
		AddInputField("Username", "", 20, nil, nil).
		AddInputField("Vault name", "", 20, nil, nil).
		AddPasswordField("New Password", "", 20, '*', nil).
		AddPasswordField("Masterpass", "", 20, '*', nil)
	addForm.AddButton("Save", func() {
		url := getInputFieldText(addForm, 0)
		username := getInputFieldText(addForm, 1)
		vault := getInputFieldText(addForm, 2)
		password := getInputFieldText(addForm, 3)
		masterPassword := getInputFieldText(addForm, 4)
		_, err := exec.Command("../main/main","add", "--vault", vault, "--masterpass",
		masterPassword, "--url", url, "--username", username, "--password", password).Output()
		check(err)
	}).
	AddButton("Quit", func() {
		app.Stop()
	})

	var inputLenForm = tview.NewForm().
		AddInputField("Lenght", "", 3, nil, nil)
	inputLenForm.AddButton("Generate", func() {
		len := getInputFieldText(inputLenForm, 0)
		_, err := exec.Command("../main/main","generate", "-l", len).Output()
		check(err)
	})

	var modifyForm = tview.NewForm().
		AddInputField("Url", "", 20, nil, nil).
		AddInputField("Username", "", 20, nil, nil).
		AddInputField("Vault name", "", 20, nil, nil).
		AddInputField("New username", "", 20, nil, nil).
		AddPasswordField("New Password", "", 20, '*', nil).
		AddPasswordField("Masterpass", "", 20, '*', nil)
	modifyForm.AddButton("Change password", func() {
		url := getInputFieldText(modifyForm, 0)
		username := getInputFieldText(modifyForm, 1)
		vault := getInputFieldText(modifyForm, 2)
		newPassword := getInputFieldText(modifyForm, 4)
		masterPassword := getInputFieldText(modifyForm, 5)
		_, err := exec.Command("../main/main","change", "password","--vault", vault, "--masterpass",
		masterPassword, "--url", url, "--username", username,
		"--new-password", newPassword).Output()
		check(err)
	}).
	AddButton("Change username", func() {
		url := getInputFieldText(modifyForm, 0)
		username := getInputFieldText(modifyForm, 1)
		vault := getInputFieldText(modifyForm, 2)
		newUsername := getInputFieldText(modifyForm, 3)
		masterPassword := getInputFieldText(modifyForm, 5)
		_, err := exec.Command("../main/main","change", "username","--vault", vault, "--masterpass",
		masterPassword, "--url", url, "--username", username,
		"--new-username", newUsername).Output()
		check(err)
	}).
	AddButton("Quit", func() {
		app.Stop()
	})
	
	var pages = tview.NewPages().
	AddPage("textView", textView, true, true).
	AddPage("masterPassForm", masterPassForm, true, false).
	AddPage("inputLenForm", inputLenForm, true, false).
	AddPage("cdForm", cdForm, true, false).
	AddPage("addForm", addForm, true, false).
	AddPage("modifyForm", modifyForm, true, false).
	AddPage("createForm", createForm, true, false)

	createForm.AddButton("List", func() {
		vaultName := getInputFieldText(createForm, 0)
		masterPass := getInputFieldText(createForm, 1)
		vaultList, err := exec.Command("../main/main","list", "--vault", vaultName, 
		"--masterpass", masterPass).Output()
		check(err)

		//display vault
		textView.SetText(string(vaultList))
		pages.SwitchToPage("textView")
	})

	var flex = tview.NewFlex()
	flex.AddItem(list.SetSelectedBackgroundColor(tcell.ColorDarkRed), 0, 3, true).
	AddItem(pages, 0, 4, true)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		//if ENTER is pressed
		if event.Rune() == 13 {
			if list.GetCurrentItem() == 0 {
				pages.SwitchToPage("createForm")
			} else if list.GetCurrentItem() == 1 {
				pages.SwitchToPage("inputLenForm")
			} else if list.GetCurrentItem() == 2 || list.GetCurrentItem() == 4 {
				pages.SwitchToPage("cdForm")
			} else if list.GetCurrentItem() == 3 {
				pages.SwitchToPage("addForm")
			} else if list.GetCurrentItem() == 5 {
				pages.SwitchToPage("modifyForm")
			} else if list.GetCurrentItem() == 6 {
				pages.SwitchToPage("masterPassForm")
			} else if list.GetCurrentItem() == 7 {
				pages.SwitchToPage("createForm")
			} else if list.GetCurrentItem() == 8 {
				app.Stop()
			}
		}
		return event
	 })

    if err := app.SetRoot(flex, true).EnableMouse(true).Run(); err != nil {
        panic(err)
    }
}