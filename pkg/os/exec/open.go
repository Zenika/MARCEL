package xexec

import "os/exec"

func Open(path string) error {
	return exec.Command(openCmd, path).Start()
}
