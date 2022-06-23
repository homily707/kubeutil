package model

import (
	tea "github.com/charmbracelet/bubbletea"
	"kubeutil/client"
	"strconv"
)

const (
	LOG ProType = iota + 1
	EXEC
	EDIT_CONFIG
)

const RootPath = ""

type ProType int

type Controller interface {
	HandleInput(string) (string, tea.Cmd)
}

type cmdFunc func(string) (string, tea.Cmd)

func nilCmdWrap(f func(string) string) cmdFunc {
	return func(s string) (string, tea.Cmd) {
		return f(s), nil
	}
}

type KubeController struct {
	kubeClient *client.KubeClient

	curPath string
	protype ProType
	route   map[string]cmdFunc
	backmap map[string]string
	history []string
}

func (c *KubeController) HandleInput(value string) (string, tea.Cmd) {
	if value == back {
		c.curPath = c.backmap[c.curPath]
		if n := len(c.history); n > 1 {
			value = c.history[n-2]
			c.history = c.history[:n]
		} else {
			value = ""
		}

	} else {
		c.history = append(c.history, value)
	}
	return c.route[c.curPath](value)
}

func (c *KubeController) addRoute(s string, f cmdFunc) {
	c.route[s] = f
}

func NewKubeController(kubeClient client.KubeClient, route map[string]cmdFunc, backmap map[string]string) *KubeController {
	c := &KubeController{
		kubeClient: kubeClient,
		curPath:    RootPath,
		route:      map[string]cmdFunc{},
		backmap:    map[string]string{},
	}

	c.addRoute(RootPath, nilCmdWrap(c.listFunction))
	c.addRoute("/func", nilCmdWrap(c.getFuncThenListNamespace))
	c.addRoute("/func/ns", nilCmdWrap(c.getNsThenListChoice))
	c.addRoute("/func/ns/log", nilCmdWrap(logPod))
	c.addRoute("/func/ns/exec", execPod)

	c.backmap[RootPath] = RootPath
	c.backmap["/func"] = RootPath
	c.backmap["/func/ns"] = RootPath
	c.backmap["/func/ns/log"] = "/func"
	c.backmap["/func/ns/exec"] = "/func"

	return c
}

func (c *KubeController) listFunction(input string) string {
	c.curPath = c.curPath + "/func"
	return "choose function \n" +
		"1: log \n" +
		"2: exec"
}

func (c *KubeController) getFuncThenListNamespace(input string) string {
	i, err := strconv.Atoi(input)
	if err != nil {
		return "parse index error"
	}
	c.curPath = c.curPath + "/ns"
	c.protype = ProType(i)
	return c.kubeClient.ListNamespace()
}

func (c *KubeController) getNsThenListChoice(input string) string {
	i, err := strconv.Atoi(input)
	if err != nil {
		return "parse index error"
	}
	c.kubeClient.SelectNs(i)
	switch c.protype {
	case LOG, EXEC:
		return c.kubeClient.ListCurNsPods()
	}
	return "something wrong"
}

func logPod(input string) string {

}

func execPod(input string) (string, tea.Cmd) {

}
