package model

import (
	tea "github.com/charmbracelet/bubbletea"
)

const (
	back = "BACK"
)

type (
	RootModel struct {
		screen     *ScreenModel
		controller Controller
	}
)

func NewRootModel() RootModel {
	return RootModel{
		screen: NewScreenModel(),
	}
}

func (r RootModel) Init() tea.Cmd {
	return nil
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
			case "esc":
				// change to visual
				r.screen.status = VISUALMODE
				r.screen.input.Reset()
				r.screen.input.Blur()
			case "enter":
				cmd = r.handleInput(r.screen.input.Value())
				if cmd != nil {
					return r, cmd
				}
				cmds = append(cmds, cmd)
				r.screen.input.Reset()
				r.screen.view.GotoBottom()
			}
		} else if r.screen.status == VISUALMODE {
			switch msg.String() {
			case "b", "esc":
				cmd = r.handleInput(back)
			case "i":
				r.screen.status = INPUTMODE
				r.screen.input.Focus()
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
