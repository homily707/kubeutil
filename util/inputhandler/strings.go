package inputhandler

import (
	"strconv"
	"strings"
)

func Kvsplit(s string) (int, string) {
	split := strings.IndexByte(s, ':')
	if split == -1 {
		return -1, "format wrong, can't find ':'"
	}
	i, err := strconv.Atoi(s[:split])
	if err != nil {
		return -1, "parse index error"
	}
	return i, s[split+1:]
}
