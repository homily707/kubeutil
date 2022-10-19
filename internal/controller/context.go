package controller

import (
	tea "github.com/charmbracelet/bubbletea"
	"strings"
)

type Context struct {
	input string
	output string // todo use writer to refactor
	cmd tea.Cmd

	NextFunc HandlerFunc
	builder  strings.Builder
}

func (c *Context) WriteString(s string) {
	c.builder.WriteString(s)
}

func (c *Context) WriteLine(s string) {
	c.builder.WriteString(s + "\n")
}

func (c *Context)
