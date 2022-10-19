package controller

import (
	tea "github.com/charmbracelet/bubbletea"
)

const (
	Back = "BACK"
)

var StartFunctions map[string]HandlerFunc

type FuncController struct {
	history
	Context // only support one connection

	curFunc  HandlerFunc
	rootFunc HandlerFunc
}

func (f FuncController) HandleInput(input string) (string, tea.Cmd) {
	// 1. make context
	f.input = input

	// 2. handle back
	if input == Back {
		f.back()
	}

	f.curFunc(&f.Context)
	f.record(f.input, f.output, f.curFunc)
	return f.output, f.cmd
}

type HandlerFunc func(*Context)

type HandlersChain []HandlerFunc

type history struct {
	inputs      []string
	outputs     []string
	funcHistory []HandlerFunc
}

func (h history) record(input, output string, handlerFunc HandlerFunc) {
	h.inputs = append(h.inputs, input)
	h.outputs = append(h.outputs, output)
	h.funcHistory = append(h.funcHistory, handlerFunc)
}

func (h history) back() {

}
