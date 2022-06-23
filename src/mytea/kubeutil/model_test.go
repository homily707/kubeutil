package kubeutil

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os"
	//"strconv"
	//"strings"
	"testing"
)

func Test_strsJoin(t *testing.T) {
	strs := []string{"1", "2", "3", "4", "5"}
	fmt.Println(stringsJoin(", ", strs...))
	fmt.Println(stringsJoin("\n", strs...))
	fmt.Println(stringsJoin("\r", strs...))
	fmt.Println(stringsJoin("\t", strs...))
}

func Test_str(t *testing.T) {
	s := "12345"
	println(s)
	s = s + "12345"
	println(s)
}

func Test_Model(t *testing.T) {
	//builder := strings.Builder{}
	//for i := 0; i < 100; i++ {
	//	builder.WriteString(strconv.Itoa(i) + "\n")
	//}
	//content := builder.String()

	p := tea.NewProgram(
		InitScreenModel(),
		tea.WithAltScreen(),       // use the full size of the terminal in its "alternate screen buffer"
		tea.WithMouseCellMotion(), // turn on mouse support so we can track the mouse wheel
	)

	if err := p.Start(); err != nil {
		fmt.Println("could not run program:", err)
		os.Exit(1)
	}
}
