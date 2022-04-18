package utils

import (
	"path"

	"github.com/wqvoon/cbox/pkg/flags"
	"github.com/wqvoon/cbox/pkg/rootdir"
	runtimeInfo "github.com/wqvoon/cbox/pkg/runtime/info"
	"github.com/wqvoon/cbox/pkg/storage/volume"
)

// 将 flags 中传递的 volumes 记录在 containerInfo 中
func Record(containerId string) {
	mntPath := rootdir.GetContainerMountPath(containerId)

	pathPairs := flags.GetVolumes()
	vs := make([]*volume.Volume, 0, len(pathPairs))

	for _, pathPair := range pathPairs {
		hostPath, containerPath := pathPair[0], pathPair[1]

		vs = append(vs, &volume.Volume{
			HostPath:      hostPath,
			ContainerPath: path.Join(mntPath, containerPath),
		})
	}

	runtimeInfo.GetContainerInfo(containerId).SaveVolumes(vs)
}
