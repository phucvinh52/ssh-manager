package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/gdamore/tcell/v2"
	"github.com/phucvinh52/ssh-manager/pkg/sshconfig"
	"github.com/rivo/tview"
)

var app = tview.NewApplication()

var sshHosts []*sshconfig.SSHHostFull

func newAAA() {

}

func RunExec(idx int) {
	app.Stop()
	if idx >= len(sshHosts) {
		return
	}
	fmt.Printf("Waiting connection [%s]...\n", sshHosts[idx].Host)
	cmd := exec.Command("ssh", sshHosts[idx].Host)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Run()

}

func main() {

	sshcfg, err := sshconfig.ParseSSHConfig(nil)
	if err != nil {
		panic(err)
	}
	var tSearch = make(chan string, 1)

	grid := tview.NewGrid().
		SetRows(3, 1, 0).
		SetBorders(true).
		AddItem(tview.NewTextView().SetText("HELP"), 0, 0, 1, 3, 0, 0, false)

	// box := tview.NewBox().SetBorder(true).SetTitle("SSH Manager")
	// list := tview.NewList()
	table := tview.NewTable().SetBorders(false)

	go func() {
		for {
			select {
			case <-context.Background().Done():
				app.Stop()
				os.Exit(0)
			case txSearch := <-tSearch:
				sshHosts = sshcfg.Filter(txSearch)
				table.Clear()
				// table = tview.NewTable().SetBorders(false)
				table.SetCell(0, 0, tview.NewTableCell("NAME").SetTextColor(tcell.ColorGray).SetSelectable(false))
				table.SetCell(0, 1, tview.NewTableCell("HOSTNAME").SetTextColor(tcell.ColorGray).SetSelectable(false))
				table.SetCell(0, 2, tview.NewTableCell("USER").SetTextColor(tcell.ColorGray).SetSelectable(false))
				table.SetCell(0, 3, tview.NewTableCell("PORT").SetTextColor(tcell.ColorGray).SetSelectable(false))
				for k, v := range sshHosts {
					idx := k + 1
					table.SetCell(idx, 0, tview.NewTableCell(v.Host))
					table.SetCell(idx, 1, tview.NewTableCell(v.HostName))
					table.SetCell(idx, 2, tview.NewTableCell(v.User))
					table.SetCell(idx, 3, tview.NewTableCell(v.Port))
				}
				table.SetSelectable(true, false)
				// table.SetEvaluateAllRows(true)
				table.Select(1, 0)
				// table.SetFixed(1, 1)
				table.SetSelectedFunc(func(row, column int) {
					RunExec(row - 1)
				})

				app.Draw()
				// app.SetFocus(table)
			}

		}
	}()
	tSearch <- ""

	grid.AddItem(table, 2, 0, 1, 3, 0, 0, false)
	// table.Select(0, 0).SetFixed(1, 1).SetDoneFunc(func(key tcell.Key) {
	// 	if key == tcell.KeyEscape {
	// 		app.Stop()
	// 	}
	// 	if key == tcell.KeyEnter {
	// 		table.SetSelectable(true, false)
	// 	}
	// }).SetSelectedFunc(func(row int, column int) {
	// 	table.GetCell(row, column).SetTextColor(tcell.ColorRed)
	// 	table.SetSelectable(false, false)
	// })
	//
	// list := tview.NewList().AddItem("aaa", "aaa1", 'a', func() {
	// 	app.Stop()
	// 	cmd := exec.Command("ssh", "os6")
	// 	cmd.Stdout = os.Stdout
	// 	cmd.Stdin = os.Stdin
	// 	err := cmd.Run()
	// 	if err != nil {
	// 		log.Println(err)
	// 	}
	// }).AddItem("a2", "a22", 'b', func() {
	// 	log.Println("ahihi2")
	// })
	// if err := app.SetRoot(list, true).EnableMouse(true).Run(); err != nil {
	// 	panic(err)
	// }
	inputField := tview.NewInputField().
		SetLabel("Filter: ").
		SetFieldWidth(0)
	inputField.SetChangedFunc(func(text string) {
		tSearch <- inputField.GetText()
	})

	inputField.SetDoneFunc(func(key tcell.Key) {
		row, _ := table.GetSelection()
		if row <= 0 {
			return
		}
		RunExec(row - 1)
	})

	grid.AddItem(inputField, 1, 0, 1, 3, 0, 0, true)
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlQ:
			app.Stop()
			os.Exit(0)
		case tcell.KeyCtrlSpace:
			app.SetFocus(inputField)
		case tcell.KeyDown, tcell.KeyUp, tcell.KeyLeft, tcell.KeyRight:
			app.SetFocus(table)
		}
		return event
	})
	if err := app.SetRoot(grid, true).Run(); err != nil {
		panic(err)
	}
}
