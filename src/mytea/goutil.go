package mytea

import (
	"fmt"
	//bubbles "github.com/charmbracelet/bubbles"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"reflect"
	"strings"
)

type Mymodel struct {
	str    string
	inputs string
	input  textinput.Model
}

func NewMyModel() Mymodel {
	ti := textinput.New()
	ti.Focus()
	ti.CharLimit = 128
	ti.SetCursorMode(textinput.CursorStatic)
	return Mymodel{input: ti}
}

func (m Mymodel) Init() tea.Cmd {
	return nil
}

func (m Mymodel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("msg type: %s", reflect.TypeOf(msg)))
	builder.WriteString("\n")

	var ticmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		builder.WriteString(fmt.Sprintf("type: %s, string: %s \n", msg.Type, msg.String()))
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "enter":
			builder.WriteString("get input: ")
			builder.WriteString(m.input.Value())
			m.str = builder.String()
			m.input.Reset()
			return m, nil
		}
	}
	m.input, ticmd = m.input.Update(msg)

	return m, ticmd
}

func (m Mymodel) View() string {
	builder := strings.Builder{}
	builder.WriteString(m.str)
	builder.WriteString("\n         ||||    \n")
	builder.WriteString(m.input.View())
	builder.WriteString("\n")
	return builder.String()
}
