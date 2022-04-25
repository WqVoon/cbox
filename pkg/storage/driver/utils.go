package driver

import (
	"os"

	"github.com/wqvoon/cbox/pkg/config"
	"github.com/wqvoon/cbox/pkg/flags"
	"github.com/wqvoon/cbox/pkg/log"
)

func Init() {
	const driverEnvName = "CBOX_STORAGE_DRIVER"
	defaultDriverName := config.GetDriverName()

	// 命令行的 flag 优先
	driverName := flags.GetStorageDriver()

	// 如果为空，那么尝试获取环境变量
	if len(driverName) == 0 {
		driverName = os.Getenv(driverEnvName)
	}

	// 如果为空，那么设置成默认值
	if len(driverName) == 0 {
		driverName = defaultDriverName
	}

	// 这里对一些不合法的内容会做校验，所以上面仅检测是否为空即可
	D = GetDriverByName(driverName)
}

// Register 会被具体的 StorageDriver 实现在启动时调用
// 在向全局注册的同时验证对应的实现是否满足接口的需求，类似于 `var _ Interface = Obj{}` 的编译验证方式
func Register(driver Interface) Interface {
	if registeredDrivers == nil {
		registeredDrivers = make(map[string]Interface)
	}

	registeredDrivers[driver.String()] = driver

	return driver
}

func GetDriverByName(driverName string) Interface {
	driver, isIn := registeredDrivers[driverName]
	if !isIn {
		log.Errorln("no such storage driver:", driverName)
	}

	return driver
}
