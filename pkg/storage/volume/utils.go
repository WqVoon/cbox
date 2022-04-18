package volume

import (
	"path/filepath"
	"strings"

	"github.com/wqvoon/cbox/pkg/flags"
	"github.com/wqvoon/cbox/pkg/log"
)

var parsedVolumes = []*Volume{}

// 从 flags 中获取 volume 对象，之所以要做成 Init 是因为要保证命令执行之前能够检验 volume 的格式
// 相应的 parsedVolumes 就不是线程安全的
func Init() {
	volumeArg := flags.GetVolume()
	if len(volumeArg) == 0 {
		return
	}

	// 由于标准库的 flags 不支持数组，所以这里暂时使用逗号分隔
	for _, v := range strings.Split(volumeArg, ",") {
		splitedV := strings.Split(v, ":")
		if len(splitedV) != 2 {
			log.Errorln("error format of volume definition:", v)
		}

		hostPath, containerPath := splitedV[0], splitedV[1]

		hostPath, err := filepath.Abs(hostPath)
		if err != nil {
			log.Errorln("failed to convert hostPath to abs path, err:", err)
		}

		if !filepath.IsAbs(containerPath) {
			log.Errorln("containerPath must be abs path")
		}

		parsedVolumes = append(parsedVolumes, &Volume{
			HostPath:      hostPath,
			ContainerPath: containerPath,
		})
	}
}

func GetVolumes() []*Volume { return parsedVolumes }
