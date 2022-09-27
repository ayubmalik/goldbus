package main

import (
	"fmt"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/tbrandon/mbserver"
)

func main() {
	app := app.New()

	window := app.NewWindow("Goldbus v0.1")

	host := widget.NewEntry()
	host.SetText("localhost")

	port := widget.NewEntry()
	port.SetText("1502")
	statusText := binding.NewString()

	var (
		start *widget.Button
		stop  *widget.Button
		mb    *mbserver.Server
	)

	start = widget.NewButton("Start Server", func() {

		mb = mbserver.NewServer()
		err := mb.ListenTCP("localhost:1502") // TODO
		if err != nil {
			panic("could not start server!")
		}

		statusText.Set(fmt.Sprintf("Started server. Holding Register Count: %d", len(mb.HoldingRegisters)))
		start.Disable()
		stop.Enable()
	})

	stop = widget.NewButton("Stop Server", func() {
		mb.Close()
		statusText.Set("Server stopped.")
		stop.Disable()
		start.Enable()
	})
	stop.Disable()

	top := container.NewVBox(
		widget.NewCard("Modbus Server:", "", widget.NewForm(
			widget.NewFormItem("Host", host),
			widget.NewFormItem("Port", port),
			widget.NewFormItem("", start),
			widget.NewFormItem("", stop)),
		),
		layout.NewSpacer(),
		widget.NewLabelWithData(statusText),
	)

	registerUI := newRegisterUI()
	main := container.NewVBox(top, registerUI.container)

	window.SetContent(main)
	window.CenterOnScreen()
	window.Resize(fyne.NewSize(640, 480))

	go func() {
		time.Sleep(7 * time.Second)
		registerUI.update(100)
		fmt.Printf("objects: %v\n", registerUI.register)
	}()
	window.ShowAndRun()
}

type register struct {
	name    string
	_type   string
	address int
}

type registerUI struct {
	container *fyne.Container
	*register
	update func(value int)
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
