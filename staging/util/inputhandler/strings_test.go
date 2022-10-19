package inputhandler

import (
	"fmt"
	"testing"
)

func Test_Kvsplit(t *testing.T) {
	fmt.Println(Kvsplit("4:hello1"))
	fmt.Println(Kvsplit("4:std::hello1"))
	fmt.Println(Kvsplit(":hello1"))
	fmt.Println(Kvsplit("4:"))
	fmt.Println(Kvsplit("x:hello1"))

}
