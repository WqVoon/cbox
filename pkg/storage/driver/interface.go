package driver

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
	// 保存所有注册过的 Driver，key 是 DriverName，value 是对应的 Interface 实现
	registeredDrivers map[string]Interface

	// 本次运行中被使用的 Driver，被 Init 方法赋值
	D Interface
)
