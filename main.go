package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type status int

const (
	start status = iota
	work
	short
	long
)

type keymap struct {
	start  key.Binding
	stop   key.Binding
	reset  key.Binding
	change key.Binding
	quit   key.Binding
}

type model struct {
	status       status
	timer        timer.Model
	workTime     int
	shortBreak   int
	longBreak    int
	keymap       keymap
	help         help.Model
	currentRound int
	quitting     bool
}

func (m model) GetStatus() string {
	var s string
	switch m.status {
	case 0:
		s = "START"

	case 1:
		s = "WORK"
	case 2:
		s = "SHORT BREAK"
	case 3:
		s = "LONG BREAK"
	}
	return s + "\n"
}

func (m model) getStyle(s string) string {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("#c2a8c2")).Padding(1, 2).Border(lipgloss.RoundedBorder()).Render(s)
}

func (m model) Notify() {
	notify := exec.Command("notify-send", m.GetStatus())
	sound := exec.Command("paplay", "./sound.wav")
	notify.Start()
	sound.Start()
}

func initialModel() model {
	return model{
		status:     start,
		timer:      timer.NewWithInterval(time.Minute*25, time.Second),
		workTime:   25,
		shortBreak: 5,
		longBreak:  15,
		keymap: keymap{
			start: key.NewBinding(
				key.WithKeys("s"),
				key.WithHelp("s", "start"),
			),
			stop: key.NewBinding(
				key.WithKeys("p"),
				key.WithHelp("p", "stop"),
			),
			reset: key.NewBinding(
				key.WithKeys("r"),
				key.WithHelp("r", "reset"),
			),
			change: key.NewBinding(
				key.WithKeys("c"),
				key.WithHelp("c", "change timer value"),
			),
			quit: key.NewBinding(
				key.WithKeys("q", "ctrl+c"),
				key.WithHelp("q", "quit"),
			),
		},
		help: help.New(),
	}
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return tea.SetWindowTitle("Pomodoro CLI App.")
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case timer.TickMsg:
		var cmd tea.Cmd
		m.timer, cmd = m.timer.Update(msg)
		return m, cmd

	case timer.StartStopMsg:
		var cmd tea.Cmd
		m.timer, cmd = m.timer.Update(msg)
		m.keymap.stop.SetEnabled(m.timer.Running())
		m.keymap.start.SetEnabled(!m.timer.Running())
		return m, cmd

	case timer.TimeoutMsg:
		// m.quitting = true
		switch m.status {
		case work:
			if m.currentRound >= 3 {
				m.status = long
				m.currentRound = 0
				m.timer.Timeout = time.Minute * time.Duration(m.longBreak)
			} else {
				m.status = short
				m.timer.Timeout = time.Minute * time.Duration(m.shortBreak)
			}
		default:
			m.status = work
			m.currentRound += 1
			m.timer.Timeout = time.Minute * time.Duration(m.workTime)
		}
		m.Notify()
		return m, cmd
	case inputs:
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.quit):
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, m.keymap.reset):
			m.timer.Timeout = time.Minute * time.Duration(m.workTime)
		case key.Matches(msg, m.keymap.change):
			return inputsInitialModel(), nil
		case key.Matches(msg, m.keymap.start):
			m.timer.Init()
			m.status = work
			m.Notify()
			return m, m.timer.Start()
		case key.Matches(msg, m.keymap.stop):
			return m, m.timer.Stop()
		}

	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render

func (m model) helpView() string {
	return helpStyle("\n" + m.help.ShortHelpView([]key.Binding{
		m.keymap.start,
		m.keymap.stop,
		m.keymap.reset,
		m.keymap.change,
		m.keymap.quit,
	}))
}

func (m model) View() string {
	s := m.timer.View()
	if m.timer.Timedout() {
		s = "All done!"
	}
	s += "\n"
	if !m.quitting {
		s = m.GetStatus() + "Time:" + s
		s += m.helpView()
	}
	return m.getStyle(s)
}

func main() {
	m := initialModel()
	m.keymap.stop.SetEnabled(false)
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
