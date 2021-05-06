package runner

import (
	"os"
	"os/exec"

	"github.com/stewproject/builder/util"
)

var (
	scriptEnv = []string{
		"STEW_PKG_PATH=" + util.TmpPkgPath,
		"STEW_SRC_DIR=" + util.TmpSrcPath,
	}
)

func RunScript(script string, env ...string) error {
	cmd := exec.Command(script)
	scriptEnv = append(scriptEnv, os.Environ()...)
	scriptEnv = append(scriptEnv, env...)
	cmd.Env = scriptEnv

	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	return cmd.Run()
}
