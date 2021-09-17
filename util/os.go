package util

import (
	"os"
	"path"
)

var (
	TmpRootPath    = "/var/tmp/pakket"
	TmpPath        string
	TmpSrcPath     string
	TmpPkgPath     string
	TmpPkgRootPath string
	Gid            string
	Uid            string
)

func CreateTempFolder(pkgName string) (err error) {
	TmpPath = path.Join(TmpRootPath, pkgName)
	TmpSrcPath = path.Join(TmpPath, "src")
	TmpPkgRootPath = path.Join(TmpPath, "package")
	TmpPkgPath = path.Join(TmpPkgRootPath, pkgName)

	err = os.RemoveAll(TmpPath)
	if err != nil {
		return
	}

	err = os.MkdirAll(path.Join(TmpRootPath, pkgName), 0774)
	if err != nil {
		return
	}

	return os.MkdirAll(TmpPkgPath, 0774)
}
