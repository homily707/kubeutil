package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"kubeutil/src/mytea/kubeutil"
	"os"
	//"strconv"
	//"strings"
)

func main() {
	//builder := strings.Builder{}
	//for i := 0; i < 100; i++ {
	//	builder.WriteString(strconv.Itoa(i) + "\n")
	//}
	//content := builder.String()

	p := tea.NewProgram(
		kubeutil.InitScreenModel(),
		tea.WithAltScreen(),       // use the full size of the terminal in its "alternate screen buffer"
		tea.WithMouseCellMotion(), // turn on mouse support so we can track the mouse wheel
	)

	if err := p.Start(); err != nil {
		fmt.Println("could not run program:", err)
		os.Exit(1)
	}
}
