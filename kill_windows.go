// +build windows

package runtime

import (
	"fmt"
	"log"
	"os/exec"
)

func enableGroupKill(cmd *exec.Cmd) {
}

func killGroup(pid int) error {
	defer func() {
		if e := recover(); e != nil {
			log.Println(e)
		}
	}()
	return exec.Command("taskkill", "/F", "/T", "/PID", fmt.Sprint(pid)).Run()
}
