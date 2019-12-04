// +build linux

package runtime

import (
	"os/exec"
	"syscall"
)

func enableGroupKill(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
}

func (r *Runtime) killGroup(pid int) error {
	defer func() {
		if e := recover(); e != nil {
			r.getLogWriter().Write([]byte("unknown"))
		}
	}()
	return syscall.Kill(-pid, syscall.SIGKILL)
}
