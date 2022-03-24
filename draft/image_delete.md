## 删除一个 Image 时要做哪些操作

在 Image 的 Layout 中创建 info 文件，其中保存一个 json 对象， UsedBy 参数是一个 []string，内部记录使用这个 Image 的 Container.ID

ImageInfo.MarkUsedBy 方法添加一个 containerID 到 UsedBy 数组中

ImageInfo.MarkReleasedBy 方法从 UsedBy 中删除一个 containerID

ImageInfo.CanBeDeleted 判断是否可以删除，当前仅判断 UsedBy 是否为空

提供 GetImageInfo 用于获取 ImageInfo 对象，如果 info 文件不存在那么创建它

- 创建容器时调用 ImageInfo.MarkUsedBy
- 删除容器时调用 ImageInfo.MarkReleasedBy
- 删除镜像前检查 ImageInfo.CanBeDeleted 的返回值