package runner

import (
	"os"
	"os/exec"

	"github.com/pakket-project/builder/util"
)

var (
	scriptEnv []string
)

func RunScript(script string, env ...string) (err error) {
	err = os.Chmod(script, 0755)
	if err != nil {
		return err
	}

	cmd := exec.Command(script)
	scriptEnv = append(scriptEnv, os.Environ()...)
	scriptEnv = append(scriptEnv, env...)
	scriptEnv = append(scriptEnv, "PAKKET_PKG_PATH="+util.TmpPkgPath, "PAKKET_SRC_DIR="+util.TmpSrcPath)
	scriptEnv = append(scriptEnv, "PAKKET_ARCH="+util.Arch)
	cmd.Env = scriptEnv

	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	return cmd.Run()
}
