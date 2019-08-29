// +build linux

package runtime

import (
	"github.com/swift9/ares-sdk/runtime"
	"log"
	"os/exec"
	"syscall"
	"time"
)

type CommandRuntime struct {
	runtime.Runtime
	Dir string
	Cmd *exec.Cmd
}

type logWriter struct {
	r *CommandRuntime
}

func (w *logWriter) Write(bytes []byte) (n int, err error) {
	w.r.Emit("log", string(bytes))
	return len(bytes), nil
}

func (r *CommandRuntime) Start(cmd string, args ...string) int {
	r.Cmd = exec.Command(cmd, args...)
	r.Cmd.Dir = r.Dir
	r.Cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	r.Cmd.Stdout = &logWriter{
		r: r,
	}
	r.Cmd.Stderr = &logWriter{
		r: r,
	}
	err := r.Cmd.Start()
	if err != nil {
		log.Println("exit:", err)
		r.Emit("exit", 1)
		return 1
	}
	r.Emit("ready")
	time.Sleep(1 * time.Second)
	err = r.Cmd.Wait()
	if err != nil {
		log.Println("exit:", err)
		r.Emit("exit", 1)
		return 1
	} else {
		log.Println("exit:0")
		r.Emit("exit", 0)
		return 0
	}
}

func (r *CommandRuntime) Stop() {
	r.kill()
	r.killAll()
}

func (r *CommandRuntime) kill() {
	defer func() {
		if e := recover(); e != nil {
			log.Println(e)
		}
	}()
	r.Cmd.Process.Kill()
}

func (r *CommandRuntime) killAll() {
	defer func() {
		if e := recover(); e != nil {
			log.Println(e)
		}
	}()
	syscall.Kill(-r.Cmd.Process.Pid, syscall.SIGKILL)
}

func (r *CommandRuntime) Idle() int {
	return 0
}

func (r *CommandRuntime) Health() int {
	return 0
}

func (r *CommandRuntime) Init() {
	r.CreateTime = time.Now()
}

func NewCommandRuntime(workDir string) runtime.IRuntime {
	var r runtime.IRuntime = &CommandRuntime{
		Dir: workDir,
	}
	r.Init()
	return r
}
