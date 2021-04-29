package util

import "path"

var (
	TmpPath    = "/var/tmp/stew"
	TmpSrcPath = path.Join(TmpPath, "src")
	TmpPkgPath = path.Join(TmpPath, "package")
)
