package utils

// ParseCmd 用于从数组中解析出 cmd，该函数的返回值可直接传递给 exec.Command 函数
func ParseCmd(input ...string) (name string, args []string) {
	if len(input) == 0 {
		return
	}

	if len(input) == 1 {
		return input[0], nil
	}

	return input[0], input[1:]
}
