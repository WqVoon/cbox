# 运行一个容器所需的操作（无 Cgroups 和 Network 相关）
> cbox run busybox /bin/sh [...args]
- preflight
  - 检查是否是 root 用户（需要 chroot 操作，后面可以尝试用 rootless 技术解决）
  - 补齐镜像名字，拉取镜像到 rootDir/images 内
  - 创建一个随机字符串作为 containerID
  - 解析 image 的 Env，entrypoint 之类的
  - 在 rootDir/containers 中创建镜像的目录，mnt 作为最终目录，workdir 作为上层目录，upper 作为镜像目录（TODO：Overlay 存储引擎布局，其他存储引擎再议）
  - 绑定挂载容器到 rootDir/containers/containerID/mnt 目录下
- run
  - 设置 hostName（这个提前做，用来判断是否进入了容器）
  - chroot 到 rootDir/containers/containerID/mnt
  - 切换目录到 /（此时的 / 其实就是 rootDir/containers/containerID/mnt）
  - 设置 Env 为镜像中的 Env
  - 挂载一些特殊目录（TODO：需要知道含义）
    - /dev/pts: devpts
    - /sys: sysfs
    - /proc: proc
    - /tmp: tmpfs
  - 运行用户指定的程序或者 entrypoint
- clean
  - 卸载容器的 mnt 目录
  - 删除 rootDir/containers/containerID 中的全部内容