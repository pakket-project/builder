package runner

import (
	"os"
	"os/exec"
)

var (
	scriptEnv = []string{"STEW_CARGO_ARGS=eee"}
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
