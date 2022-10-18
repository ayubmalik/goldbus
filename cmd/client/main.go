package main

import (
	"fmt"
	"log"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const (
	version = "v0.1"
	title   = "Goldbus Client"
)

func main() {

	application := app.New()

	serverUI := newServerUI()
	registersUI := newRegistersUI()
	registersUI.addRegister()
	controlUI := newControlUI()
	ui := container.NewVBox(serverUI.container, registersUI.container, layout.NewSpacer(), controlUI.container)

	window := application.NewWindow(fmt.Sprintf("%s %s", title, version))
	window.SetContent(ui)
	window.CenterOnScreen()
	window.Resize(fyne.NewSize(640, 480))
	window.ShowAndRun()
}

type serverUI struct {
	container *fyne.Container
}

func newServerUI() *serverUI {
	host := widget.NewEntry()
	host.SetText("localhost")
	port := widget.NewEntry()
	port.SetText("1501")
	content := container.New(layout.NewFormLayout(), widget.NewLabel("Host"), host, widget.NewLabel("Port"), port)
	card := widget.NewCard("Modbus Details", "", content)
	return &serverUI{container: container.NewMax(card)}
}

type registersUI struct {
	container *fyne.Container
}

func (ui *registersUI) addRegister() {
	reg := newRegisterUI()
	ui.container.Objects = append(ui.container.Objects, reg.container)
	log.Println("adding reg")
}

func newRegistersUI() *registersUI {
	var ui *registersUI
	add := widget.NewButtonWithIcon("Add Register", theme.ContentAddIcon(), func() { ui.addRegister() })
	content := container.NewVBox(container.NewHBox(add), container.NewVBox())
	ui = &registersUI{container: content}
	return ui
}

type registerUI struct {
	container *fyne.Container
	*register
	update func(value int)
}

// TODO: move pkg
type register struct {
	name    string
	_type   string
	address int
}

func newRegisterUI() *registerUI {
	n := widget.NewEntry()
	n.SetPlaceHolder("name")

	t := widget.NewSelect([]string{"INPUT", "HOLDING"}, nil)
	t.SetSelected("HOLDING")
	a := widget.NewEntry()
	a.SetPlaceHolder("address")

	value := binding.NewInt()
	v := widget.NewEntryWithData(binding.IntToString(value))
	v.Disable()

	r := &register{}
	n.OnChanged = func(s string) { r.name = s }
	t.OnChanged = func(s string) { r._type = s }
	a.OnChanged = func(s string) { r.address, _ = strconv.Atoi(s) }

	return &registerUI{
		register:  r,
		container: container.NewGridWithColumns(4, n, t, a, v),
		update: func(n int) {
			value.Set(n)
		},
	}
}

type controlUI struct {
	container *fyne.Container
}

func newControlUI() *controlUI {
	read := widget.NewButtonWithIcon("Read", theme.MediaPlayIcon(), func() {})
	stop := widget.NewButtonWithIcon("Stop", theme.MediaStopIcon(), func() {})
	container := container.NewGridWithColumns(3, layout.NewSpacer(), stop, read)
	return &controlUI{container}
}
