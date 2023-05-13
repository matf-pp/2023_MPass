package main

import (
	"os/exec"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func check(e error, pages *tview.Pages, errView *tview.TextView) {
	if e != nil {
		errView.SetText("ERROR !")
		pages.SwitchToPage("errView")
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

	var textView = tview.NewTextView().SetTextColor(tcell.ColorWhiteSmoke).SetTextAlign(tview.AlignCenter)

	var pages = tview.NewPages().
	AddPage("textView", textView, true, true)

	var errView = tview.NewTextView().SetTextColor(tcell.ColorDarkRed).SetTextAlign(tview.AlignCenter)
	pages.AddPage("errView", errView, true, false)

	var masterPassForm = tview.NewForm().
		AddInputField("Vault name", "", 20, nil, nil).
		AddPasswordField("Old masterpass", "", 20, '*', nil).
		AddPasswordField("New masterpass", "", 20, '*', nil)
	masterPassForm.AddButton("Save", func() {
		vault := getInputFieldText(masterPassForm, 0)
		oldMasterPass := getInputFieldText(masterPassForm, 1)
		newMasterPass := getInputFieldText(masterPassForm, 2)
		_, err := exec.Command("../main/main", "change", "masterpass", "--vault", vault,
			"--masterpass", oldMasterPass, "--new-masterpass", newMasterPass).Output()
		vault = strings.TrimSpace(vault)
		oldMasterPass = strings.TrimSpace(oldMasterPass)
		newMasterPass = strings.TrimSpace(newMasterPass)
		if vault != "" && oldMasterPass != "" && newMasterPass != "" {
			check(err, pages, errView)
		}
	}).
		AddButton("Quit", func() {
			app.Stop()
		})

	pages.AddPage("masterPassForm", masterPassForm, true, false)

	var createForm = tview.NewForm().
		AddInputField("Vault name", "", 20, nil, nil).
		AddPasswordField("Masterpass", "", 20, '*', nil)
	createForm.AddButton("Create", func() {
		vault := getInputFieldText(createForm, 0)
		masterPass := getInputFieldText(createForm, 1)
		_, err := exec.Command("../main/main", "create", "--vault", vault,
			"--masterpass", masterPass).Output()
		vault = strings.TrimSpace(vault)
		masterPass = strings.TrimSpace(masterPass)
		if vault != "" && masterPass != "" {
			check(err, pages, errView)
		}
	}).
		AddButton("Quit", func() {
			app.Stop()
		})

	pages.AddPage("createForm", createForm, true, false)

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
		out, err := exec.Command("../main/main", "copy", "--vault", vault,
			"--masterpass", masterPass, "--url", url, "--username", username).Output()
		url = strings.TrimSpace(url)
		username = strings.TrimSpace(username)
		vault = strings.TrimSpace(vault)
		masterPass = strings.TrimSpace(masterPass)

		if url != "" && username != "" && vault != "" && masterPass != "" {
			check(err, pages, errView)
		}

		textView.SetText(string(out))
		pages.SwitchToPage("textView")

	}).
		AddButton("Delete", func() {
			url := getInputFieldText(cdForm, 0)
			username := getInputFieldText(cdForm, 1)
			vault := getInputFieldText(cdForm, 2)
			masterPass := getInputFieldText(cdForm, 3)
			_, err := exec.Command("../main/main", "delete", "--vault", vault,
				"--masterpass", masterPass, "--url", url, "--username", username).Output()
			url = strings.TrimSpace(url)
			username = strings.TrimSpace(username)
			vault = strings.TrimSpace(vault)
			masterPass = strings.TrimSpace(masterPass)

			if url != "" && username != "" && vault != "" && masterPass != "" {
				check(err, pages, errView)
			}
		}).
		AddButton("Quit", func() {
			app.Stop()
		})

	pages.AddPage("cdForm", cdForm, true, false)

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
		_, err := exec.Command("../main/main", "add", "--vault", vault, "--masterpass",
			masterPassword, "--url", url, "--username", username, "--password", password).Output()

		url = strings.TrimSpace(url)
		username = strings.TrimSpace(username)
		vault = strings.TrimSpace(vault)
		masterPassword = strings.TrimSpace(masterPassword)
		password = strings.TrimSpace(password)

		if url != "" && username != "" && vault != "" && password != "" && masterPassword != "" {
			check(err, pages, errView)
		}
	}).
		AddButton("Quit", func() {
			app.Stop()
		})
	
	pages.AddPage("addForm", addForm, true, false)

	var inputLenForm = tview.NewForm().
		AddInputField("Lenght", "", 3, nil, nil)
	inputLenForm.AddButton("Generate", func() {
		len := getInputFieldText(inputLenForm, 0)
		out, err := exec.Command("../main/main", "generate", "-l", len).Output()
		check(err, pages, errView)

		textView.SetText(string(out))
		pages.SwitchToPage("textView")
	})

	pages.AddPage("inputLenForm", inputLenForm, true, false)

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
		_, err := exec.Command("../main/main", "change", "password", "--vault", vault, "--masterpass",
			masterPassword, "--url", url, "--username", username,
			"--new-password", newPassword).Output()

		masterPassword = strings.TrimSpace(masterPassword)
		vault = strings.TrimSpace(vault)

		if err != nil {
			if vault != "" && masterPassword != "" {
				check(err, pages, errView)
			}
		}
	}).
		AddButton("Change username", func() {
			url := getInputFieldText(modifyForm, 0)
			username := getInputFieldText(modifyForm, 1)
			vault := getInputFieldText(modifyForm, 2)
			newUsername := getInputFieldText(modifyForm, 3)
			masterPassword := getInputFieldText(modifyForm, 5)
			_, err := exec.Command("../main/main", "change", "username", "--vault", vault, "--masterpass",
				masterPassword, "--url", url, "--username", username,
				"--new-username", newUsername).Output()
			username = strings.TrimSpace(username)
			masterPassword = strings.TrimSpace(masterPassword)
			vault = strings.TrimSpace(vault)
			url = strings.TrimSpace(url)

			if err != nil {
				if username != "" && url != "" && masterPassword != "" && vault != "" {
					check(err, pages, errView)
				}
			}
		}).
		AddButton("Quit", func() {
			app.Stop()
		})

	pages.AddPage("modifyForm", modifyForm, true, false)

	createForm.AddButton("List", func() {
		vaultName := getInputFieldText(createForm, 0)
		masterPass := getInputFieldText(createForm, 1)
		vaultList, err := exec.Command("../main/main", "list", "--vault", vaultName,
			"--masterpass", masterPass).Output()

		vaultName = strings.TrimSpace(vaultName)
		masterPass = strings.TrimSpace(masterPass)

		if vaultName != "" && masterPass != "" {
			check(err, pages, errView)
		}

		//display vault
		textView.SetText(string(vaultList))
		pages.SwitchToPage("textView")
	})

	pages.AddPage("createForm", createForm, true, false)

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
