package mytea

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"testing"
)

func Test_main(t *testing.T) {
	p := tea.NewProgram(NewMyModel())
	if err := p.Start(); err != nil {
		fmt.Println("wrong")
	}
}
