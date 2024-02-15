package main

import (
	"os/exec"
	"time"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/timer"

	"github.com/charmbracelet/lipgloss"
)

type status int

const (
	start status = iota
	work
	short
	long
)

const (
	workTime   int = 25
	longBreak  int = 15
	shortBreak int = 5
)

type timerModel struct {
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

func (m timerModel) GetStatus() string {
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

func (m timerModel) getStyle(s string) string {
	timerWrapper := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(3, 1).Bold(true).Foreground(lipgloss.Color("#fff")).Background(lipgloss.Color("#74E291")).Align(lipgloss.Center)
	x, y := timerWrapper.GetFrameSize()
	return timerWrapper.Width(20 + x).Height(20 - y).Render(s)
}

func (m timerModel) Notify() {
	notify := exec.Command("notify-send", m.GetStatus())
	sound := exec.Command("paplay", "./getup.mp3")
	notify.Start()
	sound.Start()
}

func initialModel(w, s, l int) timerModel {
	return timerModel{
		status:     start,
		timer:      timer.NewWithInterval(time.Minute*time.Duration(w), time.Second),
		workTime:   w,
		shortBreak: s,
		longBreak:  l,
		keymap:     helpTimerKeys,
		help:       help.New(),
	}
}

func (m timerModel) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m timerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		m.keymap.stop.SetEnabled(false)
		return m, nil
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.quit):
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, m.keymap.reset):
			m.timer.Timeout = time.Minute * time.Duration(m.workTime)
			m.timer.Stop()

		case key.Matches(msg, m.keymap.change):
			return inputs(inputsInitialModel()), nil
		case key.Matches(msg, m.keymap.start):
			m.status = work
			return m, m.timer.Start()
		case key.Matches(msg, m.keymap.stop):
			return m, m.timer.Stop()
		}

	}
	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m timerModel) helpView() string {
	return helpStyle("\n" + m.help.ShortHelpView([]key.Binding{
		m.keymap.start,
		m.keymap.stop,
		m.keymap.reset,
		m.keymap.change,
		m.keymap.quit,
	}))
}

func (m timerModel) View() string {
	s := m.timer.View()
	s += "\n"
	s = m.getStyle(m.GetStatus() + "Time:" + s)
	w, h := wrapperStyle.GetFrameSize()
	return wrapperStyle.Width(50 + w).Height(10 - h).Render(lipgloss.JoinVertical(lipgloss.Center, "Pomodoro CLI App", s, m.helpView()))
}
