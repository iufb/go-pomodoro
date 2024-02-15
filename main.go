package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var wrapperStyle = lipgloss.NewStyle().Align(lipgloss.Center).Border(lipgloss.RoundedBorder()).Padding(4, 4)

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render

func main() {
	m := initialModel(workTime, shortBreak, longBreak)
	m.keymap.stop.SetEnabled(false)
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
