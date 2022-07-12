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
	"kubeutil/util/inputhandler"
	"kubeutil/util/math"
	"log"
	"os/exec"
	"strings"
)

var (
	CurCMS       = "CurrentConfigMaps"
	CurCM        = "CurrentConfigMap"
	CurCMKeyList = "CurrentCMKeyList"
	CurNS        = "CurrentNamespace"
	CurPODS      = "CurrentPods"
	CurPOD       = "CurrentPod"
)

type KubeClient struct {
	*kubernetes.Clientset
	store

	namespaces   []v1.Namespace
	curNamespace string
	pods         []v1.Pod
}

func NewKubeClientFromConfig(filepath string) *KubeClient {
	config, _ := clientcmd.BuildConfigFromFlags("", filepath)
	client, _ := kubernetes.NewForConfig(config)
	return &KubeClient{
		Clientset: client,
		store:     NewConcurrentHashMap()}
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
	c.Set(CurNS, c.namespaces[i].Name)
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
	c.Set(CurPOD, c.pods[i].Name)
	req := c.CoreV1().Pods(c.curNamespace).GetLogs(c.pods[i].Name, &v1.PodLogOptions{})
	body, err := req.Stream(context.Background())
	if err != nil {
		return "get log error" + err.Error()
	}
	var buf bytes.Buffer
	io.Copy(&buf, body)
	return buf.String()
}

func (c KubeClient) SearchPod(word string) string {
	podname := c.Get(CurPOD).(string)
	req := c.CoreV1().Pods(c.curNamespace).GetLogs(podname, &v1.PodLogOptions{})
	body, err := req.Stream(context.TODO())
	if err != nil {
		return "get log error" + err.Error()
	}
	var buf bytes.Buffer
	io.Copy(&buf, body)
	text := buf.Bytes()
	indexs := inputhandler.Kmp(text, []byte(word))

	builder := strings.Builder{}
	for _, i := range indexs {
		builder.Write(text[math.Max(0, i-20):math.Min(len(text), i+20)])
	}
	return builder.String()
}

func (c KubeClient) ExecPod(i int) (string, *exec.Cmd) {
	if i > len(c.pods) {
		return "index out of range", nil
	}
	cmd := exec.Command("kubectl", "exec", "-it", "-n", c.curNamespace,
		c.pods[i].Name, "--", "sh", "-c", "clear;  (bash || sh || ash)")
	return cmd.String(), cmd
}

//kubectl exec -it mysql-mysql-cluster-0 -- sh -c "mysql -u root -h mysql-mysql-cluster -p"
func (c KubeClient) LoginDPSMysql() (string, *exec.Cmd) {
	cmd := exec.Command("kubectl", "exec", "-it", "-n", c.curNamespace,
		"mysql-mysql-cluster-0", "--", "sh", "-c", "mysql -u root -h mysql-mysql-cluster -p")
	return cmd.String(), cmd
}

//kubectl exec -it redis-redis-cluster-0 -- sh -c "redis-cli -h redis-redis-cluster"
func (c KubeClient) LoginDPSRedis() (string, *exec.Cmd) {
	cmd := exec.Command("kubectl", "exec", "-it", "-n", c.curNamespace,
		"redis-redis-cluster-0", "--", "sh", "-c", "redis-cli -h redis-redis-cluster")
	return cmd.String(), cmd
}

func (c *KubeClient) ListConfigMaps() string {
	cmList, _ := c.CoreV1().ConfigMaps(c.curNamespace).List(context.TODO(), metav1.ListOptions{})
	c.Set(CurCMS, cmList)

	build := strings.Builder{}

	build.WriteString("choose a configMap\n")
	for j, item := range cmList.Items {
		build.WriteString(fmt.Sprintf("%d: %s \n", j, item.Name))
	}
	return build.String()
}

func (c *KubeClient) ShowConfigMap(i int) string {
	cmList := c.Get(CurCMS).(*v1.ConfigMapList)
	cm := cmList.Items[i]
	c.Set(CurCM, cm)

	build := strings.Builder{}

	var keyList []string
	build.WriteString("inputhandler [line]:[value] to update. e.g 3:hello\n")
	count := 0
	for k, v := range cm.Data {
		build.WriteString(fmt.Sprintf("%d %s: %s \n", count, k, v))
		count++
		keyList = append(keyList, k)
	}
	c.Set(CurCMKeyList, keyList)
	return build.String()
}

func (c *KubeClient) UpdateConfigMap(i int, value string) string {
	cmKeyList := c.Get(CurCMKeyList).([]string)
	if i > len(cmKeyList) || i < 0 {
		return "index invalid"
	}
	key := cmKeyList[i]
	cm := c.Get(CurCM).(v1.ConfigMap)
	cm.Data[key] = value
	_, err := c.CoreV1().ConfigMaps(c.Get(CurNS).(string)).Update(context.TODO(), &cm, metav1.UpdateOptions{})
	if err != nil {
		return err.Error()
	}
	return "updated"
}
