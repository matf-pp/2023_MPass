package main

import (
	"github.com/gdamore/tcell/v2"
    "github.com/rivo/tview"
	"os/exec"
	// "fmt"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
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
		AddItem("Quit", "Press to exit", 'q', func() {
			app.Stop()
		})

	var textView = tview.NewTextView()

	var masterPassForm = tview.NewForm().
		AddPasswordField("Masterpass", "", 20, '*', nil).
		AddButton("Save", func() {

		}).
		AddButton("Quit", func() {
			app.Stop()
		})

	var cdForm = tview.NewForm().
		AddInputField("Url", "", 20, nil, nil).
		AddInputField("Username", "", 20, nil, nil).
		AddPasswordField("Masterpass", "", 20, '*', nil).
		AddButton("Save", func() {

		}).
		AddButton("Quit", func() {
			app.Stop()
		})

	var addForm = tview.NewForm().
		AddInputField("Url", "", 20, nil, nil).
		AddInputField("Username", "", 20, nil, nil).
		AddPasswordField("New Password", "", 20, '*', nil).
		AddPasswordField("Masterpass", "", 20, '*', nil).
		AddButton("Save", func() {

		}).
		AddButton("Quit", func() {
			app.Stop()
		})

	var inputLenForm = tview.NewForm().
		AddInputField("Lenght", "", 3, nil, nil)
	inputLenForm.AddButton("Generate", func() {
		inputField := inputLenForm.GetFormItem(0).(*tview.InputField)
		text := inputField.GetText()
		_, err := exec.Command("../main/main","generate", "-l", text).Output()
		check(err)
	})

	var modifyForm = tview.NewForm().
		AddInputField("Url", "", 20, nil, nil).
		AddInputField("Username", "", 20, nil, nil).
		AddInputField("New username", "", 20, nil, nil).
		AddPasswordField("New Password", "", 20, '*', nil).
		AddPasswordField("Masterpass", "", 20, '*', nil).
		AddButton("Save", func() {

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
	AddPage("modifyForm", modifyForm, true, false)

	var flex = tview.NewFlex()
	flex.AddItem(list.SetSelectedBackgroundColor(tcell.ColorDarkRed), 0, 3, true).
	AddItem(pages, 0, 4, true)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'q' {
			app.Stop()
		} else if event.Rune() == 'g' {
			pages.SwitchToPage("inputLenForm")
		} else if event.Rune() == 'p' {
			pages.SwitchToPage("masterPassForm")
		} else if event.Rune() == 'c' || event.Rune() == 'd' {
			pages.SwitchToPage("cdForm")
		} else if event.Rune() == 'a' {
			pages.SwitchToPage("addForm")
		} else if event.Rune() == 'm' {
			pages.SwitchToPage("modifyForm")
		}
		// } else if event.Rune() == 'a' {
		// 	_, err := exec.Command("../main/main","add", "--url", form.GetFormItem(0).GetLabel(),
		// 	"--username", form.GetFormItem(1).GetLabel()).Output()
		// 	check(err)
		// } else if event.Rune() == 'g' {
		// 	inputField := form.GetFormItem(0).(*tview.InputField)
		// 	text := inputField.GetText()
		// 	_, err := exec.Command("../main/main","generate", "-l", text).Output()
		// 	check(err)
		// }
		return event
	 })

    if err := app.SetRoot(flex, true).EnableMouse(true).Run(); err != nil {
        panic(err)
    }
}