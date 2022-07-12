package inputhandler

func prefix(str []byte) []int {
	n := len(str)
	pre := make([]int, n)
	pre[0] = 0
	for i := 1; i < n; i++ {
		j := pre[i-1]
		for j > 0 && str[j] != str[i] {
			j = pre[j-1]
		}
		if str[i] == str[j] {
			j++
		}
		pre[i] = j
	}
	return pre
}

func Kmp(text []byte, target []byte) []int {
	pre := prefix(target)
	i, j := 0, 0
	var ans []int

	for i < len(text) {
		if text[i] == target[j] {
			i++
			j++
		} else if j == 0 {
			i++
		} else if j != 0 {
			j = pre[j-1]
		}

		if j == len(target) {
			ans = append(ans, i-len(target))
			j = pre[j-1]
		}
	}
	return ans
}
