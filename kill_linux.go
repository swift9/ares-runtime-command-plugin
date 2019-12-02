// +build linux

package runtime

import (
	"log"
	"os/exec"
	"syscall"
)

func enableGroupKill(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
}

func killGroup(pid int) error {
	defer func() {
		if e := recover(); e != nil {
			log.Println(e)
		}
	}()
	return syscall.Kill(-pid, syscall.SIGKILL)
}
