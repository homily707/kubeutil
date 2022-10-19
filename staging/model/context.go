package model

import (
	tea "github.com/charmbracelet/bubbletea"
	"kubeutil/staging/client"
	"strings"
)

type Context struct {
	input      string
	output     string
	cmd        tea.Cmd
	err        error
	errOutput  string
	inputIndex int

	curFunc        pipeFunc
	nextFunc       pipeFunc
	historyF       []pipeFunc
	historyInputs  []string
	historyOutputs []string

	kubeClient     *client.KubeClient
	startFunctions []pipeFuncWithName

	//util
	builder strings.Builder
}

type pipeFuncWithName struct {
	name string
	pipeFunc
}

func (c *Context) clearHistory() {
	c.historyF = []pipeFunc{}
	c.historyInputs = []string{}
	c.curFunc = rootFunc
	c.input = ""
}

func (c *Context) WriteString(s string) {
	c.builder.WriteString(s)
}

func (c *Context) WriteLine(s string) {
	c.builder.WriteString(s + "\n")
}

func (c *Context) FilterChan(filters ...pipeFunc) {
	for _, filter := range filters {
		filter(c)
		if c.err != nil {
			c.output = c.errOutput
			break
		}
	}
}

func (c *Context) Flush() {
	c.output = c.builder.String()
	c.builder.Reset()
}
