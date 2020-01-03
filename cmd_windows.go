// +build windows

package runtime

import (
	"fmt"
	"os/exec"
)

func enhanceCmd(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
}

func (r *Runtime) killGroup(pid int) error {
	defer func() {
		if e := recover(); e != nil {
			r.getLogWriter().Write([]byte("unknown"))
		}
	}()
	r.getLogWriter().Write([]byte("taskkill /T " + fmt.Sprint(pid)))
	return exec.Command("taskkill", "/F", "/T", "/PID", fmt.Sprint(pid)).Run()
}
