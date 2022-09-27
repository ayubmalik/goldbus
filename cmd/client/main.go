package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/goburrow/modbus"
)

const (
	title = "Goldbus client v0.1"
)

func main() {
	app := app.NewWithID("com.goldbus.Client")
	window := app.NewWindow(title)

	content := container.NewVBox(createServerUI(), createRegistersUI())
	window.SetContent(content)
	window.CenterOnScreen()

	window.Resize(fyne.NewSize(900, 660))
	window.Canvas().Scale()
	window.ShowAndRun()
}

func createServerUI() *fyne.Container {

	host := widget.NewEntry()
	host.SetText("localhost")

	port := widget.NewEntry()
	port.SetText("1501")

	serverCard := widget.NewCard("Modbus Server", "",
		widget.NewForm(
			widget.NewFormItem("Host", host),
			widget.NewFormItem("Port", port),
		),
	)

	status := widget.NewLabelWithStyle("Not connected.", fyne.TextAlignLeading, fyne.TextStyle{Bold: true, Monospace: true})

	connect := widget.NewButton("Connect", func() {

		handler := modbus.NewTCPClientHandler(fmt.Sprintf("%s:%s", host.Text, port.Text))
		handler.Timeout = 10 * time.Second
		handler.SlaveId = 0x01
		handler.Logger = log.New(os.Stdout, "test: ", log.LstdFlags)
		// Connect manually so that multiple requests are handled in one connection session
		handler.Connect()
		defer handler.Close()

		client := modbus.NewClient(handler)
		_, err := client.ReadHoldingRegisters(14505, 1)
		if err != nil {
			status.SetText(fmt.Sprintf("Connection error: %v", err))
		} else {
			status.SetText("Connected.")
		}

	})

	controlCard := container.NewVBox(connect)

	return container.NewBorder(nil, status, nil, controlCard, serverCard)
}

// registers contains the UI and logic
// for the registers in modbus
type registers struct {
	client *modbus.Client
}

func createRegistersUI() *fyne.Container {
	registers := container.NewGridWithColumns(4, widget.NewEntry(), widget.NewSelect([]string{"INPUT", "HOLDING"}, nil), widget.NewEntry())
	card := widget.NewCard("Registers", "", registers)
	return container.NewVBox(card)
}
