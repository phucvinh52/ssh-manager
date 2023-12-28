package appui

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/gdamore/tcell/v2"
	"github.com/phucvinh52/ssh-manager/pkg/sshconfig"
	"github.com/rivo/tview"
)

type AppUI struct {
	App         *tview.Application
	Grid        *tview.Grid
	CTableSSH   *tview.Table
	CFilter     *tview.InputField
	sshCfg      *sshconfig.SSHConfig
	sshFiltered []*sshconfig.SSHHostFull
}

func (app *AppUI) initStyle() {
	//set grid
	app.Grid.
		SetRows(3, 1, 0).
		SetBorders(true)

	app.CTableSSH.
		SetBorders(false)

	app.CFilter.
		SetLabel("Filter").
		SetFieldWidth(0)

}

func (app *AppUI) initDraw() {
	app.Grid.
		AddItem(tview.NewTextView().SetText("ssh-manager\nControl+Q: quit\nControl+Space: Search"), 0, 0, 1, 3, 0, 0, false).
		AddItem(app.CFilter, 1, 0, 1, 3, 0, 0, true).
		AddItem(app.CTableSSH, 2, 0, 1, 3, 0, 0, false)
}

func (app *AppUI) Start() {
	if err := app.App.SetRoot(app.Grid, true).Run(); err != nil {
		panic(err)
	}
}

func (app *AppUI) keyBinding() {
	app.CFilter.SetDoneFunc(func(key tcell.Key) {
		row, _ := app.CTableSSH.GetSelection()
		if row <= 0 {
			return
		}
		app.execSSH(row - 1)
	})
	app.CFilter.SetChangedFunc(func(text string) {
		app.sshFiltered = app.sshCfg.Filter(text)
		app.CTableSSH.Clear()
		app.CTableSSH.SetCell(0, 0, tview.NewTableCell("#").SetTextColor(tcell.ColorGray).SetSelectable(false))
		app.CTableSSH.SetCell(0, 1, tview.NewTableCell("NAME").SetTextColor(tcell.ColorGray).SetSelectable(false))
		app.CTableSSH.SetCell(0, 2, tview.NewTableCell("HOSTNAME").SetTextColor(tcell.ColorGray).SetSelectable(false))
		app.CTableSSH.SetCell(0, 3, tview.NewTableCell("USER").SetTextColor(tcell.ColorGray).SetSelectable(false))
		app.CTableSSH.SetCell(0, 4, tview.NewTableCell("PORT").SetTextColor(tcell.ColorGray).SetSelectable(false))
		for k, v := range app.sshFiltered {
			idx := k + 1
			app.CTableSSH.SetCell(idx, 0, tview.NewTableCell(fmt.Sprint(idx)))
			app.CTableSSH.SetCell(idx, 1, tview.NewTableCell(v.Host))
			app.CTableSSH.SetCell(idx, 2, tview.NewTableCell(v.HostName))
			app.CTableSSH.SetCell(idx, 3, tview.NewTableCell(v.User))
			app.CTableSSH.SetCell(idx, 4, tview.NewTableCell(v.Port))
		}
		// // table.SetEvaluateAllRows(true)
		app.CTableSSH.SetSelectable(true, false)
		app.CTableSSH.Select(1, 0)
		app.CTableSSH.SetSelectedFunc(func(row, column int) {
			app.execSSH(row - 1)
		})
	})
	app.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlQ:
			app.App.Stop()
			os.Exit(0)
		case tcell.KeyEscape, tcell.KeyCtrlSpace:
			app.App.SetFocus(app.CFilter)
		case tcell.KeyDown, tcell.KeyUp:
			app.App.SetFocus(app.CTableSSH)
		}
		return event
	})
	// app.CFilter.SetChangedFunc(func(text string) {
	// 	sshHosts := app.sshCfg.Filter(text)
	// 	app.CTableSSH.Clear()
	// 	// table = tview.NewTable().SetBorders(false)
	// 	app.CTableSSH.SetCell(0, 0, tview.NewTableCell("NAME").SetTextColor(tcell.ColorGray).SetSelectable(false))
	// 	app.CTableSSH.SetCell(0, 1, tview.NewTableCell("HOSTNAME").SetTextColor(tcell.ColorGray).SetSelectable(false))
	// 	app.CTableSSH.SetCell(0, 2, tview.NewTableCell("USER").SetTextColor(tcell.ColorGray).SetSelectable(false))
	// 	app.CTableSSH.SetCell(0, 3, tview.NewTableCell("PORT").SetTextColor(tcell.ColorGray).SetSelectable(false))
	// 	for k, v := range sshHosts {
	// 		idx := k + 1
	// 		app.CTableSSH.SetCell(idx, 0, tview.NewTableCell(v.Host))
	// 		app.CTableSSH.SetCell(idx, 1, tview.NewTableCell(v.HostName))
	// 		app.CTableSSH.SetCell(idx, 2, tview.NewTableCell(v.User))
	// 		app.CTableSSH.SetCell(idx, 3, tview.NewTableCell(v.Port))
	// 	}
	// 	// table.SetEvaluateAllRows(true)
	// 	app.CTableSSH.SetSelectable(true, false)
	// 	app.CTableSSH.Select(1, 0)
	// 	// table.SetFixed(1, 1)
	// 	app.CTableSSH.SetSelectedFunc(func(row, column int) {
	// 		log.Println("ahhi")
	// 	})
	//
	// 	app.App.Draw()
	// })
}

func CreateApp(sshCfg *sshconfig.SSHConfig) *AppUI {
	app := new(AppUI)

	app.App = tview.NewApplication()
	app.Grid = tview.NewGrid()
	app.CTableSSH = tview.NewTable()
	app.CFilter = tview.NewInputField()
	app.sshCfg = sshCfg

	app.initStyle()
	app.initDraw()

	app.keyBinding()
	return app
}

func (app *AppUI) execSSH(idx int) {
	app.App.Stop()
	if idx >= len(app.sshFiltered) {
		return
	}
	fmt.Printf("Waiting connection [%s]...\n", app.sshFiltered[idx].Host)
	cmd := exec.Command("ssh", app.sshFiltered[idx].Host)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Run()

}
