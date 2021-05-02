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

func CreateTempFolder() (tmpPath string, err error) {
	os.MkdirAll(TmpRootPath, 0774)

	tmpPath, err = os.MkdirTemp(TmpRootPath, "pkg")
	TmpPath = tmpPath
	TmpSrcPath = path.Join(TmpPath, "src")
	TmpPkgPath = path.Join(TmpPath, "package")

	os.Mkdir(TmpPkgPath, 0774)

	return
}
