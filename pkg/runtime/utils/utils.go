package utils

import (
	"os"
	"os/exec"

	"github.com/wqvoon/cbox/pkg/log"
)

// var runtimeModePattern = regexp.MustCompile(`^/proc/\d+/exe$`)
const runtimeModePattern = "/proc/self/exe"

// TODO: 除了命令的模式外，还需要检测 pid、挂载点、命名空间等特征
func IsRuntimeMode() bool {
	return len(os.Args) > 1 && runtimeModePattern == os.Args[0]
}

// 从 os.Args 中提取出 exec.Cmd 所需的 path 及 args 参数
func ExtractCmdFromOSArgs() (path string, args []string) {
	switch length := len(os.Args); length {
	case 1, 2:
		log.Errorln("error length of runtimeMode args:", length)

	case 3:
		path, args = os.Args[2], os.Args[2:]

	default:
		path, args = os.Args[3], os.Args[3:]
	}

	var err error
	path, err = exec.LookPath(path)
	if err != nil {
		log.Errorln("failed to lookPath, err:", err)
	}

	return
}
