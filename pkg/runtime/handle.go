package runtime

import (
	"os"
	"os/exec"
	"strings"

	"github.com/wqvoon/cbox/pkg/container"
	"github.com/wqvoon/cbox/pkg/log"
	"github.com/wqvoon/cbox/pkg/rootdir"
	runtimeUtils "github.com/wqvoon/cbox/pkg/runtime/utils"
	"github.com/wqvoon/cbox/pkg/utils"
	"golang.org/x/sys/unix"
)

// 这个函数只能在 RuntimeMode 下使用
func Handle() {
	log.Println("Enter RuntimeMode, args:", os.Args)

	c := container.GetContainerByID(os.Args[1])

	// TODO: 待补充其他的 ns
	namespaces := []string{"/pid", "/uts", "/ipc"}
	srcPathPrefix := "/proc/self/ns"
	dstPathPrefix := rootdir.GetContainerNSPath(c.ID)
	for _, ns := range namespaces {
		src := srcPathPrefix + ns
		dst := dstPathPrefix + ns
		if err := unix.Mount(src, dst, "", unix.MS_BIND, ""); err != nil {
			log.Errorf("failed to bind mount %q to %q, err: %v\n", src, dst, err)
		}
	}

	// 需要保证在 ExtractCmdFromOSArgs 前进行 Env 的处理，这样得到的 path 才是正确的
	os.Clearenv()
	for _, oneEnv := range c.Env {
		envPair := strings.Split(oneEnv, "=")
		key, val := envPair[0], envPair[1]
		os.Setenv(key, val)
	}

	// 这里需要保证在 ExtractCmdFromOSArgs 前进行 chroot，这样得到的 path 才是正确的
	unix.Chroot(rootdir.GetContainerMountPath(c.ID))

	// TODO: Mount 的第一个参数如果留空则宿主机上会因为解析错误而读不到这条记录，也许可以利用下
	utils.CreateDirIfNotExist("/proc")
	if err := unix.Mount("cbox-proc", "/proc", "proc", 0, ""); err != nil {
		log.Errorln("faild to mount /proc, err:", err)
	}

	path, args := runtimeUtils.ExtractCmdFromOSArgs()
	cmd := &exec.Cmd{
		Path: path,
		Args: args,
		Dir:  "/",
		Env:  c.Env,

		// TODO: 这部分的赋值应该可选
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	if err := cmd.Run(); err != nil {
		// 这里不进行 Errorln，用于避免 Ctrl+C 后 Ctrl+D 引起的常见错误
		log.Println("an error may have occurred while running container:", err)
	}
	os.Exit(0)
}
