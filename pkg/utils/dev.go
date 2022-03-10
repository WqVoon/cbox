package utils

func TODO(msg ...string) {
	if len(msg) == 0 {
		panic("TODO")
	} else {
		panic("TODO: " + msg[0])
	}
}
