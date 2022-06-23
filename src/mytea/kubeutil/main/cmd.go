package main

import (
	"bytes"
	"context"
	"io"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
)

func main1() {
	home := os.Getenv("HOME")
	kubeconfig := filepath.Join(home, ".kube", "config")
	config, _ := clientcmd.BuildConfigFromFlags("", kubeconfig)
	client, _ := kubernetes.NewForConfig(config)

	podsDefault := client.CoreV1().Pods("default")
	podList, _ := podsDefault.List(context.TODO(), metav1.ListOptions{})
	pod := podList.Items[0]
	req := client.CoreV1().Pods("default").GetLogs(pod.Name, &v1.PodLogOptions{})
	body, _ := req.Stream(context.Background())
	var buf bytes.Buffer
	io.Copy(&buf, body)
	buf.String()
	body.Close()
	//for _, pod := range podList.Items {
	//	fmt.Println(pod.Name)
	//}
}
