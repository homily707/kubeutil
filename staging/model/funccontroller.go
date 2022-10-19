package model

import (
	tea "github.com/charmbracelet/bubbletea"
	"kubeutil/staging/client"
	"os"
	"path/filepath"
	"strconv"
)

type funcController struct {
	context *Context
}

func NewFuncController() *funcController {
	home := os.Getenv("HOME")
	kubeConfig := filepath.Join(home, ".kube", "config")

	controller := funcController{
		context: &Context{
			kubeClient: client.NewKubeClientFromConfig(kubeConfig),
			startFunctions: []pipeFuncWithName{
				{"LOG", getNsThenListPodsToLog},
				{"EXEC", getNsThenListPodsToExec},
				{"CONFIG", getNsThenListConfigMaps},
			},
		},
	}
	controller.context.clearHistory()
	return &controller
}

type pipeFunc func(*Context)

func (f *funcController) HandleInput(s string) (string, tea.Cmd) {
	f.context.input = s
	f.context.err = nil
	f.context.cmd = nil
	if s == back {
		backFunc(f.context)
	}

	// execute
	f.context.curFunc(f.context)

	//history control
	f.context.historyInputs = append(f.context.historyInputs, s)
	f.context.historyF = append(f.context.historyF, f.context.curFunc)
	f.context.curFunc = f.context.nextFunc

	return f.context.output, f.context.cmd
}

func rootFunc(c *Context) {
	c.WriteLine("choose function")
	for i, f := range c.startFunctions {
		c.WriteLine(strconv.Itoa(i) + ": " + f.name)
	}
	c.nextFunc = getFuncAndListNameSpace
	c.Flush()
}

func backFunc(c *Context) {
	if len(c.historyF) != len(c.historyInputs) {
		c.clearHistory()
		return
	}
	n := len(c.historyF)
	if n > 1 {
		c.curFunc = c.historyF[n-2]
		c.input = c.historyInputs[n-2]
		c.historyInputs = c.historyInputs[:n-2]
		c.historyF = c.historyF[:n-2]
	} else {
		c.clearHistory()
	}
}

func endFunc(c *Context) {

}
