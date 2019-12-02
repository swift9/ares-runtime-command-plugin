// +build darwin

package runtime

import (
	"github.com/swift9/ares-sdk/runtime"
	"log"
	"os/exec"
	"syscall"
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

func (r *CommandRuntime) Start(command runtime.ProcessCommand) int {
	r.Cmd = exec.Command(command.Cmd, command.Args...)
	r.Cmd.Dir = r.Dir

	if len(r.Cmd.Env) == 0 {
		r.Cmd.Env = []string{}
	}
	for _, env := range command.Envs {
		r.Cmd.Env = append(r.Cmd.Env, env.Name+"="+env.Value)
	}

	r.Cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	r.Cmd.Stdout = &logWriter{
		r: r,
	}
	r.Cmd.Stderr = &logWriter{
		r: r,
	}
	err := r.Cmd.Start()
	if err != nil {
		log.Println("start error ", err)
		r.Emit("exit", 1)
		return 1
	}
	r.Emit("ready")
	err = r.Cmd.Wait()
	if err != nil {
		log.Println("wait error ", err)
		r.Emit("exit", 1)
		return 1
	} else {
		log.Println("exit")
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
	return -1
}

func (r *CommandRuntime) Health() runtime.Status {
	return runtime.NewStatusUp()
}

func (r *CommandRuntime) Init() {
}

func NewCommandRuntime(workDir string) runtime.IRuntime {
	var r runtime.IRuntime = &CommandRuntime{
		Dir: workDir,
	}
	r.Init()
	return r
}
