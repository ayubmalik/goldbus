package main

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func main() {
	model := parseFlags(os.Args)
	fmt.Println("server = ", model.server)
	fmt.Println("regs = ", model.registers)

	p := tea.NewProgram(model)
	if err := p.Start(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var (
	textStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("12")).Render
	valueStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#04B575")).Render
	helpStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render
)

type register struct {
	address int
	rtype   string
	dtype   string
	value   float32
}

type server struct {
	host    string
	port    uint16
	slaveID uint16
}

type model struct {
	server    server
	registers []register
	interval  uint16
}

func (m *model) Init() tea.Cmd {
	return readReg
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case readRegResult:
		m.setResult(msg)
		return m, readReg

	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m *model) View() string {
	s := "\n"
	for _, reg := range m.registers {
		s += fmt.Sprintf("\n%-30s %v\n", textStyle(fmt.Sprintf("%d", reg.address)), valueStyle(fmt.Sprintf("%f", reg.value)))
	}
	s += helpStyle("\nq: exit\n")
	return s
}

func (m model) setResult(result readRegResult) {
	for i, reg := range m.registers {
		if reg.address == result.a {
			m.registers[i].value = result.v
			break
		}
	}
}

type readRegResult struct {
	a int
	v float32
}

func readReg() tea.Msg {
	time.Sleep(500 * time.Millisecond)
	return readRegResult{a: 12322, v: float32(time.Now().Second())}
}
