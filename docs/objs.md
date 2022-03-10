# 各种对象的定义

## Image
```go
type Image struct {
    // 仅内部使用，用于快速定位当前镜像的 layout 位置
    rootPath string
    // 镜像 layout 的文件夹名
    Hash string
    // manifest.json 中 Config 对应的文件解码后的内容
    Config ImageConfig
    // manifest.json 解码后的内容
    Manifest ManifestType
    // Layer fs 对应的文件夹相对于 rootdir 的路径
    Layers []string
}
```

Image 对象可以通过 Get 方法来创建，内部按需调用 Pull 或者 GetFromLocal，均接受 NameTag 做参数


## Container
```go
type Container struct {
    // TODO
}
```

Container 对象可以通过 Image 的 CreateContainer 来创建，该方法内部调用 NewContainerID 来生成镜像 ID，如果
调用成功，那么会创建好 Container 的 Layout，但是还没有 Mount，该方法返回 containerID

Container 对象有 Run 方法，这里执行 Mount，控制流进入 Container 的 Entrypoint 或者用户指定的进程，该方法接受
命令行参数，用 ...string 来表示，如果不传则去运行 Entrypoint

Container 对象有 Exec 方法