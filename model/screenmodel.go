package model

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

const (
	INPUTMODE ModelStatus = iota
	VISUALMODE
)

type (
	ModelStatus int
)

type ScreenModel struct {
	view    viewport.Model
	input   textinput.Model
	ready   bool
	status  ModelStatus
	content string
}

var (
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return titleStyle.Copy().BorderStyle(b)
	}()

	statusMap = map[ModelStatus]string{
		INPUTMODE:  "--INSERT--",
		VISUALMODE: "[VISUAL]",
	}

	//_ ScreenModel = (tea.Model)(nil)
)

func NewScreenModel() *ScreenModel {
	var m ScreenModel
	m.ready = false
	m.status = INPUTMODE
	m.view = viewport.Model{}
	m.input = textinput.New()
	return &m
}

func (m ScreenModel) Init() tea.Cmd {
	return nil
}

func (m *ScreenModel) Update(msg tea.Msg) tea.Cmd {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		inputHeight := lipgloss.Height(m.inputView())
		verticalMarginHeight := headerHeight + footerHeight + inputHeight

		if m.status == INPUTMODE {
			m.input.Focus()
		} else if m.status == VISUALMODE {
			m.input.Reset()
			m.input.Blur()
		}
		if !m.ready {
			m.view = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.view.YPosition = headerHeight
			m.view.SetContent(m.content)
			m.ready = true

			m.view.YPosition = headerHeight + 1
		} else {
			m.view.Width = msg.Width
			m.view.Height = msg.Height - verticalMarginHeight
		}
	}

	m.view, cmd = m.view.Update(msg)
	cmds = append(cmds, cmd)
	m.input, cmd = m.input.Update(msg)
	cmds = append(cmds, cmd)
	return tea.Batch(cmds...)
}

func (m ScreenModel) View() string {
	return fmt.Sprintf("%s\n%s\n%s\n%s", m.headerView(), m.view.View(), m.footerView(), m.inputView())
	//return stringsJoin("\n", m.headerView(), m.View(), m.footerView(), m.inputView())
}

func (m ScreenModel) headerView() string {
	builder := strings.Builder{}
	switch m.status {
	case INPUTMODE:
		builder.WriteString("--INSERT--")
		builder.WriteString("ESC to VISUAL MODE, ")
	case VISUALMODE:
		builder.WriteString("[VISUAL]:")
		builder.WriteString("i to INSERT MODE, 'b' or ESC to last step, ")
	}
	builder.WriteString("ctrl+c to exit")
	title := titleStyle.Render(builder.String())
	line := strings.Repeat("─", max(0, m.view.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m ScreenModel) footerView() string {
	info := infoStyle.Render(fmt.Sprintf("%3.f%%", m.view.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, m.view.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

func (m ScreenModel) inputView() string {
	return m.input.View()
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
