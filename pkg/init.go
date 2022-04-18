package pkg

import (
	"github.com/wqvoon/cbox/pkg/flags"
	"github.com/wqvoon/cbox/pkg/log"
	"github.com/wqvoon/cbox/pkg/rootdir"
	"github.com/wqvoon/cbox/pkg/storage/driver"
	"github.com/wqvoon/cbox/pkg/storage/volume"
)

func Init() {
	// 首先解析 flags 参数
	flags.Init()
	// rootdir 通过 flags 参数可以得到 cbox-dir 的根目录
	rootdir.Init()
	// log 通过 flags 参数决定是否使用 debug 模式
	log.Init()
	// driver 通过 flags 参数决定使用哪个 storageDriver
	driver.Init()
	// volume 通过 flags 参数决定使用哪些 volume
	volume.Init()
}
