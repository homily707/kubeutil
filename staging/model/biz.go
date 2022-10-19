package model

import (
	tea "github.com/charmbracelet/bubbletea"
	"kubeutil/staging/util/inputhandler"
	"os/exec"
)

func getFuncAndListNameSpace(c *Context) {
	c.FilterChan(integerParseFilter)
	if err := c.err; err != nil {
		return
	}
	c.nextFunc = c.startFunctions[c.inputIndex].pipeFunc
	// TODO no ns func need to be considered
	c.output = c.kubeClient.ListNamespace()
}

func selectNs(c *Context) {
	err := c.kubeClient.SelectNs(c.inputIndex)
	if err != nil {
		c.err = err
		c.errOutput = err.Error()
	}
}

func getNsThenListPodsToLog(c *Context) {
	c.FilterChan(integerParseFilter, selectNs)
	c.output = c.kubeClient.ListCurNsPods()
	c.nextFunc = logPod
}

func logPod(c *Context) {
	f := func(c *Context) {
		c.output = c.kubeClient.LogPod(c.inputIndex)
		c.nextFunc = searchLog
	}
	c.FilterChan(integerParseFilter, f)
}

func searchLog(c *Context) {
	c.output = c.kubeClient.SearchPod(c.input)
	c.nextFunc = endFunc
}

func getNsThenListPodsToExec(c *Context) {
	c.FilterChan(integerParseFilter, selectNs)
	c.output = c.kubeClient.ListCurNsPods()
	c.nextFunc = execPod
}

func execPod(c *Context) {
	f := func(c *Context) {
		var cmd *exec.Cmd
		c.output, cmd = c.kubeClient.ExecPod(c.inputIndex)
		c.cmd = tea.Exec(tea.WrapExecCommand(cmd), nil)
		c.nextFunc = endFunc
	}
	c.FilterChan(integerParseFilter, f)
}

func getNsThenListConfigMaps(c *Context) {
	c.FilterChan(integerParseFilter, selectNs, listConfigMaps)
}

func listConfigMaps(c *Context) {
	c.output = c.kubeClient.ListConfigMaps()
	c.nextFunc = showConfigData
}

func showConfigData(c *Context) {
	f := func(c *Context) {
		c.output = c.kubeClient.ShowConfigMap(c.inputIndex)
		c.nextFunc = editConfigData
	}
	c.FilterChan(integerParseFilter, f)
}

func editConfigData(c *Context) {
	i, value := inputhandler.Kvsplit(c.input)
	c.output = c.kubeClient.UpdateConfigMap(i, value)
	c.nextFunc = endFunc
}
