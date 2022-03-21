package utils

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"

	"github.com/wqvoon/cbox/pkg/log"
)

// 把压缩包 tarball 解压到 target 指定的文件夹，target 如果不存在则会创建
func Untar(tarball, target string) {
	CreateDirIfNotExist(target)

	// newName -> oldName，先保存对应关系，最后处理
	hardLinks := make(map[string]string)

	var reader io.ReadCloser
	var err error

	if reader, err = os.Open(tarball); err != nil {
		log.Errorln("faild to open file, err:", err)
	}
	defer reader.Close()

	// 如果以 .gz 结尾，那么要再套一层 gzip 做解压
	if filepath.Ext(tarball) == ".gz" {
		if reader, err = gzip.NewReader(reader); err != nil {
			log.Errorln("failed to new gzip reader, err:", err)
		}
		defer reader.Close()
	}

	tarReader := tar.NewReader(reader)
	for {
		header, err := tarReader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Errorln("faild to call tarReader.Next, err:", err)
		}

		path := filepath.Join(target, header.Name)
		info := header.FileInfo()

		switch header.Typeflag {
		case tar.TypeDir: // 文件夹类型
			if err = os.MkdirAll(path, info.Mode()); err != nil {
				log.Errorln("failed to mkdir, err:", err)
			}

		case tar.TypeLink: // 硬链接，延迟处理
			oldName := filepath.Join(target, header.Linkname)
			newName := filepath.Join(target, header.Name)
			hardLinks[newName] = oldName

		case tar.TypeSymlink: // 软链接，直接连
			linkPath := filepath.Join(target, header.Name)
			if err := os.Symlink(header.Linkname, linkPath); err != nil {
				if os.IsExist(err) {
					continue
				}
				log.Errorf("failed to symlink %q to %q, err: %v\n", linkPath, header.Linkname, err)
			}

		case tar.TypeReg: // 一般文件
			CreateDirIfNotExist(filepath.Dir(path))
			file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
			if err != nil {
				log.Errorln("failed to open file, err:", err)
			}

			_, err = io.Copy(file, tarReader)
			file.Close()
			if err != nil {
				log.Errorln("failed to copy file, err:", err)
			}

		default:
			log.Errorln("unsupported file type", header.Typeflag)
		}
	}

	for k, v := range hardLinks {
		if err := os.Link(v, k); err != nil {
			log.Errorf("failed to hardlink %q to %q, err: %v\n", k, v, err)
		}
	}
}
