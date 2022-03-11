package driver

import (
	"github.com/wqvoon/cbox/pkg/flags"
	"github.com/wqvoon/cbox/pkg/log"
)

type Interface interface {
	// dst 是 Container 的 fs 目录，Mount 方法将 layerPaths 指明的镜像层一起挂载到 dst 上
	// TODO：layerPaths 的顺序由谁控制
	Mount(dst string, layerPaths ...string)

	// 卸载 Container 的 fs
	UnMount(dst string)

	// 返回 Driver 的字符串形式
	String() string
}

var (
	registeredDrivers map[string]Interface

	// 本次运行中被使用的 Driver，被 init 方法赋值
	D Interface
)

// Register 会被具体的 StorageDriver 实现在启动时调用
// 在向全局注册的同时验证对应的实现是否满足接口的需求，类似于 `var _ Interface = Obj{}` 的编译验证方式
func Register(driver Interface) Interface {
	if registeredDrivers == nil {
		registeredDrivers = make(map[string]Interface)
	}

	registeredDrivers[driver.String()] = driver

	return driver
}

func init() {
	flags.ParseAll()

	driverName := flags.GetStorageDriver()
	var isIn bool

	D, isIn = registeredDrivers[driverName]
	if !isIn {
		log.Errorln("no such storage driver:", driverName)
	}
}
