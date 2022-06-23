package client

import "k8s.io/client-go/kubernetes"

type KubeClient struct {
	kubernetes.Clientset

	curNamespace string
}

func (c KubeClient) ListNamespace() string {

}

func (c *KubeClient) SelectNs(i int) {

}

func (c KubeClient) ListCurNsPods() string {

}
