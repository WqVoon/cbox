package utils

func ReverseStringSlice(s []string) []string {
	head, tail := 0, len(s)-1
	for head < tail {
		s[head], s[tail] = s[tail], s[head]
		head++
		tail--
	}
	return s
}

func NewStringSlice(strs ...string) []string {
	newSlice := make([]string, 0, len(strs))

	newSlice = append(newSlice, strs...)

	return newSlice
}
