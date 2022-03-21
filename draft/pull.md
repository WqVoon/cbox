# 下载一个镜像所需的操作
> cbox pull busybox:latest

暂时用 crane v0.1.1 来做镜像的下载，后面可以去掉它

新建 image/cmd.go 文件，添加 Pull 方法，接受 NameTag 做参数

cbox-dir 中新增 tarballs 文件夹

ImageIdx 添加 Update 方法，接受 NameTag 和 ImgHash 做参数，返回布尔值表示是否更新

- 调用 crane.Pull 拿到 img 对象
- 调用 img.Manifest 拿到 manifest 对象
- 使用 manifest.Config.Digest.Hex[:12] 获取镜像的 hash
- 调用 ImageIdx.Update 方法更新索引，如果返回 false 那么输出提示信息后退出即可
- 调用 crane.SaveLegacy 将 tarball 保存到 tarballs 文件夹
- 在 rootdir/images 文件夹中新建一个以 hash 为名的文件夹
- 将 tarball 解析到 hash 文件夹中
- 使用 hash 拿到 cbox 的 manifest 对象
- 遍历 manifest.layer，拆分出文件夹名，在里面新建 fs 文件夹
- 将 layer 中的每个路径指向的 tarball 解析到上一步的 fs 中