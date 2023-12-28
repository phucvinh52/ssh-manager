package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/gdamore/tcell/v2"
	"github.com/phucvinh52/ssh-manager/pkg/sshconfig"
	"github.com/rivo/tview"
)

func main() {

	sshcfg, err := sshconfig.ParseSSHConfig(nil)
	if err != nil {
		panic(err)
	}

	sshHosts := sshcfg.Filter("")

	app := tview.NewApplication()
	grid := tview.NewGrid().
		SetRows(3).
		SetBorders(true).
		AddItem(tview.NewTextView().SetText("HELP"), 0, 0, 1, 3, 0, 0, false)

	// box := tview.NewBox().SetBorder(true).SetTitle("SSH Manager")
	// list := tview.NewList()
	table := tview.NewTable().SetBorders(false)

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
		app.Stop()
		fmt.Printf("Waiting connection [%s]...\n", sshHosts[row-1].Host)
		cmd := exec.Command("ssh", sshHosts[row-1].Host)
		cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
		cmd.Run()
	})

	grid.AddItem(table, 1, 0, 1, 3, 0, 0, true)

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
	if err := app.SetRoot(grid, true).Run(); err != nil {
		panic(err)
	}
}
