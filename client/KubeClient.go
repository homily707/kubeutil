package client

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os/exec"
	"strings"
)

type KubeClient struct {
	*kubernetes.Clientset

	namespaces   []v1.Namespace
	curNamespace string
	pods         []v1.Pod
}

func NewKubeClientFromConfig(filepath string) *KubeClient {
	config, _ := clientcmd.BuildConfigFromFlags("", filepath)
	client, _ := kubernetes.NewForConfig(config)
	return &KubeClient{Clientset: client}
}

func (c *KubeClient) ListNamespace() string {
	list, err := c.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}
	c.namespaces = list.Items

	build := strings.Builder{}
	build.WriteString("choose a namespace\n")
	for i, item := range c.namespaces {
		build.WriteString(fmt.Sprintf("%d: %s \n", i, item.Name))
	}
	return build.String()
}

func (c *KubeClient) SelectNs(i int) error {
	if i > len(c.namespaces) {
		return errors.New("index out of range")
	}
	c.curNamespace = c.namespaces[i].Name
	return nil
}

func (c *KubeClient) ListCurNsPods() string {
	podList, _ := c.CoreV1().Pods(c.curNamespace).List(context.TODO(), metav1.ListOptions{})
	c.pods = podList.Items
	build := strings.Builder{}

	build.WriteString("choose a pod\n")
	for j, item := range c.pods {
		build.WriteString(fmt.Sprintf("%d: %s \n", j, item.Name))
	}
	return build.String()
}

func (c KubeClient) LogPod(i int) string {
	if i > len(c.pods) {
		return "index out of range"
	}
	req := c.CoreV1().Pods(c.curNamespace).GetLogs(c.pods[i].Name, &v1.PodLogOptions{})
	body, err := req.Stream(context.Background())
	if err != nil {
		return "get log error" + err.Error()
	}
	var buf bytes.Buffer
	io.Copy(&buf, body)
	return buf.String()
}

func (c KubeClient) ExecPod(i int) (string, *exec.Cmd) {
	if i > len(c.pods) {
		return "index out of range", nil
	}
	cmd := exec.Command("kubectl", "exec", "-it", "-n", c.curNamespace,
		c.pods[i].Name, "--", "sh", "-c", "clear; (bash || sh || ash)")
	return "wait a moment", cmd
}
