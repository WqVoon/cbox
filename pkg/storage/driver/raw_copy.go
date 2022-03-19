package driver

import (
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"

	"github.com/wqvoon/cbox/pkg/log"
	"github.com/wqvoon/cbox/pkg/utils"
)

var _ Interface = Register(&RawCopy{})

type RawCopy struct{}

// TODO: 这里单纯做复制，暂不考虑 upper layer 删除了 lower layer 中的文件的场景
func (rc *RawCopy) Mount(dst string, layerPaths ...string) {
	for _, layerPath := range layerPaths {
		filepath.WalkDir(layerPath, func(p string, d fs.DirEntry, err error) error {
			if p == layerPath {
				return nil
			}

			dstPath := path.Join(dst, p[len(layerPath):])
			// 如果路径已经存在就跳过，一个 case 是 Container.Stop 后再 Container.Start
			if utils.PathIsExist(dstPath) {
				return nil
			}

			switch {
			case d.IsDir():
				if err := os.MkdirAll(dstPath, 0777); err != nil {
					panic(err)
				}

			case d.Type().IsRegular():
				dstFile, err := os.OpenFile(dstPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
				if err != nil {
					log.Errorln("faild to create file", dstPath, "err:", err)
				}

				srcFile, err := os.OpenFile(p, os.O_RDONLY, 0777)
				if err != nil {
					log.Errorln("faild to open file", p, "err:", err)
				}

				if _, err = io.Copy(dstFile, srcFile); err != nil {
					log.Errorln("faild to copy", srcFile, "to", dstFile, "err:", err)
				}

				dstFile.Close()
				srcFile.Close()

			default:
				// TODO: 这里假设这种文件都是符号链接
				realPath, err := os.Readlink(p)
				if err != nil {
					log.Errorln("faild to read symlink for", dstPath, "err:", err)
				}

				if err := os.Symlink(realPath, dstPath); err != nil {
					// 前面 utils.PathIsExist 返回 false 时才会执行到这里，之所以这里会有 IsExist 的情况出现，是因为容器中的软链接是相对于 / 的
					// 而这里的 / 是经过了 chroot 的，也就是说在宿主机视角下该软链指向的文件实际是位于 container layout/fs/mnt 的
					// 这种位置上的不一致会导致 utils.PathIsExist 返回 false
					if os.IsExist(err) {
						return nil
					}
					log.Errorln("faild to create symlink for", dst, "err:", err)
				}
			}
			return nil
		})
	}
}

// 由于是单纯的复制，因此 UnMount 当前不用做操作，因为目录的清理会交给 Container.Delete
func (rc *RawCopy) UnMount(dst string) {}

func (rc *RawCopy) String() string { return "raw_copy" }
