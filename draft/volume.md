# 数据卷相关的设计

```golang
type Volume struct {
	// 对于容器而言是否只读
	ReadOnly bool
	// 在宿主机上的路径（绝对路径）
	HostPath string
	// 在容器中挂载的路径（绝对路径）
	ContainerPath string
}
```

ContainerInfo 中新增 Volumes 字段，是 Volume 对象的数组

容器在 `cbox create` 时通过调用 `volume.Record` 解析 `--volume` 参数来获取数据卷，参数对应的值格式如下：
`<HostPath>:<ContainerPath>`，这样的格式算作一个 Volume，可以用 `,` 分隔出多个 Volume

对于每一个 Volume 生成一个 Volume 对象，并记录在 ContainerInfo 中

runtime 在容器启动时针对 ContainerInfo 中的每一个 Volume 调用其 Mount 方法

runtime 在容器关闭时针对 ContainerInfo 中的每一个 Volume 调用其 Unmount 方法