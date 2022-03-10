package image

import "path"

// 获取 Layer 相对于 rootdir 的完整路径
func (manifest *ManifestType) GetLayerFSPaths() []string {
	paths := make([]string, 0, len(manifest.Layers))

	for _, layerPath := range manifest.Layers {
		layerPathDir := path.Dir(layerPath)

		fullPath := path.Join(manifest.rootPath, layerPathDir, "fs")

		paths = append(paths, fullPath)
	}

	return paths
}
