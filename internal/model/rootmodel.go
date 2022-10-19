package model

import (
	tea "github.com/charmbracelet/bubbletea"
	"kubeutil/internal/controller"
)

type (
	RootModel struct {
		screen     *ScreenModel
		controller controller.Controller
	}
)

func NewRootModel() RootModel {
	return RootModel{
		screen:     NewScreenModel(),
		controller: controller.FuncController{},
	}
}

func (r RootModel) Init() tea.Cmd {
	return func() tea.Msg {
		return tea.Key{Type: tea.KeyEnter}
	}
}

func (r RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if r.screen.status == INPUTMODE {
			switch msg.String() {
			case "ctrl+c":
				return r, tea.Quit
			case "esc":
				// Back to last step
				cmd = r.handleInput(Back)
				r.screen.input.Reset()
			case "enter":
				cmd = r.handleInput(r.screen.input.Value())
				if cmd != nil {
					return r, cmd
				}
				cmds = append(cmds, cmd)
				r.screen.input.Reset()
				r.screen.view.GotoBottom()
			}
		}
	}
	cmd = r.screen.Update(msg)
	cmds = append(cmds, cmd)

	return r, tea.Batch(cmds...)
}

func (r RootModel) View() string {
	return r.screen.View()
}

func (r RootModel) handleInput(value string) tea.Cmd {
	str, cmd := r.controller.HandleInput(value)
	r.screen.view.SetContent(str)
	return cmd
}
