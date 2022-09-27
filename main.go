package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func main() {
	app := app.New()
	window := app.NewWindow("Hello World")

	window.SetContent(createGUI(window))
	window.CenterOnScreen()

	window.Resize(fyne.NewSize(600, 200))
	window.ShowAndRun()
}

func createGUI(w fyne.Window) *fyne.Container {
	start := widget.NewButton("Start Server", func() {
		dialog.ShowInformation("Start", "starting", w)
	})
	stop := widget.NewButton("Stop Server", func() {
		dialog.ShowConfirm("Start", "starting", nil, w)
	})
	ctrl := container.NewVBox(start, stop)

	server := widget.NewCard("Modbus Server", "", nil)
	top := container.NewBorder(nil, nil, nil, ctrl, server)
	return top

}
