package driver

import (
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"

	"github.com/wqvoon/cbox/pkg/log"
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
