package controller

import tea "github.com/charmbracelet/bubbletea"

type Controller interface {
	HandleInput(string) (string, tea.Cmd)
}
