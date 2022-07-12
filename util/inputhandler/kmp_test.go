package inputhandler

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func Test_kmp(t *testing.T) {
	file, _ := os.Open("/Users/hml/codes/kubeutil/go.mod")
	text, _ := ioutil.ReadAll(file)
	var target string
	for {
		fmt.Scanf("%s ", &target)
		indexs := Kmp(text, []byte(target))
		for _, i := range indexs {
			fmt.Println(string(text[max(0, i-20):min(len(text), i+20)]))
		}
	}

}

func max(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

func min(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}
