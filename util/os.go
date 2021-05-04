package util

import (
	"os"
	"path"
)

var (
	TmpRootPath = "/var/tmp/stew"
	TmpPath     string
	TmpSrcPath  string
	TmpPkgPath  string
	Gid         string
	Uid         string
)

func CreateTempFolder(pkgName string) (err error) {
	TmpPath = path.Join(TmpRootPath, pkgName)
	TmpSrcPath = path.Join(TmpPath, "src")
	TmpPkgPath = path.Join(TmpPath, "package")

	err = os.MkdirAll(path.Join(TmpRootPath, pkgName), 0774)
	if err != nil {
		return err
	}

	os.Mkdir(TmpPkgPath, 0774)
	return
}
