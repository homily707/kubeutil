package biz

import (
	tea "github.com/charmbracelet/bubbletea"
	"kubeutil/internal/controller"
	"kubeutil/staging/util/inputhandler"
	"strconv"
)

func rootFunc(c *controller.Context) {
	c.WriteLine("choose function")
	for k, v := range controller.StartFunctions {

	}
	c.NextFunc = getFuncAndListNameSpace
}

func getFuncAndListNameSpace(c *controller.Context) {
	c.nextFunc = StartFunctions[c.inputIndex].pipeFunc
	// TODO no ns func need to be considered
	c.output = c.kubeClient.ListNamespace()
}

func selectNs(c *controller.Context) {
	err := c.kubeClient.SelectNs(c.inputIndex)
	if err != nil {
		c.err = err
		c.errOutput = err.Error()
	}
}

func GetNsThenListPodsToLog(c *controller.Context) {
	c.FilterChan(integerParseFilter, selectNs)
	c.output = c.kubeClient.ListCurNsPods()
	c.nextFunc = logPod
}

func logPod(c *controller.Context) {
	f := func(c *controller.Context) {
		c.output = c.kubeClient.LogPod(c.inputIndex)
		c.nextFunc = searchLog
	}
	c.FilterChan(integerParseFilter, f)
}

func searchLog(c *controller.Context) {
	c.output = c.kubeClient.SearchPod(c.input)
	c.nextFunc = endFunc
}

func GetNsThenListPodsToExec(c *controller.Context) {
	c.FilterChan(integerParseFilter, selectNs)
	c.output = c.kubeClient.ListCurNsPods()
	c.nextFunc = execPod
}

func execPod(c *controller.Context) {
	f := func(c *controller.Context) {
		var cmd *exec.Cmd
		c.output, cmd = c.kubeClient.ExecPod(c.inputIndex)
		c.cmd = tea.Exec(tea.WrapExecCommand(cmd), nil)
		c.nextFunc = endFunc
	}
	c.FilterChan(integerParseFilter, f)
}

func GetNsThenListConfigMaps(c *controller.Context) {
	c.FilterChan(integerParseFilter, selectNs, listConfigMaps)
}

func listConfigMaps(c *controller.Context) {
	c.output = c.kubeClient.ListConfigMaps()
	c.nextFunc = showConfigData
}

func showConfigData(c *controller.Context) {
	f := func(c *controller.Context) {
		c.output = c.kubeClient.ShowConfigMap(c.inputIndex)
		c.nextFunc = editConfigData
	}
	c.FilterChan(integerParseFilter, f)
}

func editConfigData(c *controller.Context) {
	i, value := inputhandler.Kvsplit(c.input)
	c.output = c.kubeClient.UpdateConfigMap(i, value)
	c.nextFunc = endFunc
}
