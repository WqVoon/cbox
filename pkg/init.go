package pkg

import (
	"github.com/wqvoon/cbox/pkg/cgroups"
	"github.com/wqvoon/cbox/pkg/config"
	"github.com/wqvoon/cbox/pkg/flags"
	"github.com/wqvoon/cbox/pkg/log"
	"github.com/wqvoon/cbox/pkg/network"
	"github.com/wqvoon/cbox/pkg/rootdir"
	"github.com/wqvoon/cbox/pkg/storage/driver"
	"github.com/wqvoon/cbox/pkg/storage/volume"
)

func Init() {
	// 首先解析 flags 参数
	flags.Init()
	// rootdir 通过 flags 参数可以得到 cbox-dir 的根目录
	rootdir.Init()
	// config 使用了 root_dir，所以需要放在 rootdir.Init 后面
	config.Init()
	// log 通过 flags 参数与 config 决定是否使用 debug 模式
	log.Init()
	// driver 通过 flags 参数与 config 决定使用哪个 storageDriver
	driver.Init()
	// volume 通过 flags 参数决定使用哪些 volume
	volume.Init()
	// network 使用了 log 和 config，所以需要放在后面
	network.Init()
	// cgroups 使用了 config，所以需要放在 config.Init 后面
	cgroups.Init()
}
