package model

import (
	tea "github.com/charmbracelet/bubbletea"
	"kubeutil/client"
	"os"
	"path/filepath"
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

func NewKubeController() *KubeController {
	home := os.Getenv("HOME")
	kubeconfig := filepath.Join(home, ".kube", "config")

	c := &KubeController{
		kubeClient: client.NewKubeClientFromConfig(kubeconfig),
		curPath:    RootPath,
		route:      map[string]cmdFunc{},
		backmap:    map[string]string{},
	}

	c.addRoute(RootPath, nilCmdWrap(c.listFunction))
	c.addRoute("/func", nilCmdWrap(c.getFuncThenListNamespace))
	c.addRoute("/func/ns", nilCmdWrap(c.getNsThenListChoice))
	c.addRoute("/func/ns/log", nilCmdWrap(c.logPod))
	c.addRoute("/func/ns/exec", c.execPod)

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
	if err := c.kubeClient.SelectNs(i); err != nil {
		return err.Error()
	}
	switch c.protype {
	case LOG:
		c.curPath = c.curPath + "/log"
		return c.kubeClient.ListCurNsPods()
	case EXEC:
		c.curPath = c.curPath + "/exec"
		return c.kubeClient.ListCurNsPods()
	}
	return "something wrong"
}

func (c *KubeController) logPod(input string) string {
	i, err := strconv.Atoi(input)
	if err != nil {
		return "parse index error"
	}
	return c.kubeClient.LogPod(i)
}

func (c *KubeController) execPod(input string) (string, tea.Cmd) {
	i, err := strconv.Atoi(input)
	if err != nil {
		return "parse index error", nil
	}
	s, cmd := c.kubeClient.ExecPod(i)
	return s, tea.Exec(tea.WrapExecCommand(cmd), nil)
}
