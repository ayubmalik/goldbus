package main

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func main() {
	p := tea.NewProgram(initialModel())
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
	name    string
	_type   string
	address int
	value   float32
}

type server struct {
	host    string
	port    uint16
	slaveID uint8
}

type model struct {
	s    server
	regs []register
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
	for _, reg := range m.regs {
		s += fmt.Sprintf("\n%-30s %v\n", textStyle(reg.name), valueStyle(fmt.Sprintf("%f", reg.value)))
	}
	s += helpStyle("\nq: exit\n")
	return s
}

func (m model) setResult(result readRegResult) {
	for i, reg := range m.regs {
		if reg.address == result.a {
			m.regs[i].value = result.v
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
	return readRegResult{a: 12344, v: float32(time.Now().Second())}
}

func initialModel() *model {
	s := server{
		host:    "localhost",
		port:    1504,
		slaveID: 4,
	}
	r := register{
		name:    "auxiliary-active-power",
		_type:   "holding",
		address: 12344,
	}
	return &model{s: s, regs: []register{r}}
}
